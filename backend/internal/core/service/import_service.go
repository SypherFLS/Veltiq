package service

import (
	"context"
	_ "fmt"
	"io"
	_ "time"
	"veltiq/internal/core/domain"
	"veltiq/internal/core/ports"
	"time"
)
type ImportService struct {
	imports  ports.ImportStore
	receipts ports.ReceiptStore
	parser   ports.ReceiptParser
	logger   ports.Logger
}

func NewImportService(imports ports.ImportStore, receipts ports.ReceiptStore, parser ports.ReceiptParser, logger ports.Logger,) *ImportService {
	return &ImportService{
		imports:  imports,
		receipts: receipts,
		parser:   parser,
		logger:   logger,
	}
}

func (s *ImportService) RunImport(ctx context.Context, tenantID int, raw io.Reader) (string, error)  {
	if err := ValidateImport(ctx, tenantID, raw);  err != nil {
		return "", err
	}

	payload, err := io.ReadAll(raw)
	if err != nil {
		return "", domain.InitError("Import", "Reading", domain.Parse_Failed, err.Error(), domain.Local, true)
	}

	imp := domain.NewImport(tenantID, payload, time.Now().UTC())

	created, err := s.imports.CreateIfAbsent(ctx, *imp)

	_ = created
	
	return "", nil // placeholder nt fix
}

func ValidateImport(ctx context.Context, tenantID int, raw io.Reader) error {
	if raw == nil {
		err := domain.InitError ("Import", "Reading", domain.Invalid_Input, "empty input", domain.Local, false)
		return err
	}
	if tenantID <= 0 {
		err := domain.InitError ("Import", "Reading", domain.Invalid_Input, "invalid tenantID", domain.Local, false)
		return err
	}
	if ctx == nil {
		err := domain.InitError ("Import", "Reading", domain.Invalid_Input, "empty ctx", domain.Local, false)
		return err
	}
	if ctx.Err() == context.Canceled {
		err := domain.InitError ("Import", "Reading", domain.Canceled, "ctx canceled", domain.Warn, false)
		return err
	}
	if ctx.Err() == context.DeadlineExceeded {
		err := domain.InitError ("Import", "Reading", domain.Timeout, "timeout", domain.Warn, true)
		return err
	}

	return nil
}
