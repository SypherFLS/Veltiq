package cookies

import (
	"errors"
	"net/http"
)

const (
	AccessTokenName = "access_token"
	RefreshTokenName = "refresh_token"
)

var (
	ErrNoAccessToken = errors.New("access token cookie not found")
	ErrNoRefreshToken = errors.New("refresh token cookie not found")
)

type Manager struct {
	secure bool
	accessMaxAge int
	refreshMaxAge int
	sameSite http.SameSite
}

func NewManager(secure bool, accessMaxAgeSec, refreshMaxAgeSec int) *Manager {
	return &Manager{
		secure: secure,
		accessMaxAge: accessMaxAgeSec,
		refreshMaxAge: refreshMaxAgeSec,
		sameSite: http.SameSiteLaxMode,
	}
}

func (m *Manager) SetSession(w http.ResponseWriter, accessToken, refreshToken string) {
	m.setNamed(w, AccessTokenName, accessToken, m.accessMaxAge)
	m.setNamed(w, RefreshTokenName, refreshToken, m.refreshMaxAge)
}

func (m *Manager) ClearSession(w http.ResponseWriter) {
	m.setNamed(w, AccessTokenName, "", -1)
	m.setNamed(w, RefreshTokenName, "", -1)
}

func (m *Manager) AccessTokenFromRequest(r *http.Request) (string, error) {
	return m.tokenFromRequest(r, AccessTokenName, ErrNoAccessToken)
}

func (m *Manager) RefreshTokenFromRequest(r *http.Request) (string, error) {
	return m.tokenFromRequest(r, RefreshTokenName, ErrNoRefreshToken)
}

func (m *Manager) tokenFromRequest(r *http.Request, name string, errEmpty error) (string, error) {
	c, err := r.Cookie(name)
	if err != nil {
		return "", errEmpty
	}
	if c.Value == "" {
		return "", errEmpty
	}
	return c.Value, nil
}

func (m *Manager) setNamed(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name: name,
		Value: value,
		Path: "/",
		HttpOnly: true,
		Secure: m.secure,
		SameSite: m.sameSite,
		MaxAge: maxAge,
	})
}
