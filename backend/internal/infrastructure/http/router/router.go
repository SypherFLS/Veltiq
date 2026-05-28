package router

import (
	"github.com/gin-gonic/gin"

	"veltiq/internal/core/orchestrator"
	"veltiq/internal/core/ports"
	"veltiq/internal/core/service"
	"veltiq/internal/infrastructure/config"
	"veltiq/internal/infrastructure/cookies"
	v1 "veltiq/internal/infrastructure/http/handlers/api/v1"
	"veltiq/internal/infrastructure/http/middleware"
)

type Deps struct {
	Config *config.Config
	AuthService *service.AuthService
	Tokens ports.TokenManager
	Cookies *cookies.Manager
	Runner *orchestrator.Runner
	Logger ports.Logger
}

func NewRouter(deps Deps) *gin.Engine {
	router := gin.New()
	router.Use(
		middleware.CORS(deps.Config.CORS.AllowedOrigins),
		middleware.RequestID(),
		middleware.Recovery(deps.Logger),
		middleware.RequestLogger(deps.Logger),
		middleware.MaxBodyBytes(deps.Config.HTTP.MaxBodyBytes),
	)

	api := router.Group("/api/v1")

	authHandler := v1.NewAuthHandler(deps.AuthService, deps.Tokens, deps.Cookies)
	authLimit := middleware.AuthRateLimit(deps.Config.Auth.RateLimitPerMinute)

	api.POST("/register", authLimit, authHandler.Register)
	api.POST("/login", authLimit, authHandler.Login)
	api.GET("/auth/session", authHandler.VerifySession)
	api.POST("/auth/refresh", authLimit, authHandler.Refresh)
	api.POST("/auth/logout", authHandler.Logout)

	protected := api.Group("/")
	protected.Use(
		middleware.AuthMiddleware(deps.Tokens, deps.Cookies),
		middleware.TenantMiddleware(),
	)

	profileHandler := v1.NewProfileHandler()
	protected.GET("/profile", profileHandler.Get)

	importHandler := v1.NewImportHandler(deps.Runner)
	protected.GET("/imports", importHandler.List)
	protected.POST("/imports", importHandler.Upload)
	protected.GET("/imports/:id/status", importHandler.Status)
	protected.GET("/imports/:id/report", importHandler.Report)

	return router
}
