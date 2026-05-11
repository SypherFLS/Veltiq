package service

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"io"

	"veltiq/internal/core/domain"
	"veltiq/internal/core/ports"
)

// fake structs 
// need to write real mock in the next updates

type fakeImportStore struct {
	createIfAbsentCreated bool
	createIfAbsentErr     error

	setStatusErrByStatus map[domain.ImportStatus]error
	setStatusCalls       []domain.ImportStatus
	getByDocumentIDCalls int

	getStatusByIDStatus domain.ImportStatus
	getStatusByIDErr    error
}

func (f *fakeImportStore) GetStatusByID(ctx context.Context, id string) (domain.ImportStatus, error) {
	return f.getStatusByIDStatus, f.getStatusByIDErr
}

func (f *fakeImportStore) CreateIfAbsent(ctx context.Context, imp domain.Import) (bool, error) {
	return f.createIfAbsentCreated, f.createIfAbsentErr
}

func (f *fakeImportStore) SetStatus(ctx context.Context, importID string, st domain.ImportStatus) error {
	f.setStatusCalls = append(f.setStatusCalls, st)
	if f.setStatusErrByStatus == nil {
		return nil
	}
	return f.setStatusErrByStatus[st]
}

func (f *fakeImportStore) GetByDocumentID(ctx context.Context, documentID string) (domain.Import, error) {
	f.getByDocumentIDCalls++
	return domain.Import{}, nil
}


func (f *fakeImportStore) CallReport(ctx context.Context, data ports.Holder, reportCh chan domain.Err) error {
	return nil
}

type fakeReceiptStore struct {
	saveParsedErr error
	saveCalls     int
	savedImportID string
	savedReceipts []domain.Receipt
}

func (f *fakeReceiptStore) SaveParsed(ctx context.Context, importID string, receipts []domain.Receipt) error {
	f.saveCalls++
	f.savedImportID = importID
	f.savedReceipts = receipts
	return f.saveParsedErr
}

func (f *fakeReceiptStore) ReadByImport(ctx context.Context, importID string) ([]domain.Receipt, error) {
	return nil, nil
}

type fakeParser struct {
	parseErr error
	receipts []domain.Receipt
	parseCalls int
}

func (f *fakeParser) Parse(ctx context.Context, importID string, raw io.Reader) ([]domain.Receipt, error) {
	f.parseCalls++
	if f.parseErr != nil {
		return nil, f.parseErr
	}
	return f.receipts, nil
}
type fakeLogger struct{}

func (l *fakeLogger) Info(msg string, kv ...any)  {}
func (l *fakeLogger) Warn(msg string, kv ...any)  {}
func (l *fakeLogger) Error(msg string, kv ...any) {}

func newImportServiceTestDeps() (*fakeImportStore, *fakeReceiptStore, *fakeParser, *fakeLogger, *ImportService) {
	imports := &fakeImportStore{
		createIfAbsentCreated: true, // holder fix with duplicate update
	}
	receipts := &fakeReceiptStore{}
	parser := &fakeParser{}
	logger := &fakeLogger{}
	svc := NewImportService(imports, receipts, parser, logger)
	return imports, receipts, parser, logger, svc
}

func TestRunImport_HappyPath(t *testing.T) {
	imports, receipts, parser, _, svc := newImportServiceTestDeps()
	parser.receipts = []domain.Receipt{
		{Sum: 100, TypeOfPayment: domain.PaymentByCash},
	}

	importID, err := svc.RunImport(context.Background(), 10, bytes.NewReader([]byte("raw data")))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if importID == "" {
		t.Fatal("expected non-empty importID")
	}
	if parser.parseCalls != 1 {
		t.Fatalf("expected parser called once, got %d", parser.parseCalls)
	}
	if receipts.saveCalls != 1 {
		t.Fatalf("expected SaveParsed called once, got %d", receipts.saveCalls)
	}


	if len(imports.setStatusCalls) < 2 {
		t.Fatalf("expected >=2 status calls, got %d", len(imports.setStatusCalls))
	}
	if imports.setStatusCalls[0] != domain.ImportProcessing {
		t.Fatalf("expected first status %q, got %q", domain.ImportProcessing, imports.setStatusCalls[0])
	}
	last := imports.setStatusCalls[len(imports.setStatusCalls)-1]
	if last != domain.ImportDone {
		t.Fatalf("expected final status %q, got %q", domain.ImportDone, last)
	}
	if imports.getByDocumentIDCalls != 0 {
		t.Fatalf("expected no duplicate-check calls, got %d", imports.getByDocumentIDCalls)
	}
}

func TestRunImport_ParseError_SetsFailed(t *testing.T) {
	imports, receipts, parser, _, svc := newImportServiceTestDeps()
	parser.parseErr = errors.New("parse boom")

	_, err := svc.RunImport(context.Background(), 10, bytes.NewReader([]byte("raw data")))
	if err == nil {
		t.Fatal("expected error, got nil")
	}


	foundFailed := false
	for _, st := range imports.setStatusCalls {
		if st == domain.ImportFailed {
			foundFailed = true
			break
		}
	}
	if !foundFailed {
		t.Fatalf("expected ImportFailed status attempt, got calls: %#v", imports.setStatusCalls)
	}

	if receipts.saveCalls != 0 {
		t.Fatalf("expected SaveParsed not called, got %d", receipts.saveCalls)
	}
	if imports.getByDocumentIDCalls != 0 {
		t.Fatalf("expected no duplicate-check calls, got %d", imports.getByDocumentIDCalls)
	}
}

func TestRunImport_SaveParsedError_SetsFailed(t *testing.T) {
	imports, receipts, parser, _, svc := newImportServiceTestDeps()

	parser.receipts = []domain.Receipt{
		{Sum: 100, TypeOfPayment: domain.PaymentByCash},
	}
	receipts.saveParsedErr = errors.New("save failed")

	_, err := svc.RunImport(context.Background(), 10, bytes.NewReader([]byte("raw data")))
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	foundFailed := false
	for _, st := range imports.setStatusCalls {
		if st == domain.ImportFailed {
			foundFailed = true
			break
		}
	}
	if !foundFailed {
		t.Fatalf("expected ImportFailed status attempt, got calls: %#v", imports.setStatusCalls)
	}

	if receipts.saveCalls != 1 {
		t.Fatalf("expected SaveParsed called once, got %d", receipts.saveCalls)
	}
}

func TestRunImport_SetProcessingError(t *testing.T) {
	imports, _, parser, _, svc := newImportServiceTestDeps()

	imports.setStatusErrByStatus = map[domain.ImportStatus]error{
		domain.ImportProcessing: errors.New("set processing failed"),
	}

	_, err := svc.RunImport(context.Background(), 10, bytes.NewReader([]byte("raw data")))
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if len(imports.setStatusCalls) == 0 || imports.setStatusCalls[0] != domain.ImportProcessing {
		t.Fatalf("expected first status call to be %q, got %#v", domain.ImportProcessing, imports.setStatusCalls)
	}

	if parser.parseCalls != 0 {
		t.Fatalf("expected parser not called, got %d", parser.parseCalls)
	}
}

func TestGetImportStatus_HappyPath(t *testing.T) {
	imports, _, _, _, svc := newImportServiceTestDeps()
	imports.getStatusByIDStatus = domain.ImportDone

	st, err := svc.GetImportStatus(context.Background(), "imp-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if st != domain.ImportDone {
		t.Fatalf("expected status %q, got %q", domain.ImportDone, st)
	}
}

func TestGetImportStatus_StoreError(t *testing.T) {
	imports, _, _, _, svc := newImportServiceTestDeps()
	imports.getStatusByIDErr = errors.New("db down")

	_, err := svc.GetImportStatus(context.Background(), "imp-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}