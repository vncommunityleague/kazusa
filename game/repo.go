package game

import "github.com/redis/rueidis/om"

type Repository interface {
	GetOsuGameDataCache() om.Repository[OsuDataRedis]
}
