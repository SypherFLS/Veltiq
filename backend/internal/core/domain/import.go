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
	Id int
	TenantdID int
	Status ImportStatus
	ErrorCode string // поменять на структурированные ошибки
	UpdatedAt time.Time
	Time time.Time
}
