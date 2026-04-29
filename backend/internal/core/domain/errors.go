package domain

import (
	"fmt"
)

type ErrorCode string // типизация ошибок, расширить

const (
	Invalid_Input ErrorCode = "invalid input"
	Import_Failed ErrorCode = "import failed"
	Status_Update_Failed ErrorCode = "status update failed"
	Parse_Failed ErrorCode = "parse failed"
	Canceled ErrorCode = "canceled"
	Store_Failed ErrorCode = "store failed"
	Timeout ErrorCode = "timeout"
)

type Rang string

const (
	Fatal         Rang = "Fatal"
	Insignificant Rang = "insignificant"
	Local         Rang = "local"
	Warn          Rang = "Warn"
	// прописать ранги значимости
)

type Err struct {
	ModuleName string
	ProcessName string
	Code ErrorCode
	Cause error
	ErrorMessage string
	Severity Rang
	Retryable bool
}

func (e *Err) Error() string {
	return fmt.Sprintf("error: module: %s, process: %s, message: %s", e.ModuleName, e.ProcessName, e.ErrorMessage)
}

func InitError(moduleName string, processName string, code ErrorCode, cause error, errorMessage string, severity Rang, retryable bool) *Err {
	return &Err{
		ModuleName: moduleName,
		ProcessName: processName,
		Code: code,
		Cause: cause,
		ErrorMessage: errorMessage,
		Severity: severity,
		Retryable:retryable,
	}
}
