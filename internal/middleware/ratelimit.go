package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/illtamer/grapevine/internal/database"
	"github.com/illtamer/grapevine/internal/model"
	"golang.org/x/time/rate"
)

type RateLimiterStore struct {
	mu       sync.RWMutex
	limiters map[string]*rate.Limiter
}

var store = &RateLimiterStore{
	limiters: make(map[string]*rate.Limiter),
}

func (s *RateLimiterStore) Get(key string, r rate.Limit, burst int) *rate.Limiter {
	s.mu.RLock()
	lim, ok := s.limiters[key]
	s.mu.RUnlock()
	if ok {
		return lim
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	lim = rate.NewLimiter(r, burst)
	s.limiters[key] = lim
	return lim
}

func (s *RateLimiterStore) Update(key string, r rate.Limit, burst int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if lim, ok := s.limiters[key]; ok {
		lim.SetLimit(r)
		lim.SetBurst(burst)
	}
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		suffix := c.Param("suffix")
		if suffix == "" {
			c.Next()
			return
		}

		var cfg model.Config
		var svc model.Service
		if err := database.DB.Where("suffix = ?", suffix).First(&svc).Error; err != nil {
			c.Next()
			return
		}
		if err := database.DB.Where("service_id = ?", svc.ID).First(&cfg).Error; err != nil {
			c.Next()
			return
		}

		if cfg.RateLimit <= 0 {
			c.Next()
			return
		}

		// 按 suffix + IP 维度限流
		key := suffix + ":" + c.ClientIP()
		limiter := store.Get(key, rate.Limit(cfg.RateLimit), cfg.Burst)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
