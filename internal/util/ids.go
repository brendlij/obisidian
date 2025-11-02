package util

import (
	"crypto/rand"
	"encoding/hex"
)

func RandID() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
