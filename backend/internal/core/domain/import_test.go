package domain

import (
	"time"
	"testing"
)

func TestNewImport_Service(t *testing.T) {
	now := time.Date(2026, time.May, 11, 12, 0, 0, 0, time.UTC)
	payload := []byte("receipt-raw-data")

	imp := NewImport(42, payload, now)

	if imp == nil {
		t.Fatal("imp nil bruh")
	}

	if imp.ID == "" {
		t.Fatal("empty ID")
	}

	if imp.TenantID != 42 {
		t.Fatalf("expected tenantID=42, got %d", imp.TenantID)
	}

	if imp.Status != ImportPending {
		t.Fatalf("expected status=%q, got %q", ImportPending, imp.Status)
	}

	if !imp.CreatedAt.Equal(now) || !imp.UpdatedAt.Equal(now) {
		t.Fatalf("expected timestamps to equal now=%v", now)
	}

	if imp.DocumentID == "" {
		t.Fatal("expected non-empty document ID")
	}
}

func TestBuildDocumentID_Deterministic(t *testing.T) {
	imp := &Import{TenantID: 1}
	payload := []byte("same-payload")
	id1 := imp.BuildDocumentID(payload)
	id2 := imp.BuildDocumentID(payload)
	if id1 == "" || id2 == "" {
		t.Fatal("expected non-empty hash")
	}
	if id1 != id2 {
		t.Fatalf("expected same hashes, got %q and %q", id1, id2)
	}
}
func TestBuildDocumentID_DifferentByTenant(t *testing.T) {
	payload := []byte("same-payload")
	imp1 := &Import{TenantID: 1}
	imp2 := &Import{TenantID: 2}
	id1 := imp1.BuildDocumentID(payload)
	id2 := imp2.BuildDocumentID(payload)
	if id1 == id2 {
		t.Fatalf("expected different hashes for different tenants, got %q", id1)
	}
}
func TestGetImportStatus_ReturnsCurrentStatus(t *testing.T) {
	imp := &Import{Status: ImportDone}
	got := imp.GetImportStatus()
	if got != ImportDone {
		t.Fatalf("expected %q, got %q", ImportDone, got)
	}
}