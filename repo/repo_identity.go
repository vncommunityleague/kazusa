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
