package ports

import (
	"context"
	"io"
	"veltiq/internal/core/domain"
)

type ReceiptParser interface {
	Parse(ctx context.Context, importID string, raw io.Reader) ([]domain.Receipt, error)
}
