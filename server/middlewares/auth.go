package middlewares

import (
	"github.com/ArcticOJ/blizzard/v0/cache/stores"
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/models/user"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/ArcticOJ/blizzard/v0/server/session"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"strings"
)

func invalidate(ctx *http.Context, next echo.HandlerFunc) error {
	ctx.Set("user", nil)
	ctx.DeleteCookie("session")
	return next(ctx)
}

func Authentication() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if authHeader := c.Request().Header.Get("Authorization"); strings.HasPrefix(authHeader, "Bearer") {
				authToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
				if len(authToken) > 0 {
					var usr user.User
					if e := db.Database.NewSelect().Model(&usr).Where("api_key = ?", authToken).Column("id").Scan(c.Request().Context()); e == nil {
						c.Set("user", usr.ID)
						return next(c)
					}
				}
			}
			ctx := &http.Context{
				Context: c,
			}
			cookie, e := ctx.Cookie("session")
			if e != nil || cookie.Value == "" {
				return invalidate(ctx, next)
			}
			if uid := session.Decrypt(cookie.Value); uid != uuid.Nil {
				u := stores.Users.GetMinimal(c.Request().Context(), uid)
				if u == nil {
					return invalidate(ctx, next)
				}
				ctx.Set("id", uid)
				ctx.Set("user", u)
			}
			return next(c)
		}
	}
}
