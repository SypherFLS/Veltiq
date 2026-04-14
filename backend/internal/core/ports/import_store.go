package ports

import (
	"veltiq/internal/core/domain"
	"context"
)

type ImportStore interface {
	Get(ctx context.Context, args ...any) (any, error)
	Update(ctx context.Context) error
	Validate(ctx context.Context, args ...any) error
	CallReport(ctx context.Context, data Holder, reportCh chan domain.Err) error
}

