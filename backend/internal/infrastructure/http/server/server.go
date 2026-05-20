package server

import (
	"context"
	"errors"
	"net/http"

	"veltiq/internal/infrastructure/config"
)

func Run(ctx context.Context, cfg *config.Config, handler http.Handler) error {
	srv := &http.Server{
		Addr: cfg.HTTP.Address,
		Handler: handler,
		ReadTimeout: cfg.HTTP.Timeout,
		WriteTimeout: cfg.HTTP.Timeout,
		IdleTimeout: cfg.HTTP.IdleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.Timeout)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
