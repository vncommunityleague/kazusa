package registry

import (
	"context"

	"github.com/vncommunityleague/kazusa/connection"
	"github.com/vncommunityleague/kazusa/internal"
	"github.com/vncommunityleague/kazusa/ory"
	"github.com/vncommunityleague/kazusa/repo"
)

type Default struct {
	repo.Repository

	kratos ory.Kratos

	connectionHandler *connection.Handler
	connectionManager *connection.Manager
}

func NewRegistryDefault() Registry {
	return &Default{
		Repository: repo.NewRepository(),
	}
}

func (m *Default) RegisterPublicRoutes(ctx context.Context, router *internal.PublicRouter) {
	m.ConnectionHandler().RegisterPublicRoutes(router)
}

func (m *Default) Kratos() ory.Kratos {
	if m.kratos == nil {
		m.kratos = ory.NewDefaultKratos(m)
	}

	return m.kratos
}

func (m *Default) ConnectionHandler() *connection.Handler {
	if m.connectionHandler == nil {
		m.connectionHandler = connection.NewHandler(m)
	}

	return m.connectionHandler
}

func (m *Default) ConnectionManager() *connection.Manager {
	if m.connectionManager == nil {
		m.connectionManager = connection.NewManager(m)
	}

	return m.connectionManager
}
