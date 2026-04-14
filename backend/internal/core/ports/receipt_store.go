package ports		

import (
	"veltiq/internal/core/domain"
	"context"
)

type ReceiptStore interface {
	SaveParsed(ctx context.Context, receipts []domain.Receipt) error
	Read(ctx context.Context,) ([]domain.Receipt, error)
}