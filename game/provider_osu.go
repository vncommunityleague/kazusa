package game

import (
	"context"
	"fmt"
	"strings"
	"time"

	"net/http"

	"github.com/vncommunityleague/kazusa/connection"
	"github.com/vncommunityleague/kazusa/internal"
)

type osuGameProvider struct {
	connection.SelfAuthProvider

	d providerDependencies
}

func NewOsuProvider(d providerDependencies) Provider {
	return &osuGameProvider{
		SelfAuthProvider: connection.NewOsuSelfAuthProvider(),
		d:                d,
	}
}

type OsuData struct {
	baseData

	Level       uint32 `json:"level"`
	RankedScore uint64 `json:"ranked_score"`
	TotalScore  uint64 `json:"total_score"`

	HitAccuracy      float64 `json:"hit_accuracy"`
	PerformancePoint float64 `json:"performance_point"`
}

type OsuDataRedis struct {
	Key       string    `json:"key" redis:",key"`
	Ver       int64     `json:"ver" redis:",ver"`
	ExpiresAt time.Time `json:"exat" redis:",exat"`

	OsuData
}

func toOsuDataRedis(mode string, data OsuData) *OsuDataRedis {
	return &OsuDataRedis{
		Key:       mode + "_" + fmt.Sprint(data.ID),
		ExpiresAt: time.Now().Add(5 * time.Minute),
		OsuData:   data,
	}
}

type osuUserStaticsLevel struct {
	Current uint32 `json:"current,omitempty"`
}

type osuUserStatistics struct {
	PlayTime uint64 `json:"play_time,omitempty"`

	GlobalRank  uint32 `json:"global_rank,omitempty"`
	CountryRank uint32 `json:"country_rank,omitempty"`

	Level       osuUserStaticsLevel `json:"level,omitempty"`
	RankedScore uint64              `json:"ranked_score,omitempty"`
	TotalScore  uint64              `json:"total_score,omitempty"`

	HitAccuracy      float64 `json:"hit_accuracy,omitempty"`
	PerformancePoint float64 `json:"pp,omitempty"`
}

func (p *osuGameProvider) GetMultiUserGameData(ctx context.Context, ids []string, extraQuery *ExtraQuery) (any, error) {
	var result []*OsuData
	var savings []*OsuDataRedis

	for _, id := range ids {
		cache, err := p.d.GetOsuGameDataCache().Fetch(ctx, id)

		if err != nil && cache != nil {
			result = append(result, &cache.OsuData)
		} else {
			data, err := p.fetchGameData(ctx, id, extraQuery.Mode)
			if err != nil {
				return nil, err
			}

			result = append(result, data)
			savings = append(savings, toOsuDataRedis(extraQuery.Mode, *data))
		}
	}

	p.d.GetOsuGameDataCache().SaveMulti(ctx, savings...)

	return &result, nil
}

func (p *osuGameProvider) fetchGameData(ctx context.Context, id string, mode string) (*OsuData, error) {
	c, err := p.ClientCredentials()
	if err != nil {
		return nil, err
	}

	if mode == "" || strings.EqualFold(mode, "std") {
		mode = "osu"
	}
	url := strings.ReplaceAll(connection.OsuGetUserURL, "{user}", id)
	url = strings.ReplaceAll(url, "{mode}", mode)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	httpClient := c.Client(ctx)
	var data struct {
		Username string `json:"username,omitempty"`

		Statistics osuUserStatistics `json:"statistics"`
	}
	if err := internal.RequestOAuthData(httpClient, req, &data); err != nil {
		return nil, err
	}

	return &OsuData{
		baseData: baseData{
			ID:          id,
			Username:    data.Username,
			PlayTime:    data.Statistics.PlayTime,
			GlobalRank:  data.Statistics.GlobalRank,
			RegionRank:  data.Statistics.CountryRank,
			CountryRank: data.Statistics.CountryRank,
		},
		Level:            data.Statistics.Level.Current,
		RankedScore:      data.Statistics.RankedScore,
		TotalScore:       data.Statistics.TotalScore,
		HitAccuracy:      data.Statistics.HitAccuracy,
		PerformancePoint: data.Statistics.PerformancePoint,
	}, nil
}
