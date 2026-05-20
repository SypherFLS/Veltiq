package parser

import (
	"context"
	"io"

	"veltiq/internal/core/domain"
)

type StubParser struct{}

func NewStub() *StubParser {
	return &StubParser{}
}

func (p *StubParser) Parse(ctx context.Context, importID string, raw io.Reader) ([]domain.Receipt, error) {
	_ = ctx
	_, _ = io.ReadAll(raw)

	return []domain.Receipt{
		{
			ImportID: importID,
			TenantID: "",
			StoreID: 1,
			TypeOfPayment: domain.PaymentByCard,
			Sum: 1000,
			Cashier: "stub",
		},
		{
			ImportID: importID,
			TenantID: "",
			StoreID: 2,
			TypeOfPayment: domain.PaymentByCash,
			Sum: 500,
			Cashier: "stub",
		},
	}, nil
}
