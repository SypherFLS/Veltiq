package orchestrator

import (
	"context"
	"io"

	"veltiq/internal/core/domain"
	"veltiq/internal/core/service"
)

type Runner struct {
	importService *service.ImportService
	reportService *service.ReportService
}

func NewRunner(
	importService *service.ImportService,
	reportService *service.ReportService,
) *Runner {
	return &Runner{
		importService: importService,
		reportService: reportService,
	}
}

func (r *Runner) StartImport(ctx context.Context, tenantID int, raw io.Reader) (string, error) {
	return r.importService.RunImport(ctx, tenantID, raw)
}

func (r *Runner) GetImportStatus(ctx context.Context, importID string) (domain.ImportStatus, error) {
	return r.importService.GetImportStatus(ctx, importID)
}

func (r *Runner) GetImportReport(ctx context.Context, importID string) (domain.Report, error) {
	return r.reportService.BuildReport(ctx, importID)
}