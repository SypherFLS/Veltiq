package ports

import (
	"context"

	"veltiq/internal/core/domain"
)

type ImportWorkflow interface {
	PersistResults(ctx context.Context, importID string, receipts []domain.Receipt, items []domain.ReceiptItem) error
}
