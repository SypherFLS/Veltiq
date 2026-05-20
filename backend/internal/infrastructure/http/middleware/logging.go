package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"veltiq/internal/core/ports"
)

func RequestLogger(log ports.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		log.Info("http_request",
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
			"request_id", c.GetString(RequestIDKey),
		)
	}
}
