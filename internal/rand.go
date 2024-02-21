package internal

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func RandomBytesInHex(count int) (string, error) {
	buf := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic("Impl error handler for randomBytesInHex")
	}

	return hex.EncodeToString(buf), nil
}
