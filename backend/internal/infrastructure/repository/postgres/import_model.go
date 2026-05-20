package postgres

import (
	"time"

	"veltiq/internal/core/domain"
)

type importRecord struct {
	ID string `gorm:"column:id;primaryKey;type:uuid"`
	DocumentID string `gorm:"column:document_id;not null;uniqueIndex:idx_import_tenant_document,priority:2"`
	TenantID string `gorm:"column:tenant_id;type:uuid;not null;uniqueIndex:idx_import_tenant_document,priority:1"`
	Status string `gorm:"column:status;not null"`
	ErrorCode string `gorm:"column:error_code"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (importRecord) TableName() string {
	return "imports"
}

func importFromDomain(i domain.Import) importRecord {
	return importRecord{
		ID: i.ID,
		DocumentID: i.DocumentID,
		TenantID: i.TenantID,
		Status: string(i.Status),
		ErrorCode: i.ErrorCode,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}

func importToDomain(r importRecord) domain.Import {
	return domain.Import{
		ID: r.ID,
		DocumentID: r.DocumentID,
		TenantID: r.TenantID,
		Status: domain.ImportStatus(r.Status),
		ErrorCode: r.ErrorCode,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
