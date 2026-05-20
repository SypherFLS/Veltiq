package analytics

import (
	"context"

	"veltiq/internal/core/domain"
)

type StubAnalyzer struct{}

func NewStub() *StubAnalyzer {
	return &StubAnalyzer{}
}

func (a *StubAnalyzer) Analyze(ctx context.Context, importID string, receipts []domain.Receipt) (domain.ReportSummary, error) {
	_ = ctx
	_ = importID

	return domain.ReportSummary{
		ReceiptsCount: len(receipts),
		TotalSum: 0,
		CashSum: 0,
		CardSum: 0,
		IsStub: true,
	}, nil
}
