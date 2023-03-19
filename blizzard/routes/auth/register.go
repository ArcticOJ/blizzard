package auth

import (
	"backend/blizzard/core"
	"backend/blizzard/db/models/shared"
	"backend/blizzard/models"
	"github.com/jackc/pgerrcode"
	"github.com/uptrace/bun/driver/pgdriver"
)

type (
	RegisterForm struct {
		DisplayName  string `json:"displayName,omitempty"`
		Handle       string `json:"handle"`
		Email        string `json:"email"`
		Password     string `json:"password"`
		Organization string `json:"organization,omitempty"`
	}
)

// TODO: Validate req before processing

func Register(ctx *models.Context) models.Response {
	var req RegisterForm
	if ctx.Bind(&req) != nil {
		return ctx.Bad("Malformed request payload.")
	}
	r, e := core.HashConfig.HashEncoded([]byte(req.Password))
	if e != nil {
		return ctx.InternalServerError("Could not hash provided password.")
	}
	_, err := ctx.Database.NewInsert().Model(&shared.User{
		DisplayName:  req.DisplayName,
		Handle:       req.Handle,
		Email:        req.Email,
		Password:     string(r),
		Organization: req.Organization,
	}).Returning("NULL").Exec(ctx.Request().Context())
	if err != nil {
		if err, ok := err.(pgdriver.Error); ok && err.Field('C') == pgerrcode.UniqueViolation {
			return ctx.Err(403, "User with the same email or username already exists.", nil)
		}
		return ctx.InternalServerError("Request failed with unexpected error.")
	}
	return ctx.Success()
}
