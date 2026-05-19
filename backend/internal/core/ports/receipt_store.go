package ports		

import (
	"veltiq/internal/core/domain"
	"context"
)

type ReceiptStore interface {
	SaveParsed(ctx context.Context, importID string, receipts []domain.Receipt) error
	ReadByImport(ctx context.Context, importID string) ([]domain.Receipt, error)
}