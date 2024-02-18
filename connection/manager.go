package connection

import (
	"context"

	"github.com/google/uuid"
)

type (
	managerDependencies interface {
		Repository
	}

	ManagementProvider interface {
		ConnectionManager() *Manager
	}

	Manager struct {
		d managerDependencies
	}
)

func NewManager(d managerDependencies) *Manager {
	return &Manager{d: d}
}

func (m *Manager) LinkOrCreateConnections(ctx context.Context, id uuid.UUID, provider string, providerConn interface{}) error {
	conns := &UserConnections{
		ID: id,
	}

	if provider == "osu" {
		conns.Osu = providerConn.(OsuConnection)
	}

	if err := m.d.SaveConnections(ctx, conns); err != nil {
		return err
	}

	return nil
}
