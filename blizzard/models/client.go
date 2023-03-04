package models

import (
	"backend/blizzard/pb"
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type PolarClient struct {
	pb.DRPCPolarClient
	Address string `json:"-"`
}

func (client *PolarClient) Ping(ctx context.Context) (health *pb.PolarHealth, responseTime float64) {
	start := time.Now()
	if client.DRPCPolarClient == nil {
		return nil, -1
	}
	health, e := client.Health(ctx, nil)
	logrus.Error(e)
	if e != nil {
		return nil, -1
	}
	responseTime = float64(time.Since(start).Nanoseconds()) / 1e6
	return
}
