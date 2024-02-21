package identity

import (
	"time"

	"github.com/google/uuid"
)

type Identity struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v1()"`
	// Socials
	DiscordId string `json:"discord_id"`
	// Games
	OsuId string `json:"osu_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
