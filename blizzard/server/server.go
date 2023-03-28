package server

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/logger"
	"blizzard/blizzard/middlewares"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/routes/auth"
	"blizzard/blizzard/routes/contests"
	"blizzard/blizzard/routes/feeds"
	"blizzard/blizzard/routes/oauth"
	"blizzard/blizzard/routes/problems"
	"blizzard/blizzard/routes/root"
	"blizzard/blizzard/routes/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"net/http"
)

var Map = map[string]extra.RouteMap{
	"/problems": problems.Map,
	"/feeds":    feeds.Map,
	"/contests": contests.Map,
	"/auth":     auth.Map,
	"/oauth":    oauth.Map,
	"/user":     user.Map,
	"/":         root.Map,
}

func createHandler(handler extra.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		res := handler(&extra.Context{
			Context: c,
		})
		if c.Response().Committed {
			return nil
		}
		if res == nil {
			return c.NoContent(http.StatusNoContent)
		} else {
			return c.JSONPretty(res.StatusCode(), res.Body(), "\t")
		}
	}
}

func CreateServer() *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code, message := http.StatusInternalServerError, "Internal Server Error"
		if er, ok := err.(*echo.HTTPError); ok {
			code = er.Code
			message = er.Message.(string)
		}
		_ = c.JSONPretty(code, models.Error{Code: code, Message: message}, "\t")
	}
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(config.Config.RateLimit))))
	e.Pre(middleware.RemoveTrailingSlash())
	if config.Config.EnableCORS {
		e.Use(middleware.CORS())
	}
	if config.Config.Debug {
		e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:    true,
			LogStatus: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				logger.Logger.Debug().
					Str("URI", v.URI).
					Int("status", v.Status).
					Msg("request")
				return nil
			},
		}))
	}
	e.Use(middlewares.Authentication(config.Config.PrivateKey))
	for route, group := range Map {
		g := e.Group(route)
		for r, sub := range group {
			for _, m := range sub.Methods {
				method := m.ToString()
				handler := createHandler(sub.Handler)
				if route == "/" {
					e.Add(method, r, handler)
				} else if r == "/" {
					e.Add(method, route, handler)
				} else {
					g.Add(method, r, handler)
				}
			}
		}
	}
	return e
}