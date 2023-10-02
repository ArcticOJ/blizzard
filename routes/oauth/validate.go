package oauth

import (
	"blizzard/db"
	"blizzard/db/models/user"
	"blizzard/logger/debug"
	"blizzard/oauth"
	"blizzard/oauth/providers"
	"blizzard/server/http"
	"blizzard/utils"
	"blizzard/utils/crypto"
	"crypto/hmac"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun/driver/pgdriver"
	"strings"
)

type oauthValidationRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

func HandleLink(ctx *http.Context, provider string, res *providers.UserInfo) http.Response {
	uuid := ctx.GetUUID()
	if _, e := db.Database.NewInsert().Model(&user.OAuthConnection{
		UserID:   *uuid,
		Username: res.Username,
		ID:       res.ID,
		Provider: provider,
	}).Exec(ctx.Request().Context()); e != nil {
		var err pgdriver.Error
		if errors.As(e, &err) && err.Field('C') == pgerrcode.UniqueViolation {
			return ctx.Forbid("This ID is already bound to another account.")
		}
		debug.Dump(e)
	}
	return nil
}

func HandleLogin(ctx *http.Context, provider string, res *providers.UserInfo, state []string) http.Response {
	var c []user.OAuthConnection
	if e := db.Database.NewSelect().
		Model(&c).
		Where("provider = ? AND id = ?", provider, res.ID).
		Limit(1).
		Scan(ctx.Request().Context()); e != nil {
		debug.Dump(e)
		return ctx.NotFound("This ID is not linked to any accounts.")
	} else {
		if len(c) != 1 {
			return ctx.NotFound("This ID is not linked to any accounts.")
		}
		ctx.Authenticate(c[0].UserID, len(state) == 3 && state[2] == "1")
	}
	return nil
}

func Validate(ctx *http.Context) http.Response {
	name := ctx.Param("provider")
	if prov, ok := oauth.Conf[name]; ok {
		var body oauthValidationRequest
		if ctx.Bind(&body) != nil {
			return ctx.Bad("Malformed body.")
		}
		state := strings.Split(body.State, "#")
		if len(state) < 2 || len(state) > 3 {
			return ctx.Bad("Invalid OAuth state.")
		}
		action := utils.DecodeBase64ToString(state[0])
		hash := utils.DecodeBase64ToBytes(state[1])
		if !utils.ArrayIncludes(oauth.AllowedActions, action) {
			return ctx.Bad("Invalid action, supported actions are: " + strings.Join(oauth.AllowedActions, ", "))
		}
		if hash == nil || (action != "link" && !hmac.Equal(crypto.Hash([]byte(action)), hash)) {
			return ctx.Bad("OAuth hash mismatch.")
		} else if action == "link" {
			u := ctx.GetUUID().String()
			h := crypto.Hash([]byte(action + "_" + strings.ReplaceAll(u, "-", "")))
			if !hmac.Equal(h, hash) {
				return ctx.Bad("OAuth hash mismatch.")
			}
		}
		if action == "link" && ctx.RequireAuth() {
			return nil
		}
		token, e := prov.Exchange(ctx.Request().Context(), body.Code)
		if e != nil {
			return ctx.Bad("Failed to exchange for access token.", echo.Map{
				"action": action,
			})
		}
		client := prov.Client(ctx.Request().Context(), token)
		res := oauth.UserInfoHandler[name](client)
		if res == nil {
			return ctx.Bad("Unable to get user info.")
		}
		switch action {
		case "link":
			if r := HandleLink(ctx, name, res); r != nil {
				return r
			}
			break
		case "login":
			// TODO: login callback url
			if r := HandleLogin(ctx, name, res, state); r != nil {
				return r
			}
			break
		}
		return ctx.Respond(echo.Map{
			"action": action,
			"user":   res,
		})
	}
	return ctx.Bad("Unsupported OAuth provider.")
}
