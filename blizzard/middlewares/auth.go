package middlewares

import (
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/user"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"strings"
)

func invalidate(ctx *extra.Context, next echo.HandlerFunc) error {
	ctx.Set("user", nil)
	ctx.DeleteCookie("session")
	return next(ctx)
}

func Authentication(secret string) echo.MiddlewareFunc {
	key := []byte(secret)
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
			ctx := &extra.Context{
				Context: c,
			}
			jt, e := ctx.Cookie("session")
			if e != nil {
				return invalidate(ctx, next)
			}
			token, err := jwt.ParseWithClaims(jt.Value, &models.Session{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return key, nil
			})
			if err != nil {
				return invalidate(ctx, next)
			}
			if session, ok := token.Claims.(*models.Session); ok && token.Valid {
				if ok, e := db.Database.NewSelect().Model((*user.User)(nil)).Where("id = ?", session.UUID).Exists(ctx.Request().Context()); ok && e == nil {
					ctx.Set("user", session.UUID)
				} else {
					return invalidate(ctx, next)
				}
			} else {
				return invalidate(ctx, next)
			}
			return next(c)
		}
	}
}
