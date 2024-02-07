package user

import (
	"encoding/base64"
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/schema/user"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/ArcticOJ/blizzard/v0/utils"
	"github.com/labstack/echo/v4"
	"strconv"
	"time"
)

func ApiKey(ctx *http.Context) http.Response {
	if ctx.RequireAuth() {
		return nil
	}
	switch ctx.Method() {
	case http.Get:
		var apiKey string
		uuid := ctx.GetUUID()
		if db.Database.NewSelect().Model(&user.User{ID: uuid}).Column("id").WherePK().Column("api_key").Scan(ctx.Request().Context(), apiKey) != nil {
			apiKey = ""
		}
		return ctx.Respond(echo.Map{
			"apiKey": apiKey,
		})
	case http.Patch:
		uuid := ctx.GetUUID()
		now := base64.RawStdEncoding.EncodeToString([]byte(strconv.FormatInt(time.Now().UTC().Unix(), 10)))
		hash := utils.Rand(10, "")
		if hash == "" {
			return ctx.InternalServerError("Could not generate an API key.")
		}
		apiKey := "arctic." + hash + now
		if _, e := db.Database.NewUpdate().Model(&user.User{
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
