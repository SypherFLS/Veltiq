package service

import (
	"veltiq/internal/core/domain"
	"veltiq/internal/core/ports"
	"context"
	"time"
)

type ReportService struct {
	imports  ports.ImportStore
	receipts ports.ReceiptStore
	logger   ports.Logger
}

func NewReportService(imports ports.ImportStore, receipts ports.ReceiptStore, logger ports.Logger) *ReportService {
	return &ReportService{
		imports:  imports,
		receipts: receipts,
		logger:   logger,
	}
}

func (s *ReportService) BuildReport(ctx context.Context, importID string) (domain.Report, error) {
	if ctx == nil {
		return domain.Report{}, domain.InitError("Report", "Build", domain.Invalid_Input, nil, "empty ctx", domain.Local, false)
	}
	if importID == "" {
		return domain.Report{}, domain.InitError("Report", "Build", domain.Invalid_Input, nil, "empty importID", domain.Local, false)
	}

	status, err := s.imports.GetStatusByID(ctx, importID)
	if err != nil {
		return domain.Report{}, domain.InitError("Report", "GetStatus", domain.Store_Failed, err, err.Error(), domain.Local, true)
	}

	receipts, err := s.receipts.ReadByImport(ctx, importID)
	if err != nil {
		return domain.Report{}, domain.InitError("Report", "ReadReceipts", domain.Store_Failed, err, err.Error(), domain.Local, true)
	}

	var totalSum, cashSum, cardSum int
	for _, r := range receipts {
		totalSum += r.Sum
		switch r.TypeOfPayment {
		case domain.PaymentByCash:
			cashSum += r.Sum
		case domain.PaymentByCard:
			cardSum += r.Sum
		}
	}

	report := domain.Report{
		ImportID:  importID,
		Status:    status,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Data: []any{
			map[string]any{"receiptsCount": len(receipts)},
			map[string]any{"totalSum": totalSum},
			map[string]any{"cashSum": cashSum},
			map[string]any{"cardSum": cardSum},
		},
	}

	return report, nil
}