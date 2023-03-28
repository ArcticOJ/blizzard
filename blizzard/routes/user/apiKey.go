package user

import (
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/shared"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/utils"
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"strconv"
	"strings"
	"time"
)

func ApiKey(ctx *extra.Context) models.Response {
	switch ctx.Method() {
	case models.Get:
		return ctx.Respond(echo.Map{
			"apiKey": ctx.GetUser("api_key").ApiKey,
		})
	case models.Patch:
		uuid := ctx.GetUUID()
		now := strings.TrimRight(base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(time.Now().UTC().Unix(), 10))), "=")
		hash := utils.Rand(10, "")
		if hash == "" {
			return ctx.InternalServerError("Could not generate an API key.")
		}
		apiKey := "arctic." + hash + now
		if _, e := db.Database.NewUpdate().Model(&shared.User{
			ApiKey: apiKey,
		}).Column("api_key").Where("id = ?", uuid).Returning("NULL").Exec(ctx.Request().Context()); e != nil {
			return ctx.InternalServerError("Could not update your API key.")
		}
		return ctx.Respond(echo.Map{
			"apiKey": apiKey,
		})
	}
	return nil
}
