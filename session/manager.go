package session

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/securecookie"
)

const (
	ErrUnauthorized = "unauthorized"
)

type (
	managerDependencies interface {
		Repository
		SecureCookieProvider
	}

	ManagerProvider interface {
		SessionManager() *Manager
	}

	Manager struct {
		r managerDependencies
	}
)

const (
	CookieName = "session-token"
)

func NewManager(d managerDependencies) *Manager {
	return &Manager{
		d,
	}
}

func (m *Manager) IssueCookie(w http.ResponseWriter, r *http.Request, s *Session) {
	if encoded, err := securecookie.EncodeMulti(CookieName, s.Token, m.r.SecureCookie()); err == nil {
		cookie := &http.Cookie{
			Name:    CookieName,
			Value:   encoded,
			Path:    "/",
			Expires: s.ExpiresAt.UTC(),
			MaxAge:  0,
		}

		http.SetCookie(w, cookie)
	} else {
		http.Error(w, "Failed to issue cookie", http.StatusInternalServerError)
	}
}

func (m *Manager) RefreshCookie(w http.ResponseWriter, r *http.Request, s *Session) {
	m.IssueCookie(w, r, s)
}

func (m *Manager) GetSessionFromRequest(ctx context.Context, r *http.Request) (*Session, error) {
	token := m.extractToken(r)
	if token == "" {
		return nil, errors.New("no token found")
	}

	return m.r.GetSessionByToken(ctx, token)
}

func (m *Manager) extractToken(r *http.Request) string {
	ctx := r.Context()

	cookie, err := r.Cookie(CookieName)
	if err == nil {
		var token string
		if err = m.r.SecureCookie().Decode(CookieName, cookie.Value, &token); err == nil {
			return token
		}
	}

	token, _ := bearerTokenFromRequest(r.WithContext(ctx))
	return token
}

// From ory/kratos
func bearerTokenFromRequest(r *http.Request) (string, bool) {
	parts := strings.Split(r.Header.Get("Authorization"), " ")

	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1], true
	}

	return "", false
}
