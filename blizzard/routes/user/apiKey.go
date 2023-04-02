package user

import (
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/users"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/utils"
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"strconv"
	"time"
)

func ApiKey(ctx *extra.Context) models.Response {
	if ctx.RequireAuth() {
		return nil
	}
	switch ctx.Method() {
	case models.Get:
		return ctx.Respond(echo.Map{
			"apiKey": ctx.GetUser("api_key").ApiKey,
		})
	case models.Patch:
		uuid := ctx.GetUUID()
		now := base64.RawStdEncoding.EncodeToString([]byte(strconv.FormatInt(time.Now().UTC().Unix(), 10)))
		hash := utils.Rand(10, "")
		if hash == "" {
			return ctx.InternalServerError("Could not generate an API key.")
		}
		apiKey := "arctic." + hash + now
		if _, e := db.Database.NewUpdate().Model(&users.User{
			ApiKey: apiKey,
		}).Column("api_key").Where("uuid = ?", uuid).Returning("NULL").Exec(ctx.Request().Context()); e != nil {
			return ctx.InternalServerError("Could not update your API key.")
		}
		return ctx.Respond(echo.Map{
			"apiKey": apiKey,
		})
	}
	return nil
}
