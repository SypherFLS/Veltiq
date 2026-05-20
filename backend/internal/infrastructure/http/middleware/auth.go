package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"veltiq/internal/core/ports"
)

func AuthMiddleware(
	tokens ports.TokenManager,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		cookie, err := c.Request.Cookie("access_token")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, err := tokens.Verify(cookie.Value)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userID", userID)

		c.Next()
	}
}