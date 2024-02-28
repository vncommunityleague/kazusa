package internal

import (
	"encoding/base64"
	"github.com/gorilla/securecookie"
)

func RandomString(length int) string {
	return base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(length))
}
