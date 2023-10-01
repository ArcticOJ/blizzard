package http

import (
	"blizzard/config"
	"blizzard/models/extra"
	"blizzard/routes"
	"blizzard/server/http/middlewares"
	"github.com/labstack/echo/v4"
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

func Register(e *echo.Echo) {
	g := e.Group("/api", middlewares.Authentication(config.Config.PrivateKey))
	// TODO: migrate to redis-based rate limiter
	//e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(config.Config.RateLimit))))
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
