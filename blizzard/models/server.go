package models

import (
	"backend/blizzard/core"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"time"
)

type Server struct {
	*echo.Echo
	Database      *bun.DB
	BootTimestamp time.Time
	Config        *core.Config
	Polar         map[string]PolarClient
}

func (s Server) Uptime() int64 {
	return int64(time.Since(s.BootTimestamp).Round(time.Second).Seconds())
}

func (s Server) Listen() {
	s.HideBanner = true
	s.Logger.Fatal(s.Start(fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)))
}
