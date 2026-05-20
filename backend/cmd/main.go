package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"veltiq/internal/app"
	"veltiq/internal/infrastructure/config"
	httpserver "veltiq/internal/infrastructure/http/server"
)

func main() {
	cfg := config.MustLoad()

	application, err := app.Bootstrap(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("server started on", cfg.HTTP.Address)

	if err := httpserver.Run(ctx, cfg, application.Router); err != nil {
		log.Fatal(err)
	}

	log.Println("server stopped")
}
