package service

import (
	"context"
	"errors"
	"testing"

	"veltiq/internal/core/domain"
	"veltiq/internal/core/ports"
)

type fakeReportImportStore struct {
	status domain.ImportStatus
	statusErr error
}

func (f *fakeReportImportStore) GetStatusByID(ctx context.Context, id string) (domain.ImportStatus, error) {
	return f.status, f.statusErr
}

func (f *fakeReportImportStore) CreateIfAbsent(ctx context.Context, imp domain.Import) (bool, error) {
	return true, nil
}
func (f *fakeReportImportStore) SetStatus(ctx context.Context, importID string, st domain.ImportStatus) error {
	return nil
}
func (f *fakeReportImportStore) GetByDocumentID(ctx context.Context, documentID string) (domain.Import, error) {
	return domain.Import{}, nil
}
func (f *fakeReportImportStore) CallReport(ctx context.Context, data ports.Holder, reportCh chan domain.Err) error {
	return nil
}

type fakeReportReceiptStore struct {
	receipts []domain.Receipt
	readErr  error
}

func (f *fakeReportReceiptStore) SaveParsed(ctx context.Context, importID string, receipts []domain.Receipt) error {
	return nil
}
func (f *fakeReportReceiptStore) ReadByImport(ctx context.Context, importID string) ([]domain.Receipt, error) {
	if f.readErr != nil {
		return nil, f.readErr
	}
	return f.receipts, nil
}

type fakeReportLogger struct{}

func (l *fakeReportLogger) Info(msg string, kv ...any)  {}
func (l *fakeReportLogger) Warn(msg string, kv ...any)  {}
func (l *fakeReportLogger) Error(msg string, kv ...any) {}

func newReportServiceTestDeps() (*fakeReportImportStore, *fakeReportReceiptStore, *fakeReportLogger, *ReportService) {
	imports := &fakeReportImportStore{status: domain.ImportDone}
	receipts := &fakeReportReceiptStore{}
	logger := &fakeReportLogger{}
	svc := NewReportService(imports, receipts, logger)
	return imports, receipts, logger, svc
}

func TestBuildReport_HappyPath_Aggregates(t *testing.T) {
	imports, receipts, _, svc := newReportServiceTestDeps()

	imports.status = domain.ImportDone
	receipts.receipts = []domain.Receipt{
		{Sum: 100, TypeOfPayment: domain.PaymentByCash},
		{Sum: 250, TypeOfPayment: domain.PaymentByCard},
		{Sum: 50, TypeOfPayment: domain.PaymentByCash},
	}

	report, err := svc.BuildReport(context.Background(), "imp-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if report.ImportID != "imp-1" {
		t.Fatalf("expected importID imp-1, got %s", report.ImportID)
	}
	if report.Status != domain.ImportDone {
		t.Fatalf("expected status %q, got %q", domain.ImportDone, report.Status)
	}

	if len(report.Data) == 0 {
		t.Fatal("expected non-empty report data")
	}
}

func TestBuildReport_GetStatusError(t *testing.T) {
	imports, _, _, svc := newReportServiceTestDeps()
	imports.statusErr = errors.New("status store failed")

	_, err := svc.BuildReport(context.Background(), "imp-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestBuildReport_ReadByImportError(t *testing.T) {
	_, receipts, _, svc := newReportServiceTestDeps()
	receipts.readErr = errors.New("read receipts failed")

	_, err := svc.BuildReport(context.Background(), "imp-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestBuildReport_InvalidInput(t *testing.T) {
	_, _, _, svc := newReportServiceTestDeps()

	_, err := svc.BuildReport(nil, "imp-1")
	if err == nil {
		t.Fatal("expected error for nil ctx")
	}

	_, err = svc.BuildReport(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty importID")
	}
}