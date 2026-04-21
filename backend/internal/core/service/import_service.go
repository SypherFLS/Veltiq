package service

import (
	"context"
	"fmt"
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

	s.logger.Info(fmt.Sprintf("importing %v started at %v", tenantID, time.Now()), ctx) // пока хз, может передавать время отдельно

	if err := ValidateImport(ctx, tenantID, raw);  err != nil {
		s.logger.Error(err.Error(), ctx)
		return "", err
	}

	payload, err := io.ReadAll(raw)

	if err != nil {
		osh := domain.InitError("Import", "Reading", domain.Parse_Failed, err, err.Error(), domain.Local, false)
		s.logger.Error(osh.Error(), ctx)
		return "", osh
	}

	imp := domain.NewImport(tenantID, payload, time.Now().UTC())

	created, err := s.imports.CreateIfAbsent(ctx, *imp)
	
	if err != nil {
		osh := domain.InitError("Import", "creating new import", domain.Import_Failed,  err, err.Error(), domain.Local, false)
		s.logger.Error(osh.Error(), ctx)
		return "", osh
	}

	if created {

	}

	// проверку на дупликат

	if err = s.imports.SetStatus(ctx, imp.ID, domain.ImportProcessing); err != nil {
		osh := domain.InitError("Import", "Status updaiting failed", domain.Status_Update_Failed, err, err.Error(), domain.Local, true)
		s.logger.Error(osh.Error(), ctx)
		return "", osh
	}

	receipts, err := s.parser.Parse(ctx, imp.ID, bytes.NewReader(payload))

	s.receipts.SaveParsed(ctx, imp.ID, receipts)

	s.imports.SetStatus(ctx, imp.ID, domain.ImportDone)
	
	return imp.ID, nil
} // без логов, обработки внутренних ошибок проверок на дубликат и тд, доработать, минимально работоспособность обеспечена

func ValidateImport(ctx context.Context, tenantID int, raw io.Reader) error {
	var osh error 
	if raw == nil {
		err := domain.InitError ("Import", "Reading",  domain.Invalid_Input, osh, "empty input", domain.Local, false)
		return err
	}
	if tenantID <= 0 {
		err := domain.InitError ("Import", "Reading", domain.Invalid_Input, osh, "invalid tenantID", domain.Local, false)
		return err
	}
	if ctx == nil {
		err := domain.InitError ("Import", "Reading", domain.Invalid_Input, osh, "empty ctx", domain.Local, false)
		return err
	}
	if ctx.Err() == context.Canceled {
		err := domain.InitError ("Import", "Reading", domain.Canceled, osh, "ctx canceled", domain.Warn, false)
		return err
	}
	if ctx.Err() == context.DeadlineExceeded {
		err := domain.InitError ("Import", "Reading", domain.Timeout, osh, "timeout", domain.Warn, true)
		return err
	}

	return nil
}
