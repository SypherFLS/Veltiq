package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestInitError_SetsAllFields(t *testing.T) {
	cause := errors.New("db error")

	err := InitError("Import", "Save", Store_Failed, cause, "save failed", Local, true)

	if err.ModuleName != "Import" {
		t.Fatalf("expected ModuleName=Import, got %s", err.ModuleName)
	}
	if err.ProcessName != "Save" {
		t.Fatalf("expected ProcessName=Save, got %s", err.ProcessName)
	}
	if err.Code != Store_Failed {
		t.Fatalf("expected code=%q, got %q", Store_Failed, err.Code)
	}
	if err.Cause != cause {
		t.Fatal("expected cause to match input cause")
	}
	if err.ErrorMessage != "save failed" {
		t.Fatalf("expected error message save failed, got %s", err.ErrorMessage)
	}
	if err.Severity != Local {
		t.Fatalf("expected severity=%q, got %q", Local, err.Severity)
	}
	if !err.Retryable {
		t.Fatal("expected retryable=true")
	}
}

func TestErr_Error_ContainsMainParts(t *testing.T) {
	err := InitError("Import", "Parsing", Parse_Failed, nil, "parse failed", Local, false)

	msg := err.Error()

	if !strings.Contains(msg, "module: Import") {
		t.Fatalf("unexpected msg: %s", msg)
	}
	if !strings.Contains(msg, "process: Parsing") {
		t.Fatalf("unexpected msg: %s", msg)
	}
	if !strings.Contains(msg, "message: parse failed") {
		t.Fatalf("unexpected msg: %s", msg)
	}
}