package domain

import (
	"time"
)

type ImportStatus string

const (
	ImportPending    ImportStatus = "pending"
	ImportProcessing ImportStatus = "processing"
	ImportDone       ImportStatus = "done"
	ImportFailed     ImportStatus = "failed"
)


type Import struct {
	ID int
	TenantID int
	Status ImportStatus
	ErrorCode string // TODO: заменить на структурированные ошибки
	UpdatedAt time.Time
	CreatedAt time.Time
}
