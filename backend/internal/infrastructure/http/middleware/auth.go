package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"veltiq/internal/core/ports"
	"veltiq/internal/infrastructure/cookies"
)

func AuthMiddleware(tokens ports.TokenManager, cookieManager *cookies.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := cookieManager.AccessTokenFromRequest(c.Request)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, err := tokens.VerifyAccess(token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}
