package oauth

import (
	"encoding/base64"
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/oauth"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/ArcticOJ/blizzard/v0/utils"
	"github.com/ArcticOJ/blizzard/v0/utils/crypto"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"strings"
)

func CreateUrl(ctx *http.Context) http.Response {
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
