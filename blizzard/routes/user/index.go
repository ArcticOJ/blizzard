package user

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"crypto/md5"
	"fmt"
)

func Index(ctx *extra.Context) models.Response {
	if user := ctx.GetUser(); user != nil {
		user.Avatar = fmt.Sprintf("%x", md5.Sum([]byte(user.Email)))
		return ctx.Respond(user)
	}
	return ctx.NotFound("User not found")
}
