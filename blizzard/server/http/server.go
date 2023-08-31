package http

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/logger"
	"blizzard/blizzard/middlewares"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/routes"
	"blizzard/blizzard/validator"
	"errors"
	"fmt"
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
	if config.Config.EnableCORS {
		e.Use(middleware.CORS())
	}
	if config.Config.Debug {
		e.Use(middleware.BodyDump(func(c echo.Context, req, res []byte) {
			logger.Logger.Debug().Str("url", c.Request().RequestURI).Bytes("req", req).Bytes("res", res).Msg("body")
		}))
		e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:    true,
			LogStatus: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				logger.Logger.Debug().
					Str("url", v.URI).
					Int("status", v.Status).
					Dur("latency", v.Latency).
					Msg("req")
				return nil
			},
		}))
	}
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		code, message := http.StatusInternalServerError, "Internal Server Error"
		var er *echo.HTTPError
		if errors.As(err, &er) {
			code = er.Code
			message = er.Message.(string)
		}
		err = c.JSON(code, models.Error{Code: code, Message: message})
		if err != nil {
			e.Logger.Error(err)
		}
	}
	e.Validator = validator.New()
	rConf := middleware.RecoverConfig{
		DisablePrintStack: true,
		DisableStackAll:   true,
	}
	if config.Config.Debug {
		rConf = middleware.RecoverConfig{
			DisableStackAll: true,
			LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
				logger.Logger.Err(err).Str("url", c.Request().URL.RequestURI()).Send()
				fmt.Println(string(stack))
				return nil
			},
		}
	}
	e.Use(middleware.RecoverWithConfig(rConf))
	// TODO: migrate to redis-based rate limiter
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(config.Config.RateLimit))))
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middlewares.Authentication(config.Config.PrivateKey))
	for route, group := range routes.Map {
		g := e.Group(route)
		// TODO: allow middlewares
		for r, sub := range group {
			for _, m := range sub.Methods {
				method := m.ToString()
				handler := createHandler(sub.Handler)
				if route == "/" {
					// add handler for root routes
					e.Add(method, r, handler)
				} else if r == "/" {
					// handle apex routes
					e.Add(method, route, handler)
				} else {
					g.Add(method, r, handler)
				}
			}
		}
	}
	return e
}
