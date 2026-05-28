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

func (r *ReceiptRepository) SaveItems(ctx context.Context, importID string, items []domain.ReceiptItem) error {
	if len(items) == 0 {
		return nil
	}

	records := make([]receiptItemRecord, 0, len(items))
	for _, it := range items {
		it.ImportID = importID
		records = append(records, receiptItemFromDomain(it))
	}

	return r.db.WithContext(ctx).CreateInBatches(records, 500).Error
}

func (r *ReceiptRepository) ReadItemsByImport(ctx context.Context, importID string) ([]domain.ReceiptItem, error) {
	var records []receiptItemRecord
	err := r.db.WithContext(ctx).Where("import_id = ?", importID).Find(&records).Error
	if err != nil {
		return nil, err
	}

	out := make([]domain.ReceiptItem, 0, len(records))
	for _, rec := range records {
		out = append(out, receiptItemToDomain(rec))
	}
	return out, nil
}
