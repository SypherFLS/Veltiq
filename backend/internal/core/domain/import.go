package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type ImportStatus string

const (
	ImportPending ImportStatus = "pending"
	ImportProcessing ImportStatus = "processing"
	ImportPartialFail ImportStatus = "partial_failed"
	ImportDone ImportStatus = "done"
	ImportFailed ImportStatus = "failed"
)

type Import struct {
	ID string
	DocumentID string
	TenantID string
	Status ImportStatus
	ErrorCode string
	UpdatedAt time.Time
	CreatedAt time.Time
}

func (i *Import) GetImportStatus() ImportStatus {
	return i.Status
}

func NewImport(tenantID string, payload []byte, now time.Time) *Import {
	imp := &Import{
		ID: uuid.New().String(),
		TenantID: tenantID,
		Status: ImportPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	imp.DocumentID = imp.BuildDocumentID(payload)
	return imp
}

func (i *Import) BuildDocumentID(payload []byte) string {
	h := sha256.New()
	h.Write([]byte(i.TenantID))
	h.Write([]byte{':'})
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}
