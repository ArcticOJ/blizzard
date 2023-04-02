package oauth

import (
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/users"
	"blizzard/blizzard/logger/debug"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/oauth"
	"blizzard/blizzard/oauth/providers"
	"blizzard/blizzard/utils"
	"blizzard/blizzard/utils/crypto"
	"crypto/hmac"
	"github.com/labstack/echo/v4"
	"strings"
)

type oauthValidationRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

func HandleLink(ctx *extra.Context, provider string, res *providers.UserInfo) models.Response {
	uuid := ctx.GetUUID()
	var conn []users.OAuthConnection
	if e := db.Database.
		NewSelect().
		Model((*users.OAuthConnection)(nil)).
		Where("id = ?", res.Id).
		WhereOr("user_id = ? AND provider = ?", uuid, provider).
		Scan(ctx.Request().Context(), &conn); e != nil {
		debug.Dump(e)
		return ctx.InternalServerError(e.Error())
	} else {
		if len(conn) == 0 {
			if _, e := db.Database.NewInsert().Model(&users.OAuthConnection{
				UserID:   *uuid,
				Username: res.Username,
				ID:       res.Id,
				Provider: provider,
			}).Returning("NULL").Exec(ctx.Request().Context()); e != nil {
				debug.Dump(e)
				return ctx.InternalServerError("Could not bind this UUID to current account.")
			}
		} else if len(conn) == 1 {
			if conn[0].UserID != *uuid {
				return ctx.Forbid("This UUID is already bound to another account.")
			}
			if conn[0].Username != res.Username {
				if _, e := db.Database.NewUpdate().Model(&users.OAuthConnection{
					Username: res.Username,
				}).Column("username").Where("id = ? AND user_id = ?", res.Id, uuid).Returning("NULL").Exec(ctx.Request().Context()); e != nil {
					debug.Dump(e)
					return ctx.InternalServerError("Could not update account info.")
				}
			}
		} else {
			// TODO: handle 1+ connections with the same id? lol
			return ctx.Forbid("Unexpected error.")
		}
	}
	return nil
}

func HandleLogin(ctx *extra.Context, provider string, res *providers.UserInfo, state []string) models.Response {
	var c []users.OAuthConnection
	if e := db.Database.NewSelect().
		Model((*users.OAuthConnection)(nil)).
		Where("provider = ? AND id = ?", provider, res.Id).
		Limit(1).
		Scan(ctx.Request().Context(), &c); e != nil {
		debug.Dump(e)
		return ctx.NotFound("This UUID is not linked with any accounts.")
	} else {
		if len(c) != 1 {
			return ctx.NotFound("This UUID is not linked with any accounts.")
		}
		ctx.Authenticate(c[0].UserID, len(state) == 3 && state[2] == "1")
	}
	return nil
}

func Validate(ctx *extra.Context) models.Response {
	// TODO: handle oauth callback
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
			if r := HandleLogin(ctx, name, res, state); r != nil {
				return r
			}
		}
		return ctx.Respond(echo.Map{
			"action": action,
			"user":   res,
		})
	}
	return ctx.Bad("Unsupported OAuth provider.")
}
