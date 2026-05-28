package ports

import (
	"context"
	"veltiq/internal/core/domain"
)

type ReceiptStore interface {
	SaveParsed(ctx context.Context, importID string, receipts []domain.Receipt) error
	SaveItems(ctx context.Context, importID string, items []domain.ReceiptItem) error
	ReadByImport(ctx context.Context, importID string) ([]domain.Receipt, error)
	ReadItemsByImport(ctx context.Context, importID string) ([]domain.ReceiptItem, error)
}
