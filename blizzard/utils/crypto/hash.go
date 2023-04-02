package crypto

import (
	"blizzard/blizzard/config"
	"crypto/hmac"
	"crypto/sha256"
	"hash"
)

var hmacHash hash.Hash

func Init() {
	hmacHash = hmac.New(sha256.New, []byte(config.Config.PrivateKey))
}

func Hash(buf []byte) []byte {
	hmacHash.Reset()
	hmacHash.Write(buf)
	return hmacHash.Sum(nil)
}
