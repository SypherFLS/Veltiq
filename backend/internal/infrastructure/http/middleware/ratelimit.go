package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	mu sync.Mutex
	limit int
	window time.Duration
	buckets map[string][]time.Time
}

func newRateLimiter(limitPerMinute int) *rateLimiter {
	if limitPerMinute <= 0 {
		limitPerMinute = 20
	}
	return &rateLimiter{
		limit: limitPerMinute,
		window: time.Minute,
		buckets: make(map[string][]time.Time),
	}
}

func (r *rateLimiter) allow(key string) bool {
	now := time.Now()
	cutoff := now.Add(-r.window)

	r.mu.Lock()
	defer r.mu.Unlock()

	times := r.buckets[key]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	if len(filtered) >= r.limit {
		r.buckets[key] = filtered
		return false
	}
	filtered = append(filtered, now)
	r.buckets[key] = filtered
	return true
}

func AuthRateLimit(limitPerMinute int) gin.HandlerFunc {
	rl := newRateLimiter(limitPerMinute)
	return func(c *gin.Context) {
		key := c.ClientIP()
		if !rl.allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}
		c.Next()
	}
}
