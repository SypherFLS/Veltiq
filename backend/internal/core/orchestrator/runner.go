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

func (r *Runner) StartImport(ctx context.Context, tenantID string, raw io.Reader) (string, error) {
	return r.importService.RunImport(ctx, tenantID, raw)
}

func (r *Runner) GetImport(ctx context.Context, importID, tenantID string) (domain.Import, error) {
	return r.importService.GetImport(ctx, importID, tenantID)
}

func (r *Runner) GetImportStatus(ctx context.Context, importID, tenantID string) (domain.ImportStatus, error) {
	return r.importService.GetImportStatus(ctx, importID, tenantID)
}

func (r *Runner) ListImports(ctx context.Context, tenantID string, limit int) ([]domain.Import, int64, error) {
	return r.importService.ListImports(ctx, tenantID, limit)
}

func (r *Runner) GetImportReport(ctx context.Context, importID, tenantID string) (domain.Report, error) {
	return r.reportService.BuildReport(ctx, importID, tenantID)
}
