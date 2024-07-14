package repo

import (
	"context"

	"github.com/redis/rueidis/om"

	"github.com/vncommunityleague/kazusa/connection"
)

func (r *Default) GetConnectionFlows() om.Repository[connection.Flow] {
	return r.ConnectionFlowRepo
}

func (r *Default) SaveConnections(ctx context.Context, conns *connection.Connections) error {
	return r.d.DB.WithContext(ctx).Save(conns).Error
}

func (r *Default) GetConnectionsByID(ctx context.Context, id string) (*connection.Connections, error) {
	var i connection.Connections
	if err := r.d.DB.WithContext(ctx).First(&i, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &i, nil
}
