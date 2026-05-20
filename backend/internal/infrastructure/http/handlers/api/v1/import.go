package v1

import (
	"errors"
	"io"
	"net/http"
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
		"import_id": importID,
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

	c.JSON(http.StatusOK, gin.H{
		"import_id": imp.ID,
		"status": imp.Status,
	})
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
	ImportID string `json:"import_id"`
	Status domain.ImportStatus `json:"status"`
	Data domain.ReportSummary `json:"data"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func reportResponse(r domain.Report) reportJSON {
	return reportJSON{
		ImportID: r.ImportID,
		Status: r.Status,
		Data: r.Data,
		CreatedAt: r.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: r.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
