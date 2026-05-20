package ports

import (
	"context"
	"veltiq/internal/core/domain"
)

type ReceiptStore interface {
	SaveParsed(ctx context.Context, importID string, receipts []domain.Receipt) error
	ReadByImport(ctx context.Context, importID string) ([]domain.Receipt, error)
}
