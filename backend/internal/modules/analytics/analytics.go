package analytics

import (
	"context"
	"sort"
	"time"

	"veltiq/internal/core/domain"
)

const (
	// thresholdMonitor — товар считается «следить» от этого числа дней без продаж.
	thresholdMonitor = 14
	// thresholdDiscount — рекомендуем распродажу со скидкой.
	thresholdDiscount = 30
	// thresholdWriteoff — рекомендуем списать.
	thresholdWriteoff = 60

	maxIlliquidItems = 100
)

type Analyzer struct{}

func New() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(ctx context.Context, importID string, receipts []domain.Receipt) (domain.ReportSummary, error) {
	_ = ctx
	_ = importID

	var totalSum, cashSum, cardSum int
	for _, r := range receipts {
		totalSum += r.Sum
		switch r.TypeOfPayment {
		case domain.PaymentByCash:
			cashSum += r.Sum
		case domain.PaymentByCard:
			cardSum += r.Sum
		}
	}

	return domain.ReportSummary{
		ReceiptsCount: len(receipts),
		TotalSum:      totalSum,
		CashSum:       cashSum,
		CardSum:       cardSum,
		IsStub:        false,
	}, nil
}

type productAgg struct {
	name          string
	category      string
	totalQuantity int
	lastSaleAt    time.Time
}

func (a *Analyzer) BuildInsights(ctx context.Context, importID string, items []domain.ReceiptItem) (domain.Insights, error) {
	_ = ctx

	if len(items) == 0 {
		return domain.Insights{
			ImportID:    importID,
			GeneratedAt: time.Now().UTC(),
			Items:       []domain.IlliquidItem{},
		}, nil
	}

	bySKU := make(map[string]*productAgg, 64)
	var periodEnd time.Time

	for _, it := range items {
		if it.SoldAt.After(periodEnd) {
			periodEnd = it.SoldAt
		}

		agg, ok := bySKU[it.SKU]
		if !ok {
			agg = &productAgg{
				name:     it.Name,
				category: it.Category,
			}
			bySKU[it.SKU] = agg
		}
		agg.totalQuantity += it.Quantity
		if it.SoldAt.After(agg.lastSaleAt) {
			agg.lastSaleAt = it.SoldAt
			if it.Name != "" {
				agg.name = it.Name
			}
			if it.Category != "" {
				agg.category = it.Category
			}
		}
	}

	if periodEnd.IsZero() {
		periodEnd = time.Now().UTC()
	}

	insights := make([]domain.IlliquidItem, 0, len(bySKU))
	for sku, agg := range bySKU {
		days := int(periodEnd.Sub(agg.lastSaleAt).Hours() / 24)
		if days < thresholdMonitor {
			continue
		}

		rec, note := classify(days, agg.totalQuantity)

		insights = append(insights, domain.IlliquidItem{
			SKU:                sku,
			Name:               agg.name,
			Category:           agg.category,
			SalesQuantity:      agg.totalQuantity,
			DaysWithoutSale:    days,
			LastSaleAt:         agg.lastSaleAt,
			Recommendation:     rec,
			RecommendationNote: note,
		})
	}

	sort.Slice(insights, func(i, j int) bool {
		if insights[i].DaysWithoutSale != insights[j].DaysWithoutSale {
			return insights[i].DaysWithoutSale > insights[j].DaysWithoutSale
		}
		return insights[i].SalesQuantity < insights[j].SalesQuantity
	})

	if len(insights) > maxIlliquidItems {
		insights = insights[:maxIlliquidItems]
	}

	return domain.Insights{
		ImportID:    importID,
		GeneratedAt: time.Now().UTC(),
		Items:       insights,
	}, nil
}

func classify(days, quantity int) (domain.RecommendationKind, string) {
	switch {
	case days >= thresholdWriteoff:
		return domain.RecommendationWriteoff, "Списать или вернуть поставщику — продаж нет ≥ 60 дней"
	case days >= thresholdDiscount:
		if quantity == 0 {
			return domain.RecommendationWriteoff, "Совсем не продаётся — рассмотреть списание"
		}
		return domain.RecommendationDiscount, "Распродажа со скидкой 15–30%"
	default:
		return domain.RecommendationMonitor, "Следить за динамикой — продаж нет ≥ 14 дней"
	}
}
