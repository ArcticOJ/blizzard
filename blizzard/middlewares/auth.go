package middlewares

import (
	"backend/blizzard/db/models/shared"
	"backend/blizzard/models"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"strings"
)

func invalidate(ctx models.Context) error {
	ctx.Set("user", nil)
	ctx.DeleteCookie("session")
	return ctx.JSONPretty(401, ctx.Unauthorized().Body(), "\t")
}

func Authentication(secret string, server *models.Server) echo.MiddlewareFunc {
	key := []byte(secret)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.HasPrefix(c.Path(), "/auth/") {
				return next(c)
			}
			if authHeader := c.Request().Header.Get("Authorization"); strings.HasPrefix(authHeader, "Bearer") {
				authToken := strings.TrimSpace(strings.TrimSuffix(authHeader, "Bearer"))
				if len(authToken) > 0 {
					var user shared.User
					if e := server.Database.NewSelect().Model(&user).Where("apiKey = ?", authToken).Column("id").Scan(c.Request().Context()); e == nil {
						c.Set("user", user.ID)
						return next(c)
					}
				}
			}
			ctx := models.Context{
				Server:  server,
				Context: c,
			}
			jt, e := ctx.Cookie("session")
			if e != nil {
				return invalidate(ctx)
			}
			token, err := jwt.ParseWithClaims(jt.Value, &models.Session{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return key, nil
			})
			if err != nil {
				return invalidate(ctx)
			}
			if session, ok := token.Claims.(*models.Session); ok && token.Valid {
				ctx.Set("user", session.UUID)
			} else {
				return invalidate(ctx)
			}
			return next(c)
		}
	}
}
