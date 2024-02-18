package oidc

import "context"

type Repository interface {
	UpsertFlow(ctx context.Context, state string, flow *Flow) error

	GetFlow(ctx context.Context, state string) (*Flow, error)
}
