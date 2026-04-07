package domain

import (
	"fmt"
)

type Err struct {
	ModuleName string
	ProcessName string
	ErrorMessage string
}

func (e *Err) Error() string {
	return fmt.Sprintf("error: module: %s, process: %s, message: %s", e.ModuleName, e.ProcessName, e.ErrorMessage)
}

