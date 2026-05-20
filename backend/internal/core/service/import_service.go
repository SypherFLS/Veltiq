package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	"veltiq/internal/core/domain"
	"veltiq/internal/core/ports"
)

type ImportService struct {
	imports ports.ImportStore
	workflow ports.ImportWorkflow
	parser ports.ReceiptParser
	logger ports.Logger
}

func NewImportService(
	imports ports.ImportStore,
	workflow ports.ImportWorkflow,
	parser ports.ReceiptParser,
	logger ports.Logger,
) *ImportService {
	return &ImportService{
		imports: imports,
		workflow: workflow,
		parser: parser,
		logger: logger,
	}
}

func (s *ImportService) GetImport(ctx context.Context, importID, tenantID string) (domain.Import, error) {
	imp, err := s.imports.GetByID(ctx, importID)
	if err != nil {
		return domain.Import{}, err
	}
	if imp.TenantID != tenantID {
		return domain.Import{}, domain.ErrImportForbidden
	}
	return imp, nil
}

func (s *ImportService) GetImportStatus(ctx context.Context, importID, tenantID string) (domain.ImportStatus, error) {
	imp, err := s.GetImport(ctx, importID, tenantID)
	if err != nil {
		return "", err
	}
	return imp.Status, nil
}

func (s *ImportService) RunImport(ctx context.Context, tenantID string, raw io.Reader) (string, error) {
	s.logger.Info("import_started", "tenantID", tenantID, "at", time.Now().UTC())

	if err := ValidateImport(ctx, tenantID, raw); err != nil {
		s.logger.Error("import_validation_failed", "err", err.Error())
		return "", err
	}

	payload, err := io.ReadAll(raw)
	if err != nil {
		osh := domain.InitError("Import", "Reading", domain.Parse_Failed, err, err.Error(), domain.SeverityLocal, false)
		s.logger.Error("import_read_failed", "err", osh.Error())
		return "", osh
	}

	imp := domain.NewImport(tenantID, payload, time.Now().UTC())

	exists, err := s.imports.CreateIfAbsent(ctx, *imp)
	if err != nil {
		wrapped := s.wrapErr("CreateImport", domain.Import_Failed, err, err.Error(), domain.SeverityLocal, false)
		s.logger.Error("create_import_failed", "tenantID", tenantID, "err", wrapped.Error())
		return "", wrapped
	}

	if exists {
		existing, docErr := s.imports.GetByDocumentID(ctx, tenantID, imp.DocumentID)
		if docErr != nil {
			s.logger.Error("create_import_duplicate_lookup_failed", "tenantID", tenantID, "err", docErr.Error())
			return "", domain.ErrImportDuplicate
		}
		s.logger.Info("create_import_duplicate", "tenantID", tenantID, "importID", existing.ID)
		return "", &domain.ImportDuplicateError{ExistingImportID: existing.ID}
	}

	if err := s.imports.SetStatus(ctx, imp.ID, domain.ImportProcessing); err != nil {
		wrapped := s.wrapErr("SetProcessing", domain.Status_Update_Failed, err, err.Error(), domain.SeverityLocal, true)
		s.trySetFailed(ctx, imp.ID, wrapped)
		s.logger.Error("set_processing_failed", "importID", imp.ID, "err", wrapped.Error())
		return "", wrapped
	}

	receipts, err := s.parser.Parse(ctx, imp.ID, bytes.NewReader(payload))
	if err != nil {
		wrapped := s.wrapErr("Parsing", domain.Parse_Failed, err, err.Error(), domain.SeverityLocal, false)
		s.trySetFailed(ctx, imp.ID, wrapped)
		s.logger.Error("parse_failed", "importID", imp.ID, "tenantID", tenantID, "err", wrapped.Error())
		return "", wrapped
	}

	for i := range receipts {
		receipts[i].ImportID = imp.ID
		receipts[i].TenantID = tenantID
	}

	if err := s.workflow.PersistResults(ctx, imp.ID, receipts); err != nil {
		wrapped := s.wrapErr("PersistResults", domain.Store_Failed, err, err.Error(), domain.SeverityLocal, true)
		s.trySetFailed(ctx, imp.ID, wrapped)
		s.logger.Error("persist_results_failed", "importID", imp.ID, "err", wrapped.Error())
		return "", wrapped
	}

	s.logger.Info("import_finished", "importID", imp.ID, "tenantID", tenantID)
	return imp.ID, nil
}

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
	severity domain.Severity,
	retryable bool,
) error {
	return domain.InitError("Import", process, code, cause, msg, severity, retryable)
}

func ValidateImport(ctx context.Context, tenantID string, raw io.Reader) error {
	if raw == nil {
		return domain.InitError("Import", "Reading", domain.Invalid_Input, nil, "empty input", domain.SeverityLocal, false)
	}
	if tenantID == "" {
		return domain.InitError("Import", "Reading", domain.Invalid_Input, nil, "invalid tenantID", domain.SeverityLocal, false)
	}
	if ctx == nil {
		return domain.InitError("Import", "Reading", domain.Invalid_Input, nil, "empty ctx", domain.SeverityLocal, false)
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return domain.InitError("Import", "Reading", domain.Canceled, nil, "ctx canceled", domain.SeverityWarn, false)
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return domain.InitError("Import", "Reading", domain.Timeout, nil, "timeout", domain.SeverityWarn, true)
	}
	return nil
}
