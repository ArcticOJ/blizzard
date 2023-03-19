package user

import (
	"backend/blizzard/db/models/shared"
	"backend/blizzard/models"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"github.com/labstack/echo/v4"
	"strconv"
	"strings"
	"time"
)

func GenerateApiKey(l uint8) (string, error) {
	bytes := make([]byte, l)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func ApiKey(ctx *models.Context) models.Response {
	switch ctx.Method() {
	case models.Get:
		return ctx.Respond(echo.Map{
			"apiKey": ctx.GetUser("api_key").ApiKey,
		})
	case models.Patch:
		uuid := ctx.GetUUID()
		now := strings.TrimRight(base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(time.Now().UTC().Unix(), 10))), "=")
		hash, e := GenerateApiKey(8)
		if e != nil {
			return ctx.InternalServerError("Could not generate an API key.")
		}
		apiKey := "arctic." + hash + now
		if _, e := ctx.Database.NewUpdate().Model(&shared.User{
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
