package game

import (
	"context"
	"errors"
)

type (
	providerDependencies interface {
		Repository
	}

	GameProvider interface {
		GetMultiUserGameData(ctx context.Context, ids []string) (any, error)
	}
)

var gameProviders = map[string]func(d providerDependencies) GameProvider{
	"osu": NewOsuGameProvider,
}

func GetGameProvider(id string, d providerDependencies) (GameProvider, error) {
	if p, ok := gameProviders[id]; ok {
		return p(d), nil
	}

	return nil, errors.New("game provider is not found or not supported")
}
