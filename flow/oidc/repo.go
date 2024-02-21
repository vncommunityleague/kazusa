package oidc

import "context"

type Repository interface {
	UpsertOIDCFlow(ctx context.Context, state string, flow *Flow) error

	GetAndDeleteOIDCFlow(ctx context.Context, state string) (*Flow, error)
}
