package jwt

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secretKey string
}

func NewManager(secret string) *Manager {
	return &Manager{
		secretKey: secret,
	}
}
type Claims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

func (m *Manager) Generate(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(15 * time.Minute),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(m.secretKey))
}

func (m *Manager) Verify(tokenString string) (string, error) {

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
		return "", err
	}

	return claims.UserID, nil
}