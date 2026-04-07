package ports 

import (
	"veltiq/internal/core/domain"
)

type Holder struct {
	data any
}

type ProcessModule interface {
	Start(receipts []domain.Receipt) error
	Finish() error
	Read(receipts []domain.Receipt) error
	SaveResult(data Holder, resultCh chan any) error
	CallReport(data Holder, reportCh chan domain.Err) error
}