package connection

import (
	"context"

	"github.com/redis/rueidis/om"
)

type Repository interface {
	GetConnectionFlows() om.Repository[Flow]

	SaveConnections(ctx context.Context, conns *Connections) error

	GetConnectionsByID(ctx context.Context, key string) (*Connections, error)
}
