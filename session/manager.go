package session

import (
	"context"
	"net/http"
	"strings"
)

type (
	managerDependencies interface {
		Repository
	}

	ManagerProvider interface {
		SessionManager() *Manager
	}

	Manager struct {
		r managerDependencies
	}
)

func (m *Manager) IssueToken(ctx context.Context) {

}

func (m *Manager) GetSessionFromRequest(ctx context.Context, r *http.Request) (*Session, error) {
	token := m.extractToken(r)
	if token == "" {
		return nil, nil
	}

	return m.r.GetSessionByToken(ctx, token)
}

func (m *Manager) extractToken(r *http.Request) string {
	ctx := r.Context()

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
