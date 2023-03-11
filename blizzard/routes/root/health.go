package root

import (
	"backend/blizzard/build"
	"backend/blizzard/models"
	"backend/blizzard/pb"
	"github.com/labstack/echo/v4"
	"time"
)

type JudgeInfo struct {
	Name    string             `json:"name"`
	IsAlive bool               `json:"isAlive"`
	Latency float64            `json:"latency"`
	Uptime  int64              `json:"uptime"`
	Version string             `json:"version"`
	Specs   *pb.Specifications `json:"specs,omitempty"`
}

func Health(ctx *models.Context) models.Response {
	ctx.Response().Header().Add("Timing-Allow-Origin", "*")
	var judgesInfo []JudgeInfo
	for name, client := range ctx.Server.Igloo {
		health, respTime := client.Ping(ctx.Request().Context())
		if health == nil || client.DRPCIglooClient == nil {
			judgesInfo = append(judgesInfo, JudgeInfo{
				Name:    name,
				IsAlive: false,
				Latency: -1,
				Uptime:  -1,
				Version: "Unknown",
				Specs:   nil,
			})
			continue
		}
		now := time.Now().UTC()
		judgesInfo = append(judgesInfo, JudgeInfo{
			Name:    name,
			IsAlive: true,
			Latency: respTime,
			Uptime:  now.Unix() - health.BootTimestamp.Seconds,
			Version: health.Version,
			Specs:   health.Specs,
		})
	}
	if judgesInfo == nil {
		judgesInfo = []JudgeInfo{}
	}
	return ctx.Respond(echo.Map{
		"version": build.Version,
		"uptime":  ctx.Server.Uptime(),
		"judges":  judgesInfo,
	})
}
