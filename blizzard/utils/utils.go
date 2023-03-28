package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func Rand(l uint8, fallback string) string {
	bytes := make([]byte, l)
	if _, err := rand.Read(bytes); err != nil {
		return fallback
	}
	return hex.EncodeToString(bytes)
}
