package auth

import (
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/models/user"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/matthewhartstonge/argon2"
	"strings"
)

type loginRequest struct {
	Handle     string `json:"handle"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe,omitempty" validate:""`
}

// TODO: Validate req before processing

func Login(ctx *http.Context) http.Response {
	var req loginRequest
	if ctx.Bind(&req) != nil {
		return ctx.Bad("Invalid credentials.")
	}
	var usr user.User
	handle := strings.ToLower(req.Handle)
	if e := db.Database.NewSelect().Model(&usr).Where("handle = ?", handle).Column("id", "password").Scan(ctx.Request().Context()); e != nil {
		return ctx.NotFound("User not found.")
	}
	if r, _ := argon2.VerifyEncoded([]byte(req.Password), []byte(usr.Password)); !r {
		return ctx.Bad("Wrong password.")
	} else {
		return ctx.Authenticate(usr.ID, req.RememberMe)
	}
}
