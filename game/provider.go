package game

import (
	"context"
	"errors"
)

type (
	providerDependencies interface {
		Repository
	}

	Provider interface {
		GetMultiUserGameData(ctx context.Context, ids []string, query *ExtraQuery) (any, error)
	}
)

type baseData struct {
	ID       string `json:"id"`
	Username string `json:"username"`

	PlayTime uint64 `json:"play_time"`

	GlobalRank  uint32 `json:"global_rank"`
	RegionRank  uint32 `json:"region_rank"`
	CountryRank uint32 `json:"country_rank"`
}

type ExtraQuery struct {
	Mode string
}

var gameProviders = map[string]func(d providerDependencies) Provider{
	"osu": NewOsuProvider,
}

func GetGameProvider(id string, d providerDependencies) (Provider, error) {
	if p, ok := gameProviders[id]; ok {
		return p(d), nil
	}

	return nil, errors.New("game provider is not found or not supported")
}
