package ports

import (
	"veltiq/internal/core/domain"
)

type ImportModule interface {
	SaveParsed(data []domain.Receipt) error
	Read(data ...any) error
	Parse(data ...any) ([]domain.Receipt, error)
	Validate(data ...any) error
	CallRep(data Holder, ReportChan chan domain.Err) error
}

