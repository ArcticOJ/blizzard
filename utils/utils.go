package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"strconv"
)

func Rand(l uint8, fallback string) string {
	bytes := make([]byte, l)
	if _, err := rand.Read(bytes); err != nil {
		return fallback
	}
	return hex.EncodeToString(bytes)
}

func ArrayFill[T any](val T, count int) (arr []T) {
	arr = make([]T, count)
	for i := range arr {
		arr[i] = val
	}
	return
}

func DecodeBase64ToString(b64 string) string {
	b, e := base64.RawStdEncoding.DecodeString(b64)
	if e != nil {
		return ""
	}
	return string(b)
}

func DecodeBase64ToBytes(b64 string) []byte {
	buf := make([]byte, base64.RawStdEncoding.DecodedLen(len(b64)))
	d, err := base64.RawStdEncoding.Decode(buf, []byte(b64))
	if err != nil {
		return nil
	}
	return buf[:d]
}

func ParseInt(s string) int {
	v, e := strconv.Atoi(s)
	if e != nil {
		return 0
	}
	return v
}
