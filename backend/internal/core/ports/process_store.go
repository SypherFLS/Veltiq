package ports 

import (
	"veltiq/internal/core/domain"
	"context"
)

type Holder struct {
	data any
}

type ProcessModule interface {
	Start(ctx context.Context, receipts []domain.Receipt) error
	Finish(ctx context.Context,) error
	SetStatus(ctx context.Context, ProcessID string, st domain.ImportStatus) error
	Read(ctx context.Context, receipts []domain.Receipt) error
	SaveResult(ctx context.Context, data Holder, resultCh chan any) error
	CallReport(ctx context.Context, data Holder, reportCh chan domain.Err) error
}