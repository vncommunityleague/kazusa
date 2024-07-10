package game

import (
	"context"
	"fmt"
	"time"

	"github.com/vncommunityleague/kazusa/connection"
	"github.com/vncommunityleague/kazusa/internal"
)

type osuGameProvider struct {
	connection.SelfAuthProvider

	d providerDependencies
}

func NewOsuGameProvider(d providerDependencies) GameProvider {
	return &osuGameProvider{
		SelfAuthProvider: connection.NewOsuSelfAuthProvider(),
		d:                d,
	}
}

type OsuData struct {
	GameData

	ID uint32 `json:"id,omitempty" redis:"-" gorm:"-"`

	Level uint32 `json:"level"`
	Score uint64 `json:"score"`

	HitAccuracy      float64 `json:"hit_accuracy"`
	PerformancePoint float64 `json:"pp"`
}

type OsuDataRedis struct {
	internal.ReidsBase
	OsuData
}

func toOsuDataRedis(data OsuData) *OsuDataRedis {
	return &OsuDataRedis{
		ReidsBase: internal.ReidsBase{
			Key:       fmt.Sprint(data.ID),
			ExpiresAt: time.Now().Add(5 * time.Minute),
		},
		OsuData: data,
	}
}

func (p *osuGameProvider) GetMultiUserGameData(ctx context.Context, ids []string) (any, error) {
	c, err := p.ClientCredentials()
	if err != nil {
		return nil, err
	}

	result := []*OsuData{}
	remainingIds := []string{}

	for _, id := range ids {
		cache, err := p.d.GetOsuGameDataCache().Fetch(ctx, id)

		if err == nil {
			remainingIds = append(remainingIds, id)
		} else {
			result = append(result, &cache.OsuData)
		}
	}

	if len(remainingIds) > 0 {
		httpClient := c.Client(ctx)

		var data struct {
			users []OsuData
		}

		if err := internal.RequestOAuthData(httpClient, connection.OsuGetUsersURL, data); err != nil {
			return nil, err
		}

		savings := []*OsuDataRedis{}
		for _, element := range data.users {
			result = append(result, &element)
			savings = append(savings, toOsuDataRedis(element))
		}

		p.d.GetOsuGameDataCache().SaveMulti(ctx, savings...)
	}

	return &result, nil
}
