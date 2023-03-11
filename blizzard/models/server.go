package models

import (
	"backend/blizzard/core"
	"backend/blizzard/pb"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"go.arsenm.dev/drpc/muxconn"
	"net"
	"time"
)

type Server struct {
	*echo.Echo
	Database      *bun.DB
	BootTimestamp time.Time
	Config        *core.Config
	Igloo         IglooClusters
}

func (s Server) Uptime() int64 {
	return int64(time.Since(s.BootTimestamp).Round(time.Second).Seconds())
}

func (s Server) TrySelectClient() (ok bool, client *IglooClient) {
	for name, cluster := range s.Igloo {
		// TODO: reconnect once a client is expired
		if cluster.DRPCIglooClient == nil {
			dial, e := net.DialTimeout("tcp", cluster.Address, time.Second*3)
			if e != nil {
				logrus.Error(e)
				continue
			}
			conn, e := muxconn.New(dial)
			if e == nil {
				cluster.DRPCIglooClient = pb.NewDRPCIglooClient(conn)
				s.Igloo[name] = cluster
				return true, &cluster
			}
		}
	}
	return false, nil
}

func (s Server) Listen() {
	s.HideBanner = true
	s.Logger.Fatal(s.Start(fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)))
}
