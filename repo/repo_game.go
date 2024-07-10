package repo

import (
	"github.com/redis/rueidis/om"
	"github.com/vncommunityleague/kazusa/game"
)

func (r *Default) GetOsuGameDataCache() om.Repository[game.OsuDataRedis] {
	return r.OsuGameDataRepo
}
