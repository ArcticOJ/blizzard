package server

import (
	"github.com/ArcticOJ/blizzard/v0/config"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/ArcticOJ/blizzard/v0/routes"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/ArcticOJ/blizzard/v0/server/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	nethttp "net/http"
)

func createHandler(route http.Route) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := &http.Context{
			Context: c,
		}
		if route.HasFlag(http.RouteAuth) && ctx.RequireAuth() {
			return nil
		}
		res := route.Handler(ctx)
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
	if config.Config.Blizzard.RateLimit > 0 {
		//g.Use(middlewares.RateLimit())
	}
	if config.Config.Debug {
		g.Use(middleware.BodyDump(func(c echo.Context, req, res []byte) {
			if len(req) > 128 {
				req = []byte("<payload too long>")
			}
			if len(res) > 128 {
				res = []byte("<payload too long>")
			}
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
	for group, _routes := range routes.Map {
		curGroup := g
		// create a dedicated group for non-apex routes
		if group != "/" {
			curGroup = g.Group(group)
		}
		for _, route := range _routes {
			handler := createHandler(route)
			// treat routes like /api/users as apex endpoints
			if route.Path == "/" {
				g.Add(route.Method, group, handler)
			} else {
				curGroup.Add(route.Method, group, handler)
			}
		}
	}
}
