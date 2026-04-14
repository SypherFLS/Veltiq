package domain

import (
	"fmt"
)

type ErrorCode string // типизация ошибок, разширить

const (
	Invalid_Input ErrorCode = "invalid input"
	Parse_Failed ErrorCode = "parse failed"
	Store_Failed ErrorCode = "store failed"
	Timeout ErrorCode = "timeout"
)

type Rang string

const (
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

