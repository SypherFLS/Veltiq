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

func (s *ImportService) GetImportStatus(ctx context.Context, importID string) (domain.ImportStatus, error) {
	if ctx == nil {
		return domain.ImportPending, domain.InitError("Import", "GetStatus", domain.Invalid_Input, nil, "empty ctx", domain.Local, false,
		)
	}
	if err := ctx.Err(); err != nil {
		if err == context.Canceled {
			return domain.ImportPending, domain.InitError("Import", "GetStatus", domain.Canceled, err, "ctx canceled", domain.Warn, false,
			)
		}
		if err == context.DeadlineExceeded {
			return domain.ImportPending, domain.InitError("Import", "GetStatus", domain.Timeout, err, "timeout", domain.Warn, true,
			)
		}
	}

	status, err := s.imports.GetStatusByID(ctx, importID)

	if err != nil {
		wrapped := s.wrapErr("GetStatus", domain.Import_Failed, err, err.Error(), domain.Local, false)
		s.logger.Error("get_status_failed", "importID", importID, "err", wrapped.Error())
		return domain.ImportPending, wrapped 
	}
	
	return status, nil
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
		wrapped := s.wrapErr("CreateImport", domain.Import_Failed, err, err.Error(), domain.Local, false)
		s.logger.Error("create_import_failed", "tenantID", tenantID, "err", wrapped.Error())
		return "", wrapped
	}

	if created {
		wrapped := s.wrapErr("CreateImport", domain.Import_Failed, err, "already exist", domain.Local, false)
		s.logger.Error("create_import_failed", "tenantID", tenantID, "err", wrapped.Error())
		return "", wrapped
	} // переработать напрочь пока заглушка ошибки

	

	if err := s.imports.SetStatus(ctx, imp.ID, domain.ImportProcessing); err != nil {
		wrapped := s.wrapErr("SetProcessing", domain.Status_Update_Failed, err, err.Error(), domain.Local, true)
		s.trySetFailed(ctx, imp.ID, wrapped)
		s.logger.Error("set_processing_failed", "importID", imp.ID, "err", wrapped.Error())
		return "", wrapped
	}

	receipts, err := s.parser.Parse(ctx, imp.ID, bytes.NewReader(payload))

	if err != nil {
		wrapped := s.wrapErr("Parsing", domain.Parse_Failed, err, err.Error(), domain.Local, false)
		s.trySetFailed(ctx, imp.ID, wrapped)
		s.logger.Error("parse_failed", "importID", imp.ID, "tenantID", tenantID, "err", wrapped.Error())
		return "", wrapped
	}

	if err := s.receipts.SaveParsed(ctx, imp.ID, receipts); err != nil {
		wrapped := s.wrapErr("SaveParsed", domain.Store_Failed, err, err.Error(), domain.Local, true)
		s.trySetFailed(ctx, imp.ID, wrapped)
		s.logger.Error("save_parsed_failed", "importID", imp.ID, "err", wrapped.Error())
		return "", wrapped
	}


	if err := s.imports.SetStatus(ctx, imp.ID, domain.ImportDone); err != nil {
		wrapped := s.wrapErr("SetDone", domain.Status_Update_Failed, err, err.Error(), domain.Local, true)
		s.trySetFailed(ctx, imp.ID, wrapped)
		s.logger.Error("set_done_failed", "importID", imp.ID, "err", wrapped.Error())
		return "", wrapped
	}
	
	return imp.ID, nil
} // доработать, минимально работоспособность обеспечена

func (s *ImportService) trySetFailed(ctx context.Context, importID string, cause error) {
	if importID == "" {
		return
	}
	if err := s.imports.SetStatus(ctx, importID, domain.ImportFailed); err != nil {
		s.logger.Error(
			"set_failed_status_error",
			"importID", importID,
			"cause", cause.Error(),
			"statusErr", err.Error(),
		)
	}
}

func (s *ImportService) wrapErr(
	process string,
	code domain.ErrorCode,
	cause error,
	msg string,
	severity domain.Rang,
	retryable bool,
) error {
	return domain.InitError("Import", process, code, cause, msg, severity, retryable)
}

func ValidateImport(ctx context.Context, tenantID int, raw io.Reader) error {
	if raw == nil {
		err := domain.InitError ("Import", "Reading",  domain.Invalid_Input, nil, "empty input", domain.Local, false)
		return err
	}
	if tenantID <= 0 {
		err := domain.InitError ("Import", "Reading", domain.Invalid_Input, nil, "invalid tenantID", domain.Local, false)
		return err
	}
	if ctx == nil {
		err := domain.InitError ("Import", "Reading", domain.Invalid_Input, nil, "empty ctx", domain.Local, false)
		return err
	}
	if ctx.Err() == context.Canceled {
		err := domain.InitError ("Import", "Reading", domain.Canceled, nil, "ctx canceled", domain.Warn, false)
		return err
	}
	if ctx.Err() == context.DeadlineExceeded {
		err := domain.InitError ("Import", "Reading", domain.Timeout, nil, "timeout", domain.Warn, true)
		return err
	}

	return nil
}

