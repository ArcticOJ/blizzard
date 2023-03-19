package root

import (
	"backend/blizzard/build"
	"backend/blizzard/models"
	"backend/blizzard/pb"
	"github.com/labstack/echo/v4"
	"time"
)

type JudgeInfo struct {
	Name    string  `json:"name"`
	IsAlive bool    `json:"isAlive"`
	Latency float64 `json:"latency"`
	Uptime  int64   `json:"uptime"`
	*pb.InstanceSpecification
}

func Health(ctx *models.Context) models.Response {
	ctx.Response().Header().Add("Timing-Allow-Origin", "*")
	var judgesInfo []JudgeInfo
	for name := range ctx.Igloo {
		client, ok := models.Renew(ctx.Request().Context(), ctx.Igloo, name)
		if !ok {
			judgesInfo = append(judgesInfo, JudgeInfo{
				Name:    name,
				IsAlive: false,
				Latency: -1,
				Uptime:  -1,
			})
			continue
		}
		specs, respTime := client.Ping(ctx.Request().Context())
		now := time.Now().UTC()
		info := JudgeInfo{
			Name:                  name,
			IsAlive:               true,
			Latency:               respTime,
			Uptime:                now.Unix() - specs.BootTimestamp.Seconds,
			InstanceSpecification: specs,
		}
		info.BootTimestamp = nil
		judgesInfo = append(judgesInfo, info)
	}
	if judgesInfo == nil {
		judgesInfo = []JudgeInfo{}
	}
	return ctx.Respond(echo.Map{
		"version": build.Version,
		"uptime":  ctx.Uptime(),
		"judges":  judgesInfo,
	})
}
