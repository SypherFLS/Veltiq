package domain

import (
	"time"
)

type Report struct {
	ImportID      string
	Status        ImportStatus
	Data []any
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ErrorCode     string
}