package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"veltiq/internal/infrastructure/http/middleware"
)

type ProfileHandler struct{}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{}
}

func (h *ProfileHandler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user_id": c.GetString(middleware.UserIDKey),
		"tenant_id": c.GetString(middleware.TenantIDKey),
	})
}
