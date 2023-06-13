package root

import (
	"blizzard/blizzard/build"
	"blizzard/blizzard/judge"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/pb/igloo"
	"blizzard/blizzard/server/utils"
	"github.com/labstack/echo/v4"
	"time"
)

type JudgeInfo struct {
	Name    string  `json:"name"`
	IsAlive bool    `json:"isAlive"`
	Latency float64 `json:"latency"`
	Uptime  int64   `json:"uptime"`
	*igloo.InstanceSpecification
}

func Health(ctx *extra.Context) models.Response {
	ctx.Response().Header().Add("Timing-Allow-Origin", "*")
	var judgesInfo []JudgeInfo
	for name := range judge.Igloo {
		client, ok := judge.Renew(ctx.Request().Context(), judge.Igloo, name)
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
		"uptime":  utils.Uptime(),
		"judges":  judgesInfo,
	})
}
