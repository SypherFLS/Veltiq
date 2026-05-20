package logging

import (
	"log/slog"
	"os"

	"veltiq/internal/core/ports"
)

type SlogLogger struct {
	log *slog.Logger
}

func New() ports.Logger {
	return &SlogLogger{
		log: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
	}
}

func (l *SlogLogger) Info(msg string, kv ...any) {
	l.log.Info(msg, kv...)
}

func (l *SlogLogger) Warn(msg string, kv ...any) {
	l.log.Warn(msg, kv...)
}

func (l *SlogLogger) Error(msg string, kv ...any) {
	l.log.Error(msg, kv...)
}
