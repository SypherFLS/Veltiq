package v1

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"veltiq/internal/core/domain"
	"veltiq/internal/core/orchestrator"
	"veltiq/internal/infrastructure/http/httperrors"
	"veltiq/internal/infrastructure/http/middleware"
)

type ImportHandler struct {
	runner *orchestrator.Runner
}

func NewImportHandler(runner *orchestrator.Runner) *ImportHandler {
	return &ImportHandler{runner: runner}
}

type importItemJSON struct {
	ID        string              `json:"id"`
	Status    domain.ImportStatus `json:"status"`
	ErrorCode string              `json:"errorCode,omitempty"`
	CreatedAt string              `json:"createdAt"`
	UpdatedAt string              `json:"updatedAt"`
}

type importListJSON struct {
	Items  []importItemJSON `json:"items"`
	Total  int64            `json:"total"`
	Cursor *string          `json:"cursor"`
}

func importToJSON(i domain.Import) importItemJSON {
	out := importItemJSON{
		ID:        i.ID,
		Status:    i.Status,
		ErrorCode: i.ErrorCode,
	}
	if !i.CreatedAt.IsZero() {
		out.CreatedAt = i.CreatedAt.UTC().Format(time.RFC3339)
	}
	if !i.UpdatedAt.IsZero() {
		out.UpdatedAt = i.UpdatedAt.UTC().Format(time.RFC3339)
	}
	return out
}

func (h *ImportHandler) List(c *gin.Context) {
	tenantID, ok := tenantFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit := 20
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}

	imports, total, err := h.runner.ListImports(c.Request.Context(), tenantID, limit)
	if err != nil {
		httperrors.Write(c, err)
		return
	}

	items := make([]importItemJSON, 0, len(imports))
	for _, imp := range imports {
		items = append(items, importToJSON(imp))
	}

	c.JSON(http.StatusOK, importListJSON{
		Items:  items,
		Total:  total,
		Cursor: nil,
	})
}

func (h *ImportHandler) Upload(c *gin.Context) {
	tenantID, ok := tenantFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	body, err := readCSVBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer body.Close()

	importID, err := h.runner.StartImport(c.Request.Context(), tenantID, body)
	if err != nil {
		httperrors.Write(c, err)
		return
	}

	status, err := h.runner.GetImportStatus(c.Request.Context(), importID, tenantID)
	if err != nil {
		httperrors.Write(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":     importID,
		"status": status,
	})
}

func (h *ImportHandler) Status(c *gin.Context) {
	tenantID, ok := tenantFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	importID := c.Param("id")
	imp, err := h.runner.GetImport(c.Request.Context(), importID, tenantID)
	if err != nil {
		httperrors.Write(c, err)
		return
	}

	c.JSON(http.StatusOK, importToJSON(imp))
}

func (h *ImportHandler) Report(c *gin.Context) {
	tenantID, ok := tenantFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	importID := c.Param("id")
	report, err := h.runner.GetImportReport(c.Request.Context(), importID, tenantID)
	if err != nil {
		httperrors.Write(c, err)
		return
	}

	c.JSON(http.StatusOK, reportResponse(report))
}

type illiquidItemJSON struct {
	SKU                string `json:"sku"`
	Name               string `json:"name"`
	Category           string `json:"category,omitempty"`
	SalesQuantity      int    `json:"salesQuantity"`
	DaysWithoutSale    int    `json:"daysWithoutSale"`
	LastSaleAt         string `json:"lastSaleAt,omitempty"`
	Recommendation     string `json:"recommendation,omitempty"`
	RecommendationNote string `json:"recommendationNote,omitempty"`
}

type insightsJSON struct {
	ImportID    string             `json:"importId"`
	GeneratedAt string             `json:"generatedAt"`
	Items       []illiquidItemJSON `json:"items"`
}

func (h *ImportHandler) Insights(c *gin.Context) {
	tenantID, ok := tenantFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	importID := c.Param("id")
	insights, err := h.runner.GetImportInsights(c.Request.Context(), importID, tenantID)
	if err != nil {
		httperrors.Write(c, err)
		return
	}

	items := make([]illiquidItemJSON, 0, len(insights.Items))
	for _, it := range insights.Items {
		lastSale := ""
		if !it.LastSaleAt.IsZero() {
			lastSale = it.LastSaleAt.UTC().Format(time.RFC3339)
		}
		items = append(items, illiquidItemJSON{
			SKU:                it.SKU,
			Name:               it.Name,
			Category:           it.Category,
			SalesQuantity:      it.SalesQuantity,
			DaysWithoutSale:    it.DaysWithoutSale,
			LastSaleAt:         lastSale,
			Recommendation:     string(it.Recommendation),
			RecommendationNote: it.RecommendationNote,
		})
	}

	c.JSON(http.StatusOK, insightsJSON{
		ImportID:    insights.ImportID,
		GeneratedAt: insights.GeneratedAt.UTC().Format(time.RFC3339),
		Items:       items,
	})
}

func tenantFromContext(c *gin.Context) (string, bool) {
	v, ok := c.Get(middleware.TenantIDKey)
	if !ok {
		return "", false
	}
	tenantID, ok := v.(string)
	return tenantID, ok && tenantID != ""
}

func readCSVBody(c *gin.Context) (io.ReadCloser, error) {
	file, err := c.FormFile("file")
	if err == nil {
		f, openErr := file.Open()
		if openErr != nil {
			return nil, openErr
		}
		return f, nil
	}

	if c.Request.Body == nil || c.Request.ContentLength == 0 {
		return nil, errors.New("empty csv body")
	}
	return io.NopCloser(c.Request.Body), nil
}

type reportJSON struct {
	ImportID  string               `json:"importId"`
	Status    domain.ImportStatus  `json:"status"`
	Data      domain.ReportSummary `json:"data"`
	CreatedAt string               `json:"createdAt"`
	UpdatedAt string               `json:"updatedAt"`
	ErrorCode string               `json:"errorCode,omitempty"`
}

func reportResponse(r domain.Report) reportJSON {
	return reportJSON{
		ImportID:  r.ImportID,
		Status:    r.Status,
		Data:      r.Data,
		CreatedAt: r.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: r.UpdatedAt.UTC().Format(time.RFC3339),
		ErrorCode: r.ErrorCode,
	}
}
