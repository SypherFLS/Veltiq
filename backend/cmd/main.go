package main

import (
	"log"

	"veltiq/internal/core/service"

	"veltiq/internal/infrastructure/config"
	"veltiq/internal/infrastructure/http/router"
	"veltiq/internal/infrastructure/jwt"
	"veltiq/internal/infrastructure/repository/postgres"
	"veltiq/internal/infrastructure/security"
)

func main() {
	cfg := config.MustLoad()

	db, err := postgres.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := postgres.NewUserRepository(db)

	passwordManager := &security.PasswordManager{}

	jwtManager := jwt.NewManager(
		cfg.JWT.SecretKey,
	)

	authService := service.NewAuthService(
		userRepo,
		passwordManager,
		jwtManager,
	)

	r := router.NewRouter(
		authService,
		jwtManager,
	)

	log.Println(
		"server started on",
		cfg.HTTP.Address,
	)

	err = r.Run(cfg.HTTP.Address)
	if err != nil {
		log.Fatal(err)
	}
}