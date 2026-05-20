package middleware

import (
	"github.com/gin-gonic/gin"

	"veltiq/internal/core/domain"
)

func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(UserIDKey)
		if userID != "" {
			c.Set(TenantIDKey, domain.TenantIDFromUserID(userID))
		}
		c.Next()
	}
}
