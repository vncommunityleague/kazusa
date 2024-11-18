package connection

import (
	"context"

	"github.com/redis/rueidis/om"
)

type Repository interface {
	GetConnectionFlows() om.Repository[Flow]

	SaveConnection(ctx context.Context, con *Connection) error
}
