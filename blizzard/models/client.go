package models

import (
	"backend/blizzard/pb"
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type (
	IglooClient struct {
		pb.DRPCIglooClient
		Address string `json:"-"`
	}
	IglooClusters = map[string]IglooClient
)

func (client *IglooClient) Ping(ctx context.Context) (health *pb.IglooHealth, responseTime float64) {
	start := time.Now()
	if client.DRPCIglooClient == nil {
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
