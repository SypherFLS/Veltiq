package orchestrator

import (
	"bytes"
	"context"
	"testing"

	"veltiq/internal/core/domain"
	"veltiq/internal/core/service"
	"veltiq/internal/core/testhelpers"
)

func newRunnerTestDeps() (*testhelpers.FakeImportStore, *testhelpers.FakeReceiptStore, *testhelpers.FakeParser, *testhelpers.FakeLogger, *Runner) {
	importStore := &testhelpers.FakeImportStore{
		CreateIfAbsentCreated: true,
		GetStatusByIDStatus:   domain.ImportDone,
	}
	receiptStore := &testhelpers.FakeReceiptStore{}
	parser := &testhelpers.FakeParser{}
	logger := &testhelpers.FakeLogger{}

	importService := service.NewImportService(importStore, receiptStore, parser, logger)
	reportService := service.NewReportService(importStore, receiptStore, logger)
	runner := NewRunner(importService, reportService)

	return importStore, receiptStore, parser, logger, runner
}

func TestRunner_StartImport_HappyPath(t *testing.T) {
	importStore, receiptStore, parser, _, runner := newRunnerTestDeps()
	parser.Receipts = []domain.Receipt{
		{Sum: 100, TypeOfPayment: domain.PaymentByCash},
	}

	importID, err := runner.StartImport(context.Background(), 42, bytes.NewReader([]byte("raw payload")))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if importID == "" {
		t.Fatal("expected non-empty import ID")
	}
	if parser.ParseCalls != 1 {
		t.Fatalf("expected parser to be called once, got %d", parser.ParseCalls)
	}
	if receiptStore.SaveCalls != 1 {
		t.Fatalf("expected SaveParsed to be called once, got %d", receiptStore.SaveCalls)
	}
	if importStore.GetByDocumentIDCalls != 0 {
		t.Fatalf("expected no duplicate path usage, got %d calls", importStore.GetByDocumentIDCalls)
	}
}

func TestRunner_GetImportStatus_Delegates(t *testing.T) {
	_, _, _, _, runner := newRunnerTestDeps()

	status, err := runner.GetImportStatus(context.Background(), "imp-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if status != domain.ImportDone {
		t.Fatalf("expected status %q, got %q", domain.ImportDone, status)
	}
}

func TestRunner_GetImportReport_Delegates(t *testing.T) {
	_, receiptStore, _, _, runner := newRunnerTestDeps()
	receiptStore.ReadByImportResult = []domain.Receipt{
		{Sum: 100, TypeOfPayment: domain.PaymentByCash},
		{Sum: 250, TypeOfPayment: domain.PaymentByCard},
	}

	report, err := runner.GetImportReport(context.Background(), "imp-1")
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

