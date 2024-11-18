package registry

import (
	"context"

	"github.com/vncommunityleague/kazusa/connection"
	"github.com/vncommunityleague/kazusa/internal"
	"github.com/vncommunityleague/kazusa/repo"
)

type Default struct {
	repo.Repository

	connectionHandler *connection.Handler
}

func NewRegistryDefault() Registry {
	return &Default{
		Repository: repo.NewRepository(),
	}
}

func (m *Default) RegisterPublicRoutes(ctx context.Context, router *internal.PublicRouter) {
	m.ConnectionHandler().RegisterPublicRoutes(router)
}

func (m *Default) ConnectionHandler() *connection.Handler {
	if m.connectionHandler == nil {
		m.connectionHandler = connection.NewHandler(m)
	}

	return m.connectionHandler
}
