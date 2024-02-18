package internal

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

func RandomString(length int) string {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(k)
}
