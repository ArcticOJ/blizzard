package user

import (
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/user"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/uptrace/bun"
	"strings"
)

func Info(ctx *extra.Context) models.Response {
	handle := strings.ToLower(ctx.Param("handle"))
	if len(handle) == 0 {
		return ctx.Bad("Invalid handle.")
	}
	var c []user.User
	fmt.Println(db.Database.NewSelect().Model(&c).Where("handle = ?", handle).Relation("Roles").Relation("Connections", func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("show_in_profile = true")
	}).ExcludeColumn("email_verified", "api_key", "password").Scan(ctx.Request().Context()))
	if len(c) != 1 {
		return ctx.NotFound("User not found.")
	}
	h := md5.Sum([]byte(c[0].Email))
	c[0].Avatar = hex.EncodeToString(h[:])
	c[0].Email = ""
	return ctx.Respond(c[0])
}
