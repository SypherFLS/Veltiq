package parser

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"veltiq/internal/core/domain"
	"veltiq/internal/core/ports"
)

// CSVParser parses a flat CSV with one product per line.
//
// Expected columns (case-insensitive, order from the header):
//
//	receipt_id,date,store_id,payment,sku,product_name,category,quantity,unit_price
//
// Where:
//   - receipt_id: any unique identifier of the receipt (string)
//   - date: RFC3339 / "2006-01-02 15:04:05" / "2006-01-02 15:04" / "2006-01-02"
//   - store_id: integer
//   - payment: "card" | "cash" (case-insensitive, "карта"/"наличные" тоже принимаются)
//   - sku: string
//   - product_name: string
//   - category: string (optional, can be empty)
//   - quantity: integer
//   - unit_price: integer (price in main currency units, e.g. ₽)
type CSVParser struct{}

func NewCSV() *CSVParser {
	return &CSVParser{}
}

var requiredColumns = []string{
	"receipt_id",
	"date",
	"store_id",
	"payment",
	"sku",
	"product_name",
	"category",
	"quantity",
	"unit_price",
}

var dateLayouts = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
	"2006-01-02",
	"02.01.2006 15:04",
	"02.01.2006",
}

func (p *CSVParser) Parse(ctx context.Context, importID string, raw io.Reader) (ports.ParsedImport, error) {
	if raw == nil {
		return ports.ParsedImport{}, errors.New("empty reader")
	}

	reader := csv.NewReader(raw)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1
	reader.ReuseRecord = false

	headerRow, err := reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return ports.ParsedImport{}, errors.New("empty csv")
		}
		return ports.ParsedImport{}, fmt.Errorf("read header: %w", err)
	}

	headers := make(map[string]int, len(headerRow))
	for i, h := range headerRow {
		headers[strings.ToLower(strings.TrimSpace(h))] = i
	}
	for _, col := range requiredColumns {
		if _, ok := headers[col]; !ok {
			return ports.ParsedImport{}, fmt.Errorf("missing required column: %s", col)
		}
	}

	idx := func(name string) int { return headers[name] }

	type receiptKey struct {
		external string
	}

	type receiptAgg struct {
		receipt domain.Receipt
		seen    bool
	}

	receiptsByKey := make(map[receiptKey]*receiptAgg, 64)
	items := make([]domain.ReceiptItem, 0, 128)

	lineNumber := 1
	for {
		if err := ctx.Err(); err != nil {
			return ports.ParsedImport{}, err
		}

		row, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return ports.ParsedImport{}, fmt.Errorf("row %d: %w", lineNumber+1, err)
		}
		lineNumber++

		if isBlankRow(row) {
			continue
		}

		externalID := strings.TrimSpace(safeAt(row, idx("receipt_id")))
		if externalID == "" {
			return ports.ParsedImport{}, fmt.Errorf("row %d: empty receipt_id", lineNumber)
		}

		soldAt, err := parseDate(strings.TrimSpace(safeAt(row, idx("date"))))
		if err != nil {
			return ports.ParsedImport{}, fmt.Errorf("row %d: invalid date: %w", lineNumber, err)
		}

		storeID, err := parseInt(safeAt(row, idx("store_id")))
		if err != nil {
			return ports.ParsedImport{}, fmt.Errorf("row %d: invalid store_id: %w", lineNumber, err)
		}

		paymentType, err := parsePayment(safeAt(row, idx("payment")))
		if err != nil {
			return ports.ParsedImport{}, fmt.Errorf("row %d: %w", lineNumber, err)
		}

		sku := strings.TrimSpace(safeAt(row, idx("sku")))
		if sku == "" {
			return ports.ParsedImport{}, fmt.Errorf("row %d: empty sku", lineNumber)
		}

		name := strings.TrimSpace(safeAt(row, idx("product_name")))
		category := strings.TrimSpace(safeAt(row, idx("category")))

		quantity, err := parseInt(safeAt(row, idx("quantity")))
		if err != nil || quantity <= 0 {
			return ports.ParsedImport{}, fmt.Errorf("row %d: invalid quantity", lineNumber)
		}

		unitPrice, err := parseInt(safeAt(row, idx("unit_price")))
		if err != nil || unitPrice < 0 {
			return ports.ParsedImport{}, fmt.Errorf("row %d: invalid unit_price", lineNumber)
		}

		totalPrice := unitPrice * quantity

		key := receiptKey{external: externalID}
		agg, ok := receiptsByKey[key]
		if !ok {
			agg = &receiptAgg{
				receipt: domain.Receipt{
					ImportID:      importID,
					ExternalID:    externalID,
					StoreID:       storeID,
					TypeOfPayment: paymentType,
					CreatedAt:     soldAt,
					Sum:           0,
					Items:         make([]domain.Product, 0, 2),
				},
				seen: true,
			}
			receiptsByKey[key] = agg
		}

		agg.receipt.Sum += totalPrice
		agg.receipt.Items = append(agg.receipt.Items, domain.Product{
			SKU:        sku,
			Name:       name,
			Category:   category,
			Quantity:   quantity,
			UnitPrice:  unitPrice,
			TotalPrice: totalPrice,
		})

		items = append(items, domain.ReceiptItem{
			ImportID:          importID,
			ReceiptExternalID: externalID,
			SKU:               sku,
			Name:              name,
			Category:          category,
			Quantity:          quantity,
			UnitPrice:         unitPrice,
			TotalPrice:        totalPrice,
			PaymentType:       paymentType,
			SoldAt:            soldAt,
		})
	}

	receipts := make([]domain.Receipt, 0, len(receiptsByKey))
	for _, agg := range receiptsByKey {
		receipts = append(receipts, agg.receipt)
	}

	if len(receipts) == 0 {
		return ports.ParsedImport{}, errors.New("no data rows in csv")
	}

	return ports.ParsedImport{
		Receipts: receipts,
		Items:    items,
	}, nil
}

func safeAt(row []string, i int) string {
	if i < 0 || i >= len(row) {
		return ""
	}
	return row[i]
}

func isBlankRow(row []string) bool {
	for _, v := range row {
		if strings.TrimSpace(v) != "" {
			return false
		}
	}
	return true
}

func parseInt(v string) (int, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, errors.New("empty number")
	}
	return strconv.Atoi(v)
}

func parseDate(v string) (time.Time, error) {
	if v == "" {
		return time.Time{}, errors.New("empty date")
	}
	for _, layout := range dateLayouts {
		if t, err := time.Parse(layout, v); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised date format: %q", v)
}

func parsePayment(v string) (domain.PaymentType, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "card", "карта", "безнал", "credit", "debit":
		return domain.PaymentByCard, nil
	case "cash", "наличные", "нал":
		return domain.PaymentByCash, nil
	default:
		return "", fmt.Errorf("unknown payment type: %q", v)
	}
}
