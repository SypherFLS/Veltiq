package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"veltiq/internal/core/domain"
	"veltiq/internal/core/ports"
	"veltiq/internal/core/service"
	"veltiq/internal/infrastructure/cookies"
	"veltiq/internal/infrastructure/http/httperrors"
	"veltiq/internal/infrastructure/jwt"
)

type AuthHandler struct {
	authService *service.AuthService
	tokens ports.TokenManager
	cookies *cookies.Manager
}

func NewAuthHandler(
	authService *service.AuthService,
	tokens ports.TokenManager,
	cookieManager *cookies.Manager,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		tokens: tokens,
		cookies: cookieManager,
	}
}

type RegisterRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.authService.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		httperrors.Write(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		httperrors.Write(c, err)
		return
	}

	if err := h.setSession(c, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "authenticated"})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := h.cookies.RefreshTokenFromRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}

	userID, err := h.tokens.VerifyRefresh(refreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) || errors.Is(err, jwt.ErrWrongTokenType) {
			h.cookies.ClearSession(c.Writer)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	if err := h.setSession(c, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "refresh failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "refreshed"})
}

func (h *AuthHandler) VerifySession(c *gin.Context) {
	token, err := h.cookies.AccessTokenFromRequest(c.Request)
	if err != nil {
		h.verifySessionFallback(c)
		return
	}

	userID, err := h.tokens.VerifyAccess(token)
	if err != nil {
		h.verifySessionFallback(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"user_id": userID,
		"tenant_id": domain.TenantIDFromUserID(userID),
	})
}

func (h *AuthHandler) verifySessionFallback(c *gin.Context) {
	refreshToken, err := h.cookies.RefreshTokenFromRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"valid": false,
			"error": "no session cookie",
		})
		return
	}

	userID, err := h.tokens.VerifyRefresh(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"valid": false,
			"error": "invalid or expired session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"user_id": userID,
		"tenant_id": domain.TenantIDFromUserID(userID),
		"access_expired": true,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	h.cookies.ClearSession(c.Writer)
	c.JSON(http.StatusOK, gin.H{"status": "logged_out"})
}

func (h *AuthHandler) setSession(c *gin.Context, userID string) error {
	pair, err := h.tokens.GeneratePair(userID)
	if err != nil {
		return err
	}
	h.cookies.SetSession(c.Writer, pair.AccessToken, pair.RefreshToken)
	return nil
}
