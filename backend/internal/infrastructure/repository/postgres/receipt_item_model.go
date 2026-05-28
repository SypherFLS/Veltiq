package postgres

import (
	"time"

	"veltiq/internal/core/domain"
)

type receiptItemRecord struct {
	ID                uint      `gorm:"column:id;primaryKey;autoIncrement"`
	ImportID          string    `gorm:"column:import_id;not null;index"`
	TenantID          string    `gorm:"column:tenant_id;type:uuid;not null;index"`
	ReceiptExternalID string    `gorm:"column:receipt_external_id;index"`
	SKU               string    `gorm:"column:sku;not null;index"`
	Name              string    `gorm:"column:name;not null"`
	Category          string    `gorm:"column:category"`
	Quantity          int       `gorm:"column:quantity;not null"`
	UnitPrice         int       `gorm:"column:unit_price;not null"`
	TotalPrice        int       `gorm:"column:total_price;not null"`
	PaymentType       string    `gorm:"column:payment_type"`
	SoldAt            time.Time `gorm:"column:sold_at;index"`
}

func (receiptItemRecord) TableName() string {
	return "receipt_items"
}

func receiptItemFromDomain(r domain.ReceiptItem) receiptItemRecord {
	return receiptItemRecord{
		ImportID:          r.ImportID,
		TenantID:          r.TenantID,
		ReceiptExternalID: r.ReceiptExternalID,
		SKU:               r.SKU,
		Name:              r.Name,
		Category:          r.Category,
		Quantity:          r.Quantity,
		UnitPrice:         r.UnitPrice,
		TotalPrice:        r.TotalPrice,
		PaymentType:       string(r.PaymentType),
		SoldAt:            r.SoldAt,
	}
}

func receiptItemToDomain(r receiptItemRecord) domain.ReceiptItem {
	return domain.ReceiptItem{
		ID:                int(r.ID),
		ImportID:          r.ImportID,
		TenantID:          r.TenantID,
		ReceiptExternalID: r.ReceiptExternalID,
		SKU:               r.SKU,
		Name:              r.Name,
		Category:          r.Category,
		Quantity:          r.Quantity,
		UnitPrice:         r.UnitPrice,
		TotalPrice:        r.TotalPrice,
		PaymentType:       domain.PaymentType(r.PaymentType),
		SoldAt:            r.SoldAt,
	}
}
