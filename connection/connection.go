package connection

import (
	"time"

	"github.com/google/uuid"
)

type UserConnections struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`

	Osu OsuConnection `json:"osu" gorm:"type:bytes;serializer:gob"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type OsuConnection struct {
	Id        uint32 `json:"id"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
}
