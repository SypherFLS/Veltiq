package domain

import "time"

type PaymentType string

const (
	PaymentByCard PaymentType = "card"
	PaymentByCash PaymentType = "cash"
)

type Product struct {
	Name string
	Quantity int
	UnitPrice int
	TotalPrice int
}

func (p *Product) TotalByOne() {
	p.TotalPrice = p.UnitPrice * p.Quantity
}

func (p *Product) OneByTotal() {
	if p.Quantity > 0 {
		p.UnitPrice = p.TotalPrice / p.Quantity
	}
}

type Receipt struct {
	ID int
	ImportID string
	TenantID string
	StoreID int
	TypeOfPayment PaymentType
	Sum int
	Items []Product
	Cashier string
	CreatedAt time.Time
}
