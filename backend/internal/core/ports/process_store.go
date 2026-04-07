package ports 

import (
	"veltiq/internal/core/domain"
)

type Holder struct {
	data any
}

type ProcessModule interface {
	Start(data []domain.Receipt) error 
	Read(data []domain.Receipt) error
	SaveResult(data Holder, res chan any) error 
	CallRep(data Holder, ReportChan chan domain.Err) error
}