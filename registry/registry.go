package registry

import (
	"context"
	"github.com/vncommunityleague/kazusa/connection"

	"github.com/vncommunityleague/kazusa/internal"
	"github.com/vncommunityleague/kazusa/ory"
	"github.com/vncommunityleague/kazusa/repo"
)

type Registry interface {
	repo.Repository

	ory.Provider

	connection.HandlerProvider
	connection.ManagementProvider

	RegisterPublicRoutes(ctx context.Context, public *internal.PublicRouter)
}
