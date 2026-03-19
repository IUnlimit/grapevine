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

// RateLimitHandler 可直接调用的限流处理
func RateLimitHandler(c *gin.Context) {
	suffix := c.Param("suffix")
	if suffix == "" {
		return
	}

	var cfg model.Config
	var svc model.Service
	if err := database.DB.Where("suffix = ?", suffix).First(&svc).Error; err != nil {
		return
	}
	if err := database.DB.Where("service_id = ?", svc.ID).First(&cfg).Error; err != nil {
		return
	}

	if cfg.RateLimit <= 0 {
		return
	}

	key := suffix + ":" + c.ClientIP()
	limiter := store.Get(key, rate.Limit(cfg.RateLimit), cfg.Burst)

	if !limiter.Allow() {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
	}
}

// RateLimit 返回 gin 中间件
func RateLimit() gin.HandlerFunc {
	return RateLimitHandler
}
