package auth

import (
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/users"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"github.com/matthewhartstonge/argon2"
)

type loginRequest struct {
	Handle     string `json:"handle"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe,omitempty"`
}

// TODO: Validate req before processing

func Login(ctx *extra.Context) models.Response {
	var req loginRequest
	if ctx.Bind(&req) != nil {
		return ctx.Bad("Malformed credentials.")
	}
	var user users.User
	if e := db.Database.NewSelect().Model(&user).Where("handle = ?", req.Handle).WhereOr("email = ?", req.Handle).Column("uuid", "password").Scan(ctx.Request().Context()); e != nil {
		return ctx.NotFound("User not found.")
	}
	if r, _ := argon2.VerifyEncoded([]byte(req.Password), []byte(user.Password)); !r {
		return ctx.Bad("Wrong password.")
	} else {
		return ctx.Authenticate(user.UUID, req.RememberMe)
	}
}
