package identity

import (
	"time"

	"github.com/google/uuid"
)

type Identity struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v1()"`
	// Socials
	DiscordId string `json:"discord_id" gorm:"unique"`
	// Games
	OsuId string `json:"osu_id" gorm:"unique"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
