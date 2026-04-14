package ports

import (
	"veltiq/internal/core/domain"
	"context"
)

type ImportStore interface {
	CreateIfAbsent(ctx context.Context, imp domain.Import) (bool, error) 
	SetStatus(ctx context.Context, importID string, st domain.ImportStatus) error
	GetByDocumentID(ctx context.Context, documentID string) (domain.Import, error)
	CallReport(ctx context.Context, data Holder, reportCh chan domain.Err) error
}

