package connection

import (
	"time"

	"github.com/google/uuid"
)

type Connection struct {
	ID       uint64    `json:"id"`
	UserId   uuid.UUID `json:"user_id" gorm:"type:uuid"`
	Provider string    `json:"provider"`

	ConnId    string `json:"conn_id"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
