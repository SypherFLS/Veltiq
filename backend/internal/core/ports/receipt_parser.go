package ports

import (
	"veltiq/internal/core/domain"
	"context"
	"io"
)

type ReceiptParser interface {
	Parse(ctx context.Context, importID string, raw io.Reader) ([]domain.Receipt, error)
}