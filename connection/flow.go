package connection

import "time"

type Flow struct {
	Key       string    `json:"key" redis:",key"`
	Ver       int64     `json:"ver" redis:",ver"`
	ExpiresAt time.Time `json:"exat" redis:",exat"`

	SessionId    string `json:"session_id"`
	CodeVerifier string `json:"code_verifier"`
	Url          string `json:"url"`
}
