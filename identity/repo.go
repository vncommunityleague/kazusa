package identity

import "context"

type Repository interface {
	UpsertIdentity(ctx context.Context, flow *Identity) error

	GetIdentityByID(ctx context.Context, key string) (*Identity, error)

	GetIdentityByDiscordID(ctx context.Context, discordID string) (*Identity, bool, error)

	GetIdentityByOsuID(ctx context.Context, osuID uint) (*Identity, error)
}
