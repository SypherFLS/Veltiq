package domain

import (
	"time"
	"github.com/google/uuid"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
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
	DocumentID string 
	TenantID int
	Status ImportStatus
	ErrorCode string 
	UpdatedAt time.Time
	CreatedAt time.Time
}

func NewImport(tenantID int, payload []byte, now time.Time) *Import {
	imp := &Import{
		ID:        uuid.New().String(),
		TenantID:  tenantID,
		Status:    ImportPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	imp.DocumentID = imp.BuildDocumentID(payload)
	return imp
}

func (i *Import) Update() { // обновление импорта

}


func (i *Import) BuildDocumentID(payload []byte) string { // уникален под одинаковые документы разных компаний
	h := sha256.New()
	h.Write([]byte(strconv.Itoa(i.TenantID)))
	h.Write([]byte{':'})
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}
