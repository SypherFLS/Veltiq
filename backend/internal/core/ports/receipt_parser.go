package ports

import (
	"veltiq/internal/core/domain"
	"context"
)

type ReceiptParser interface {
	Parse(ctx context.Context, args ...any) ([]domain.Receipt, error)
	CallReport(ctx context.Context, data Holder, reportCh chan domain.Err) error
}