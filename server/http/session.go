package http

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Session struct {
	UUID uuid.UUID `json:"uuid"`
	jwt.RegisteredClaims
}
