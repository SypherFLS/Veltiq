package domain

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	Invalid_Input ErrorCode = "invalid input"
	Import_Failed ErrorCode = "import failed"
	Status_Update_Failed ErrorCode = "status update failed"
	Parse_Failed ErrorCode = "parse failed"
	Canceled ErrorCode = "canceled"
	Store_Failed ErrorCode = "store failed"
	Timeout ErrorCode = "timeout"
)

type Severity string

const (
	SeverityFatal Severity = "fatal"
	SeverityInsignificant Severity = "insignificant"
	SeverityLocal Severity = "local"
	SeverityWarn Severity = "warn"
)

type Err struct {
	ModuleName string
	ProcessName string
	Code ErrorCode
	Cause error
	ErrorMessage string
	Severity Severity
	Retryable bool
}

func (e *Err) Error() string {
	return fmt.Sprintf("error: module: %s, process: %s, message: %s", e.ModuleName, e.ProcessName, e.ErrorMessage)
}

func (e *Err) Unwrap() error {
	return e.Cause
}

func InitError(
	moduleName string,
	processName string,
	code ErrorCode,
	cause error,
	errorMessage string,
	severity Severity,
	retryable bool,
) *Err {
	return &Err{
		ModuleName: moduleName,
		ProcessName: processName,
		Code: code,
		Cause: cause,
		ErrorMessage: errorMessage,
		Severity: severity,
		Retryable: retryable,
	}
}

var (
	ErrEmailTaken = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrImportNotFound = errors.New("import not found")
	ErrImportForbidden = errors.New("import forbidden")
	ErrImportNotReady = errors.New("import not ready")
	ErrImportDuplicate = errors.New("import already exists")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type ImportDuplicateError struct {
	ExistingImportID string
}

func (e *ImportDuplicateError) Error() string {
	return ErrImportDuplicate.Error()
}

func (e *ImportDuplicateError) Is(target error) bool {
	return target == ErrImportDuplicate
}
