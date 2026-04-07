package ports		

import (
	"veltiq/internal/core/domain"
)

type ReceiptStore interface {
	SaveParsed(receipts []domain.Receipt) error
	Read() ([]domain.Receipt, error)
}