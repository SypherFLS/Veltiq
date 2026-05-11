package testhelpers

import (
	"context"
	"io"

	"veltiq/internal/core/domain"
	"veltiq/internal/core/ports"
)

type FakeImportStore struct {
	CreateIfAbsentCreated bool
	CreateIfAbsentErr     error

	SetStatusErrByStatus map[domain.ImportStatus]error
	SetStatusCalls       []domain.ImportStatus

	GetStatusByIDStatus domain.ImportStatus
	GetStatusByIDErr    error

	GetByDocumentIDResult domain.Import
	GetByDocumentIDErr    error
	GetByDocumentIDCalls  int
}

func (f *FakeImportStore) GetStatusByID(ctx context.Context, id string) (domain.ImportStatus, error) {
	return f.GetStatusByIDStatus, f.GetStatusByIDErr
}

func (f *FakeImportStore) CreateIfAbsent(ctx context.Context, imp domain.Import) (bool, error) {
	return f.CreateIfAbsentCreated, f.CreateIfAbsentErr
}

func (f *FakeImportStore) SetStatus(ctx context.Context, importID string, st domain.ImportStatus) error {
	f.SetStatusCalls = append(f.SetStatusCalls, st)
	if f.SetStatusErrByStatus == nil {
		return nil
	}
	return f.SetStatusErrByStatus[st]
}

func (f *FakeImportStore) GetByDocumentID(ctx context.Context, documentID string) (domain.Import, error) {
	f.GetByDocumentIDCalls++
	if f.GetByDocumentIDErr != nil {
		return domain.Import{}, f.GetByDocumentIDErr
	}
	return f.GetByDocumentIDResult, nil
}

func (f *FakeImportStore) CallReport(ctx context.Context, data ports.Holder, reportCh chan domain.Err) error {
	return nil
}

type FakeReceiptStore struct {
	SaveParsedErr error
	SaveCalls     int
	SavedImportID string
	SavedReceipts []domain.Receipt

	ReadByImportResult []domain.Receipt
	ReadByImportErr    error
}

func (f *FakeReceiptStore) SaveParsed(ctx context.Context, importID string, receipts []domain.Receipt) error {
	f.SaveCalls++
	f.SavedImportID = importID
	f.SavedReceipts = receipts
	return f.SaveParsedErr
}

func (f *FakeReceiptStore) ReadByImport(ctx context.Context, importID string) ([]domain.Receipt, error) {
	if f.ReadByImportErr != nil {
		return nil, f.ReadByImportErr
	}
	return f.ReadByImportResult, nil
}

type FakeParser struct {
	ParseErr   error
	Receipts   []domain.Receipt
	ParseCalls int
}

func (f *FakeParser) Parse(ctx context.Context, importID string, raw io.Reader) ([]domain.Receipt, error) {
	f.ParseCalls++
	if f.ParseErr != nil {
		return nil, f.ParseErr
	}
	return f.Receipts, nil
}

type FakeLogger struct{}

func (l *FakeLogger) Info(msg string, kv ...any)  {}
func (l *FakeLogger) Warn(msg string, kv ...any)  {}
func (l *FakeLogger) Error(msg string, kv ...any) {}

