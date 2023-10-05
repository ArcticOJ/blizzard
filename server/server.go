package server

import (
	"blizzard/config"
	"blizzard/logger"
	"blizzard/routes"
	"blizzard/server/http"
	"blizzard/server/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	nethttp "net/http"
)

func createHandler(handler http.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := &http.Context{
			Context: c,
		}
		res := handler(ctx)
		if ctx.Response().Committed {
			return nil
		}
		if res == nil {
			return ctx.NoContent(nethttp.StatusNoContent)
		} else {
			return ctx.CommitResponse(res)
		}
	}
}

func Register(e *echo.Echo) {
	g := e.Group("/api", middlewares.Authentication())
	if config.Config.RateLimit > 0 {
		g.Use(middlewares.RateLimit())
	}
	if config.Config.Debug {
		g.Use(middleware.BodyDump(func(c echo.Context, req, res []byte) {
			logger.Blizzard.Debug().Str("url", c.Request().RequestURI).Bytes("req", req).Bytes("res", res).Msg("body")
		}))
		g.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:    true,
			LogStatus: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				logger.Blizzard.Debug().
					Str("url", v.URI).
					Int("status", v.Status).
					Dur("latency", v.Latency).
					Msg("request")
				return nil
			},
		}))
	}
	for route, group := range routes.Map {
		_g := g.Group(route)
		// TODO: allow middlewares
		for r, sub := range group {
			for _, m := range sub.Methods {
				method := m.ToString()
				handler := createHandler(sub.Handler)
				if route == "/" {
					// add handler for root routes
					g.Add(method, r, handler)
				} else if r == "/" {
					// handle apex routes
					g.Add(method, route, handler)
				} else {
					_g.Add(method, r, handler)
				}
			}
		}
	}
}
