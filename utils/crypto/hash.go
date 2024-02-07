package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"github.com/ArcticOJ/blizzard/v0/config"
	"hash"
)

var hmacHash hash.Hash

func init() {
	hmacHash = hmac.New(sha256.New, []byte(config.Config.Blizzard.Secret))
}

func Hash(buf []byte) []byte {
	hmacHash.Reset()
	hmacHash.Write(buf)
	return hmacHash.Sum(nil)
}
