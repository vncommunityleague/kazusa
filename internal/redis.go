package internal

import "time"

type ReidsBase struct {
	Key       string    `json:"key" redis:",key"`
	Ver       int64     `json:"ver" redis:",ver"`
	ExpiresAt time.Time `json:"exat" redis:",exat"`
}
