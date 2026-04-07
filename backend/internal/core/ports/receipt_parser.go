package ports

import (
	"veltiq/internal/core/domain"
)

type Parser interface {
	Parse(args ...any) ([]domain.Receipt, error)
	CallReport(data Holder, reportCh chan domain.Err) error
}