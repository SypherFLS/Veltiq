package service

import (
	"context"
	"errors"
	"time"

	"veltiq/internal/core/domain"
	"veltiq/internal/core/ports"
)

type ReportService struct {
	imports ports.ImportStore
	receipts ports.ReceiptStore
	analytics ports.Analytics
	logger ports.Logger
}

func NewReportService(
	imports ports.ImportStore,
	receipts ports.ReceiptStore,
	analytics ports.Analytics,
	logger ports.Logger,
) *ReportService {
	return &ReportService{
		imports: imports,
		receipts: receipts,
		analytics: analytics,
		logger: logger,
	}
}

func (s *ReportService) BuildReport(ctx context.Context, importID, tenantID string) (domain.Report, error) {
	if importID == "" {
		return domain.Report{}, domain.InitError("Report", "Build", domain.Invalid_Input, nil, "empty importID", domain.SeverityLocal, false)
	}

	imp, err := s.imports.GetByID(ctx, importID)
	if err != nil {
		if errors.Is(err, domain.ErrImportNotFound) {
			return domain.Report{}, err
		}
		return domain.Report{}, domain.InitError("Report", "GetImport", domain.Store_Failed, err, err.Error(), domain.SeverityLocal, true)
	}

	if imp.TenantID != tenantID {
		return domain.Report{}, domain.ErrImportForbidden
	}

	if imp.Status != domain.ImportDone {
		return domain.Report{}, domain.ErrImportNotReady
	}

	receipts, err := s.receipts.ReadByImport(ctx, importID)
	if err != nil {
		return domain.Report{}, domain.InitError("Report", "ReadReceipts", domain.Store_Failed, err, err.Error(), domain.SeverityLocal, true)
	}

	summary, err := s.analytics.Analyze(ctx, importID, receipts)
	if err != nil {
		return domain.Report{}, domain.InitError("Report", "Analyze", domain.Store_Failed, err, err.Error(), domain.SeverityLocal, false)
	}

	return domain.Report{
		ImportID: importID,
		Status: imp.Status,
		CreatedAt: imp.CreatedAt,
		UpdatedAt: time.Now().UTC(),
		Data: summary,
	}, nil
}

func (s *ReportService) BuildInsights(ctx context.Context, importID, tenantID string) (domain.Insights, error) {
	if importID == "" {
		return domain.Insights{}, domain.InitError("Insights", "Build", domain.Invalid_Input, nil, "empty importID", domain.SeverityLocal, false)
	}

	imp, err := s.imports.GetByID(ctx, importID)
	if err != nil {
		if errors.Is(err, domain.ErrImportNotFound) {
			return domain.Insights{}, err
		}
		return domain.Insights{}, domain.InitError("Insights", "GetImport", domain.Store_Failed, err, err.Error(), domain.SeverityLocal, true)
	}

	if imp.TenantID != tenantID {
		return domain.Insights{}, domain.ErrImportForbidden
	}

	if imp.Status != domain.ImportDone && imp.Status != domain.ImportPartialFail {
		return domain.Insights{}, domain.ErrImportNotReady
	}

	items, err := s.receipts.ReadItemsByImport(ctx, importID)
	if err != nil {
		return domain.Insights{}, domain.InitError("Insights", "ReadItems", domain.Store_Failed, err, err.Error(), domain.SeverityLocal, true)
	}

	return s.analytics.BuildInsights(ctx, importID, items)
}
