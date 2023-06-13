package http

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/logger"
	"blizzard/blizzard/middlewares"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"net/http"
)

func createHandler(handler extra.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := &extra.Context{
			Context: c,
		}
		res := handler(ctx)
		if ctx.Response().Committed {
			return nil
		}
		if res == nil {
			return ctx.NoContent(http.StatusNoContent)
		} else {
			return ctx.CommitResponse(res)
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
					Str("url", v.URI).
					Int("status", v.Status).
					Msg("req")
				return nil
			},
		}))
	}
	e.Use(middlewares.Authentication(config.Config.PrivateKey))
	for route, group := range routes.Map {
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
