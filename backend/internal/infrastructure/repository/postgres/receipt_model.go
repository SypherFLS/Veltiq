package postgres

import (
	"time"

	"veltiq/internal/core/domain"
)

type receiptRecord struct {
	ID          uint      `gorm:"column:id;primaryKey;autoIncrement"`
	ImportID    string    `gorm:"column:import_id;not null;index"`
	TenantID    string    `gorm:"column:tenant_id;type:uuid;not null;index"`
	ExternalID  string    `gorm:"column:external_id;index"`
	StoreID     int       `gorm:"column:store_id"`
	PaymentType string    `gorm:"column:payment_type"`
	Sum         int       `gorm:"column:sum"`
	Cashier     string    `gorm:"column:cashier"`
	CreatedAt   time.Time `gorm:"column:created_at;index"`
}

func (receiptRecord) TableName() string {
	return "receipts"
}

func receiptFromDomain(r domain.Receipt) receiptRecord {
	return receiptRecord{
		ImportID:    r.ImportID,
		TenantID:    r.TenantID,
		ExternalID:  r.ExternalID,
		StoreID:     r.StoreID,
		PaymentType: string(r.TypeOfPayment),
		Sum:         r.Sum,
		Cashier:     r.Cashier,
		CreatedAt:   r.CreatedAt,
	}
}

func receiptToDomain(r receiptRecord) domain.Receipt {
	return domain.Receipt{
		ID:            int(r.ID),
		ImportID:      r.ImportID,
		TenantID:      r.TenantID,
		ExternalID:    r.ExternalID,
		StoreID:       r.StoreID,
		TypeOfPayment: domain.PaymentType(r.PaymentType),
		Sum:           r.Sum,
		Cashier:       r.Cashier,
		CreatedAt:     r.CreatedAt,
	}
}
