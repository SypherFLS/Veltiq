package ports

import (
	"context"

	"veltiq/internal/core/domain"
)

type Analytics interface {
	Analyze(ctx context.Context, importID string, receipts []domain.Receipt) (domain.ReportSummary, error)
}
