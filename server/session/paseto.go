package session

import (
	"aidanwoods.dev/go-paseto"
	"github.com/ArcticOJ/blizzard/v0/config"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/google/uuid"
	"time"
)

var (
	key    paseto.V4SymmetricKey
	parser paseto.Parser
)

func init() {
	var e error
	key, e = paseto.V4SymmetricKeyFromHex(config.Config.PrivateKey)
	logger.Panic(e, "failed to decode hex-encoded private key from config, please regenerate another one using cmd/generator")
	parser = paseto.NewParser()
}

func Decrypt(cookie string) uuid.UUID {
	t, e := parser.ParseV4Local(key, cookie, nil)
	if e != nil {
		return uuid.Nil
	}
	s, e := t.GetString("id")
	if e != nil || s == "" {
		return uuid.Nil
	}
	uid, _ := uuid.Parse(s)
	return uid
}

func Encrypt(lifespan time.Duration, uid uuid.UUID) (k string, validUntil time.Time) {
	validUntil = time.Now().Add(lifespan)
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(validUntil)
	token.SetIssuer("Arctic Judge Platform")
	token.SetString("id", uid.String())
	k = token.V4Encrypt(key, nil)
	return
}
