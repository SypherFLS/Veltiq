package postgres

import (
	"context"

	"veltiq/internal/core/domain"

	"gorm.io/gorm"
)

type ReceiptRepository struct {
	db *gorm.DB
}

func NewReceiptRepository(db *gorm.DB) *ReceiptRepository {
	return &ReceiptRepository{db: db}
}

func (r *ReceiptRepository) SaveParsed(ctx context.Context, importID string, receipts []domain.Receipt) error {
	if len(receipts) == 0 {
		return nil
	}

	records := make([]receiptRecord, 0, len(receipts))
	for _, rc := range receipts {
		rc.ImportID = importID
		records = append(records, receiptFromDomain(rc))
	}

	return r.db.WithContext(ctx).Create(&records).Error
}

func (r *ReceiptRepository) ReadByImport(ctx context.Context, importID string) ([]domain.Receipt, error) {
	var records []receiptRecord
	err := r.db.WithContext(ctx).Where("import_id = ?", importID).Find(&records).Error
	if err != nil {
		return nil, err
	}

	out := make([]domain.Receipt, 0, len(records))
	for _, rec := range records {
		out = append(out, receiptToDomain(rec))
	}
	return out, nil
}
