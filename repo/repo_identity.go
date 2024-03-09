package repo

import (
	"context"
	"github.com/vncommunityleague/kazusa/identity"
)

func (r *repositoryImpl) UpsertIdentity(ctx context.Context, identity *identity.Identity) error {
	return r.d.DB.WithContext(ctx).Create(&identity).Error
}

func (r *repositoryImpl) GetIdentityByID(ctx context.Context, id string) (*identity.Identity, error) {
	var i identity.Identity
	if err := r.d.DB.WithContext(ctx).First(&i, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &i, nil
}

func (r *repositoryImpl) GetIdentityByDiscordID(ctx context.Context, discordID string) (*identity.Identity, bool, error) {
	var i identity.Identity
	result := r.d.DB.WithContext(ctx).FirstOrCreate(&i, identity.Identity{DiscordId: discordID})
	if result.Error != nil {
		return nil, false, result.Error
	}

	return &i, result.RowsAffected > 0, nil
}

func (r *repositoryImpl) GetIdentityByOsuID(ctx context.Context, osuID uint) (*identity.Identity, error) {
	var i identity.Identity
	if err := r.d.DB.WithContext(ctx).FirstOrCreate(&i, identity.Identity{OsuId: osuID}).Error; err != nil {
		return nil, err
	}

	return &i, nil
}
