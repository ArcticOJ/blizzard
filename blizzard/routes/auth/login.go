package auth

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/shared"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"github.com/golang-jwt/jwt/v4"
	"github.com/matthewhartstonge/argon2"
	"time"
)

type (
	LoginForm struct {
		Handle     string `json:"handle"`
		Password   string `json:"password"`
		RememberMe bool   `json:"rememberMe,omitempty"`
	}
)

// TODO: Validate req before processing

func Login(ctx *extra.Context) models.Response {
	var req LoginForm
	if ctx.Bind(&req) != nil {
		return ctx.Bad("Malformed credentials.")
	}
	var user shared.User
	if e := db.Database.NewSelect().Model(&user).Where("handle = ?", req.Handle).WhereOr("email = ?", req.Handle).Column("id", "password").Scan(ctx.Request().Context()); e != nil {
		return ctx.NotFound("User not found.")
	}
	if r, _ := argon2.VerifyEncoded([]byte(req.Password), []byte(user.Password)); !r {
		return ctx.Bad("Wrong password.")
	} else {
		key := []byte(config.Config.PrivateKey)
		now := time.Now()
		lifespan := now.AddDate(0, 1, 0)
		ss := &models.Session{
			UUID: user.ID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(lifespan),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    "Project Arctic",
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, ss)
		signedToken, e := token.SignedString(key)
		if e != nil {
			return ctx.InternalServerError("Could not create a new session.")
		}
		ctx.PutCookie("session", signedToken, lifespan, !req.RememberMe)
		return ctx.Success()
	}
}
