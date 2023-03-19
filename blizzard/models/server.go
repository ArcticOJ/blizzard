package models

import (
	"backend/blizzard/core"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"time"
)

type Server struct {
	*echo.Echo
	Database      *bun.DB
	BootTimestamp time.Time
	Config        *core.Config
	Igloo         IglooCluster
	Logger        zerolog.Logger
}

func (s Server) Uptime() int64 {
	return int64(time.Since(s.BootTimestamp).Round(time.Second).Seconds())
}

func (s Server) TrySelectClient(ctx context.Context) (ok bool, client *IglooClient) {
	for name := range s.Igloo {
		if client, ok := Renew(ctx, s.Igloo, name); ok {
			return true, client
		}
	}
	return false, nil
}

func (s Server) Listen() {
	s.HideBanner = true
	s.Logger.Fatal().Err(s.Start(fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)))
}
