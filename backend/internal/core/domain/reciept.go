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
	PricePerOne int
	TotalPrice int
}

func (p *Product) TotalByOne() {
	p.TotalPrice = p.PricePerOne * p.Quantity
}

func (p *Product) OneByTotal() {
	p.PricePerOne  = p.TotalPrice / p.Quantity
}

type Receipt struct {
	ID       int
	ImportID int 
	TenantID int
	StoreID int
	TypeOfPayment TOP
	Summ float64
	Items []Product
	Cashier string // реализовать подструктуру/хранение кассиров для дальнейшей обработки
	Time time.Time
}


