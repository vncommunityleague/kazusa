package game

type GameData struct {
	PlayTime uint64 `json:"play_time"`

	GlobalRank  uint32 `json:"global_rank"`
	RegionRank  uint32 `json:"region_rank"`
	CountryRank uint32 `json:"country_rank"`
}
