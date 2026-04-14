package domain

import (
	"time"
)

type ImportStatus string

const (
	ImportPending    ImportStatus = "pending"
	ImportProcessing ImportStatus = "processing"
	ImportPartialFail ImportStatus = "partial_failed"
	ImportDone       ImportStatus = "done"
	ImportFailed     ImportStatus = "failed"
)


type Import struct {
	ID string
	TenantID int
	Status ImportStatus
	ErrorCode string 
	UpdatedAt time.Time
	CreatedAt time.Time
}
