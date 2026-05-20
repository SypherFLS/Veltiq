package httperrors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"veltiq/internal/core/domain"
)

func Write(c *gin.Context, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, domain.ErrEmailTaken):
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
	case errors.Is(err, domain.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
	case errors.Is(err, domain.ErrImportNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "import not found"})
	case errors.Is(err, domain.ErrImportForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "import forbidden"})
	case errors.Is(err, domain.ErrImportNotReady):
		c.JSON(http.StatusConflict, gin.H{"error": "import not ready"})
	default:
		var dup *domain.ImportDuplicateError
		if errors.As(err, &dup) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "import already exists",
				"import_id": dup.ExistingImportID,
			})
			return
		}
		if errors.Is(err, domain.ErrImportDuplicate) {
			c.JSON(http.StatusConflict, gin.H{"error": "import already exists"})
			return
		}

		var derr *domain.Err
		if errors.As(err, &derr) {
			switch derr.Code {
			case domain.Invalid_Input, domain.Parse_Failed:
				c.JSON(http.StatusBadRequest, gin.H{"error": derr.ErrorMessage})
			case domain.Canceled, domain.Timeout:
				c.JSON(http.StatusRequestTimeout, gin.H{"error": derr.ErrorMessage})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": derr.ErrorMessage})
			}
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
