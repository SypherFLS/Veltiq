package domain

import (
	"fmt"
)

type ErrorCode string // типизация ошибок, разширить

const (
	Invalid_Input ErrorCode = "invalid input"
	Parse_Failed ErrorCode = "parse failed"
	Canceled ErrorCode = "canceled"
	Store_Failed ErrorCode = "store failed"
	Timeout ErrorCode = "timeout"
)

type Rang string

const (
	Fatal Rang = "Fatal"
	Insignificant Rang = "insignificant"
	Local Rang = "local"
	Warn Rang = "Warn"
	// прописать ранги значимости
)

type Err struct {
	ModuleName string
	ProcessName string
	Code ErrorCode
	ErrorMessage string
	Severity Rang
	Retryable bool
}

func (e *Err) Error() string {
	return fmt.Sprintf("error: module: %s, process: %s, message: %s", e.ModuleName, e.ProcessName, e.ErrorMessage)
}

func InitError(moduleName string, processName string, code ErrorCode, errorMessage string, severity Rang, retryable bool) *Err{
	return &Err {
		ModuleName : moduleName,
		ProcessName : processName,
		Code : code,
		ErrorMessage : errorMessage,
		Severity : severity,
		Retryable : retryable,
	}
}