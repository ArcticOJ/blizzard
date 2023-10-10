package middlewares

import (
	"github.com/ArcticOJ/blizzard/v0/cache/stores"
	"github.com/labstack/echo/v4"
	"strconv"
)

func RateLimit() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			allowed, totalLimit, remaining, retryAfter, nextReset := stores.RateLimiter.Limit(c.Request().Context(), c.RealIP())
			if !allowed {
				return c.JSON(429, echo.Map{
					"code":       429,
					"message":    "Rate limit exceeded",
					"resetAfter": retryAfter,
					"nextReset":  nextReset,
				})
			}
			h := c.Response().Header()
			h.Set("X-RateLimit-Limit", strconv.FormatInt(totalLimit, 10))
			h.Set("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
			h.Set("Retry-After", strconv.FormatInt(retryAfter, 10))
			h.Set("X-RateLimit-Reset", strconv.FormatInt(nextReset, 10))
			return next(c)
		}
	}
}
