package repo

import (
	"context"
	"github.com/vncommunityleague/kazusa/identity"

	"github.com/redis/rueidis"

	"github.com/vncommunityleague/kazusa/flow/oidc"
)

func (r *repositoryImpl) UpsertOIDCFlow(ctx context.Context, state string, flow *oidc.Flow) error {
	return r.d.Rds.Do(ctx, r.d.Rds.B().Set().Key(state).Value(rueidis.JSON(flow)).Build()).Error()
}

func (r *repositoryImpl) GetAndDeleteOIDCFlow(ctx context.Context, state string) (*oidc.Flow, error) {
	var flow oidc.Flow
	if err := r.d.Rds.Do(ctx, r.d.Rds.B().Get().Key(state).Build()).DecodeJSON(&flow); err != nil {
		return nil, err
	}
	r.d.Rds.Do(ctx, r.d.Rds.B().Del().Key(state).Build())

	return &flow, nil
}

func (r *repositoryImpl) GetIdentityByDiscordID(ctx context.Context, discordID string) (*identity.Identity, error) {
	var i identity.Identity
	if err := r.d.DB.WithContext(ctx).FirstOrCreate(&i, identity.Identity{DiscordId: discordID}).Error; err != nil {
		return nil, err
	}

	return &i, nil
}

func (r *repositoryImpl) GetIdentityByOsuID(ctx context.Context, osuID uint) (*identity.Identity, error) {
	var i identity.Identity
	if err := r.d.DB.WithContext(ctx).FirstOrCreate(&i, identity.Identity{OsuId: osuID}).Error; err != nil {
		return nil, err
	}

	return &i, nil
}
