package auth

import (
	"github.com/ArcticOJ/blizzard/v0/cache/stores"
	"github.com/ArcticOJ/blizzard/v0/core"
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/schema/user"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/jackc/pgerrcode"
	"github.com/uptrace/bun/driver/pgdriver"
	"strings"
)

type registerRequest struct {
	DisplayName  string `json:"displayName,omitempty"`
	Handle       string `json:"handle"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Organization string `json:"organization,omitempty"`
}

// TODO: Validate req before processing

// Register POST /register
func Register(ctx *http.Context) http.Response {
	var req registerRequest
	if ctx.Bind(&req) != nil {
		return ctx.Bad("Malformed request payload.")
	}
	r, e := core.HashConfig.HashEncoded([]byte(req.Password))
	if e != nil {
		return ctx.InternalServerError("Could not crypto provided password.")
	}
	handle := strings.TrimSpace(strings.ToLower(req.Handle))
	var uid string
	_, err := db.Database.NewInsert().Model(&user.User{
		DisplayName: req.DisplayName,
		Handle:      handle,
		Email:       strings.TrimSpace(strings.ToLower(req.Email)),
		Password:    string(r),
		//Organizations: req.Organization,
	}).Returning("id").Exec(ctx.Request().Context(), &uid)
	if err != nil {
		if err, ok := err.(pgdriver.Error); ok && err.Field('C') == pgerrcode.UniqueViolation {
			return ctx.Forbid("User with the same email or handle already exists.")
		}
		return ctx.InternalServerError("Request failed with unexpected error.")
	}
	stores.Users.Add(uid)
	return ctx.Success()
}
