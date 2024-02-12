package middlewares

import (
	"github.com/ArcticOJ/blizzard/v0/cache/stores"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/ArcticOJ/blizzard/v0/server/session"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"strings"
)

func Authentication() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if authHeader := c.Request().Header.Get("Authorization"); strings.HasPrefix(authHeader, "Bearer") {
				authToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
				if len(authToken) > 0 {
					if uid := stores.Users.ResolveApiKey(c.Request().Context(), authToken); uid != uuid.Nil {
						c.Set("id", uid)
						return next(c)
					}
				}
			}
			ctx := &http.Context{
				Context: c,
			}
			cookie, e := ctx.Cookie("session")
			if e != nil || cookie == nil || strings.TrimSpace(cookie.Value) == "" {
				return next(c)
			}
			if uid := session.Decrypt(cookie.Value); uid != uuid.Nil && stores.Users.UserExists(ctx.Request().Context(), uid) {
				ctx.Set("id", uid)
			} else {
				ctx.DeleteCookie("session")
			}
			return next(c)
		}
	}
}
