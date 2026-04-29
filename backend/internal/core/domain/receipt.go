package domain 

import (
	"time"
)

type TOP string 

const (
	PaymentByCard TOP = "card"
	PaymentByCash TOP = "cash"
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
	ImportID int 
	TenantID int
	StoreID int
	TypeOfPayment TOP
	Sum int
	Items []Product
	Cashier string // TODO: реализовать подструктуру/хранение кассиров для дальнейшей обработки
	CreatedAt time.Time
}


