package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/illtamer/grapevine/internal/database"
	"github.com/illtamer/grapevine/internal/model"
)

type CircuitState int

const (
	StateClosed   CircuitState = iota // 正常
	StateOpen                         // 熔断
	StateHalfOpen                     // 半开探测
)

type CircuitBreaker struct {
	mu            sync.Mutex
	state         CircuitState
	failures      int
	successes     int
	total         int
	threshold     int // 错误率百分比阈值
	openUntil     time.Time
	cooldown      time.Duration
	halfOpenMax   int
	windowStart   time.Time
	windowSize    time.Duration
}

func NewCircuitBreaker(threshold int) *CircuitBreaker {
	return &CircuitBreaker{
		state:       StateClosed,
		threshold:   threshold,
		cooldown:    10 * time.Second,
		halfOpenMax: 5,
		windowSize:  30 * time.Second,
		windowStart: time.Now(),
	}
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateOpen:
		if time.Now().After(cb.openUntil) {
			cb.state = StateHalfOpen
			cb.successes = 0
			cb.failures = 0
			cb.total = 0
			return true
		}
		return false
	case StateHalfOpen:
		return cb.total < cb.halfOpenMax
	default:
		return true
	}
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.successes++
	cb.total++

	if cb.state == StateHalfOpen && cb.successes >= cb.halfOpenMax {
		cb.state = StateClosed
		cb.resetWindow()
	}
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	cb.total++

	if cb.state == StateHalfOpen {
		cb.trip()
		return
	}

	// 滑动窗口检查
	if time.Since(cb.windowStart) > cb.windowSize {
		cb.resetWindow()
		cb.failures = 1
		cb.total = 1
		return
	}

	if cb.total >= 10 && cb.threshold > 0 {
		errorRate := (cb.failures * 100) / cb.total
		if errorRate >= cb.threshold {
			cb.trip()
		}
	}
}

func (cb *CircuitBreaker) trip() {
	cb.state = StateOpen
	cb.openUntil = time.Now().Add(cb.cooldown)
}

func (cb *CircuitBreaker) resetWindow() {
	cb.failures = 0
	cb.successes = 0
	cb.total = 0
	cb.windowStart = time.Now()
}

func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// 全局熔断器存储
var (
	cbMu       sync.RWMutex
	breakers   = make(map[string]*CircuitBreaker)
)

func GetBreaker(suffix string, threshold int) *CircuitBreaker {
	cbMu.RLock()
	cb, ok := breakers[suffix]
	cbMu.RUnlock()
	if ok {
		return cb
	}

	cbMu.Lock()
	defer cbMu.Unlock()
	cb = NewCircuitBreaker(threshold)
	breakers[suffix] = cb
	return cb
}

// CircuitBreakHandler 可直接调用的熔断处理（仅检查是否放行）
func CircuitBreakHandler(c *gin.Context) {
	suffix := c.Param("suffix")
	if suffix == "" {
		return
	}

	var svc model.Service
	if err := database.DB.Where("suffix = ?", suffix).First(&svc).Error; err != nil {
		return
	}
	var cfg model.Config
	if err := database.DB.Where("service_id = ?", svc.ID).First(&cfg).Error; err != nil {
		return
	}

	if cfg.CircuitBreakThreshold <= 0 {
		return
	}

	cb := GetBreaker(suffix, cfg.CircuitBreakThreshold)
	if !cb.Allow() {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "service circuit breaker open"})
		return
	}

	// 记录结果需要在响应之后，由调用方负责
	c.Set("_circuit_breaker", cb)
}

// CircuitBreakRecord 在代理完成后记录结果
func CircuitBreakRecord(c *gin.Context) {
	if v, ok := c.Get("_circuit_breaker"); ok {
		cb := v.(*CircuitBreaker)
		if c.Writer.Status() >= 500 {
			cb.RecordFailure()
		} else {
			cb.RecordSuccess()
		}
	}
}

// CircuitBreak 返回 gin 中间件
func CircuitBreak() gin.HandlerFunc {
	return func(c *gin.Context) {
		CircuitBreakHandler(c)
		if c.IsAborted() {
			return
		}
		c.Next()
		CircuitBreakRecord(c)
	}
}
