package repo

import (
	"context"

	"github.com/redis/rueidis/om"

	"github.com/vncommunityleague/kazusa/connection"
)

func (r *Default) GetConnectionFlows() om.Repository[connection.Flow] {
	return r.ConnectionFlowRepo
}

func (r *Default) SaveConnection(ctx context.Context, con *connection.Connection) error {
	return r.d.DB.WithContext(ctx).Save(con).Error
}
