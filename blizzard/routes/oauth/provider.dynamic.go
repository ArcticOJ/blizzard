package oauth

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/oauth"
	"blizzard/blizzard/utils"
	"blizzard/blizzard/utils/crypto"
	"encoding/base64"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"strings"
)

func CreateUrl(ctx *extra.Context) models.Response {
	if prov, ok := oauth.Conf[ctx.Param("provider")]; ok {
		// 0 is false, 1 is true
		remember := 0
		action := ctx.QueryParam("action")
		if ctx.QueryParam("remember") == "true" {
			remember = 1
		}
		if !utils.ArrayIncludes(oauth.AllowedActions, action) {
			return ctx.Bad("Invalid action.")
		}
		if action == "link" && ctx.RequireAuth() {
			return nil
		}
		var state string
		if action == "link" {
			hash := base64.RawStdEncoding.EncodeToString(crypto.Hash([]byte(action + "_" + strings.ReplaceAll(ctx.GetUUID().String(), "-", ""))))
			state = fmt.Sprintf("%s#%s", base64.RawStdEncoding.EncodeToString([]byte(action)), hash)
		} else {
			hash := base64.RawStdEncoding.EncodeToString(crypto.Hash([]byte(action)))
			state = fmt.Sprintf("%s#%s#%d", base64.RawStdEncoding.EncodeToString([]byte(action)), hash, remember)
		}
		return ctx.Respond(echo.Map{
			"url": prov.AuthCodeURL(state, oauth2.AccessTypeOnline),
		})
	}
	return ctx.Bad("Unsupported OAuth provider.")
}
