package service

import (
	"context"
	_ "fmt"
	"io"
	_ "time"
	"veltiq/internal/core/domain"
	"veltiq/internal/core/ports"
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

func (s *ImportService) RunImport(ctx context.Context, tenantdID int, raw io.Reader) (string, error) {
	if raw == nil {
		var err = &domain.Err {
			ModuleName : "Import",
			ProcessName : "Reading",
			Code : domain.Invalid_Input,
			Retryable: false,
		}
		_ = err
		return "", err
	}

	return "", nil // placeholder nt fix
}

