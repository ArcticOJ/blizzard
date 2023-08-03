package utils

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/logger"
	"fmt"
	"github.com/labstack/echo/v4"
	"time"
)

var BootTimestamp = time.Now()

func Uptime() int64 {
	return int64(time.Since(BootTimestamp).Round(time.Second).Seconds())
}

func Listen(router *echo.Echo) {
	router.HideBanner = true
	router.HidePort = true
	router.IPExtractor = echo.ExtractIPFromRealIPHeader()
	addr := fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port)
	logger.Logger.Info().Msgf("starting server on %s", addr)
	logger.Logger.Fatal().Err(router.Start(addr)).Send()
}
