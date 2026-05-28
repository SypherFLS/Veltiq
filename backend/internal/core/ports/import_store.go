package ports

import (
	"context"
	"veltiq/internal/core/domain"
)

type ImportStore interface {
	GetByID(ctx context.Context, id string) (domain.Import, error)
	GetStatusByID(ctx context.Context, id string) (domain.ImportStatus, error)
	CreateIfAbsent(ctx context.Context, imp domain.Import) (exists bool, err error)
	SetStatus(ctx context.Context, importID string, st domain.ImportStatus) error
	GetByDocumentID(ctx context.Context, tenantID, documentID string) (domain.Import, error)
	ListByTenant(ctx context.Context, tenantID string, limit int) ([]domain.Import, int64, error)
}
