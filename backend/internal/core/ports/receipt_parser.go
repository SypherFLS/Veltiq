package ports

import (
	"context"
	"io"
	"veltiq/internal/core/domain"
)

type ParsedImport struct {
	Receipts []domain.Receipt
	Items    []domain.ReceiptItem
}

type ReceiptParser interface {
	Parse(ctx context.Context, importID string, raw io.Reader) (ParsedImport, error)
}
