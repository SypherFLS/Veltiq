package router

import (
	"github.com/gin-gonic/gin"

	"veltiq/internal/core/ports"
	"veltiq/internal/core/service"

	v1 "veltiq/internal/infrastructure/http/handlers/api/v1"

	"veltiq/internal/infrastructure/http/middleware"
)

func NewRouter(authService *service.AuthService, tokens ports.TokenManager) *gin.Engine {
	router := gin.Default()

	authHandler := v1.NewAuthHandler(
		authService,
	)

	router.POST(
		"/register",
		authHandler.Register,
	)

	router.POST(
		"/login",
		authHandler.Login,
	)

	protected := router.Group("/")

	protected.Use(
		middleware.AuthMiddleware(tokens),
	)

	protected.GET(
		"/profile",
		func(c *gin.Context) {

			userID := c.GetString("userID")

			c.JSON(200, gin.H{
				"user_id": userID,
			})
		},
	)

	return router
}