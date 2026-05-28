package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"veltiq/internal/core/orchestrator"
	"veltiq/internal/core/ports"
	"veltiq/internal/core/service"
	"veltiq/internal/infrastructure/config"
	"veltiq/internal/infrastructure/cookies"
	"veltiq/internal/infrastructure/http/router"
	"veltiq/internal/infrastructure/jwt"
	"veltiq/internal/infrastructure/logging"
	"veltiq/internal/infrastructure/repository/postgres"
	"veltiq/internal/infrastructure/security"
	"veltiq/internal/modules/analytics"
	"veltiq/internal/modules/parser"
)

type App struct {
	Config *config.Config
	DB *gorm.DB
	Router *gin.Engine
	Logger ports.Logger
}

func Bootstrap(cfg *config.Config) (*App, error) {
	log := logging.New()

	db, err := postgres.InitDB(cfg)
	if err != nil {
		return nil, err
	}

	userRepo := postgres.NewUserRepository(db)
	importRepo := postgres.NewImportRepository(db)
	importWorkflow := postgres.NewImportWorkflow(db)

	passwordManager := &security.PasswordManager{}
	jwtManager := jwt.NewManager(cfg.JWT.SecretKey, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	cookieManager := cookies.NewManager(
		cfg.Cookies.Secure,
		int(cfg.JWT.AccessTTL.Seconds()),
		int(cfg.JWT.RefreshTTL.Seconds()),
	)

	authService := service.NewAuthService(userRepo, passwordManager)

	csvParser := parser.NewCSV()
	analyzer := analytics.New()

	importService := service.NewImportService(importRepo, importWorkflow, csvParser, log)
	reportService := service.NewReportService(importRepo, postgres.NewReceiptRepository(db), analyzer, log)
	runner := orchestrator.NewRunner(importService, reportService)

	r := router.NewRouter(router.Deps{
		Config: cfg,
		AuthService: authService,
		Tokens: jwtManager,
		Cookies: cookieManager,
		Runner: runner,
		Logger: log,
	})

	return &App{
		Config: cfg,
		DB: db,
		Router: r,
		Logger: log,
	}, nil
}
