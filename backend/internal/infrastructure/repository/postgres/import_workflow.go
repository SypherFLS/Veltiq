package postgres

import (
	"context"

	"veltiq/internal/core/domain"

	"gorm.io/gorm"
)

type ImportWorkflow struct {
	db *gorm.DB
}

func NewImportWorkflow(db *gorm.DB) *ImportWorkflow {
	return &ImportWorkflow{db: db}
}

func (w *ImportWorkflow) PersistResults(
	ctx context.Context,
	importID string,
	receipts []domain.Receipt,
	items []domain.ReceiptItem,
) error {
	return w.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		receiptRepo := NewReceiptRepository(tx)
		if err := receiptRepo.SaveParsed(ctx, importID, receipts); err != nil {
			return err
		}
		if err := receiptRepo.SaveItems(ctx, importID, items); err != nil {
			return err
		}

		importRepo := NewImportRepository(tx)
		return importRepo.SetStatus(ctx, importID, domain.ImportDone)
	})
}
