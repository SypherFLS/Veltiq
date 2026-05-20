package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"veltiq/internal/core/ports"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrWrongTokenType = errors.New("wrong token type")
)

type TokenType string

const (
	TokenTypeAccess TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type Manager struct {
	secretKey string
	accessTTL time.Duration
	refreshTTL time.Duration
}

func NewManager(secret string, accessTTL, refreshTTL time.Duration) *Manager {
	if accessTTL <= 0 {
		accessTTL = 15 * time.Minute
	}
	if refreshTTL <= 0 {
		refreshTTL = 7 * 24 * time.Hour
	}
	return &Manager{
		secretKey: secret,
		accessTTL: accessTTL,
		refreshTTL: refreshTTL,
	}
}

type Claims struct {
	UserID string `json:"sub"`
	Type TokenType `json:"typ"`
	jwt.RegisteredClaims
}

func (m *Manager) GeneratePair(userID string) (ports.TokenPair, error) {
	access, err := m.sign(userID, TokenTypeAccess, m.accessTTL)
	if err != nil {
		return ports.TokenPair{}, err
	}
	refresh, err := m.sign(userID, TokenTypeRefresh, m.refreshTTL)
	if err != nil {
		return ports.TokenPair{}, err
	}
	return ports.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (m *Manager) VerifyAccess(tokenString string) (string, error) {
	return m.verify(tokenString, TokenTypeAccess)
}

func (m *Manager) VerifyRefresh(tokenString string) (string, error) {
	return m.verify(tokenString, TokenTypeRefresh)
}

func (m *Manager) sign(userID string, typ TokenType, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		Type: typ,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *Manager) verify(tokenString string, want TokenType) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secretKey), nil
		},
	)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", ErrInvalidToken
	}
	if claims.Type != want {
		return "", ErrWrongTokenType
	}
	return claims.UserID, nil
}
