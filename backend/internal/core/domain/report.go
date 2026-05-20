package domain

import "time"

type ReportSummary struct {
	ReceiptsCount int `json:"receiptsCount"`
	TotalSum int `json:"totalSum"`
	CashSum int `json:"cashSum"`
	CardSum int `json:"cardSum"`
	IsStub bool `json:"isStub"`
	Note string `json:"note,omitempty"`
}

type Report struct {
	ImportID string
	Status ImportStatus
	Data ReportSummary
	CreatedAt time.Time
	UpdatedAt time.Time
	ErrorCode string
}
