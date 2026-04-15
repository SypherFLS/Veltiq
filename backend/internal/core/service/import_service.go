package service

import (
	"context"
	_ "fmt"
	"io"
	"veltiq/internal/core/domain"
	"veltiq/internal/core/ports"
	"bytes"
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

	// проверку на дупликат

	s.imports.SetStatus(ctx, imp.ID, domain.ImportProcessing)

	receipts, err := s.parser.Parse(ctx, imp.ID, bytes.NewReader(payload))

	s.receipts.SaveParsed(ctx, imp.ID, receipts)

	s.imports.SetStatus(ctx, imp.ID, domain.ImportDone)
	
	return imp.ID, nil
} // без логов, обработки внутренних ошибок проверок на дубликат и тд, доработать, минимально работоспособность обеспечена

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
