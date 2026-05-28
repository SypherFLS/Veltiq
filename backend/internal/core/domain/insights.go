package domain

import "time"

type RecommendationKind string

const (
	RecommendationDiscount RecommendationKind = "discount"
	RecommendationBundle   RecommendationKind = "bundle"
	RecommendationWriteoff RecommendationKind = "writeoff"
	RecommendationMonitor  RecommendationKind = "monitor"
)

type IlliquidItem struct {
	SKU                string
	Name               string
	Category           string
	SalesQuantity      int
	DaysWithoutSale    int
	LastSaleAt         time.Time
	Recommendation     RecommendationKind
	RecommendationNote string
}

type Insights struct {
	ImportID    string
	GeneratedAt time.Time
	Items       []IlliquidItem
}
