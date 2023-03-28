package judge

import (
	"blizzard/blizzard/pb"
	"context"
	"go.arsenm.dev/drpc/muxconn"
	"net"
	"time"
)

type (
	IglooClient struct {
		pb.DRPCIglooClient
		Address string `json:"-"`
	}
	IglooCluster = map[string]IglooClient
)

func Renew(ctx context.Context, cluster IglooCluster, name string) (*IglooClient, bool) {
	if client, ok := cluster[name]; ok {
		if client.DRPCIglooClient != nil {
			alive, e := client.Alive(ctx, nil)
			if e != nil {
				_ = client.DRPCConn().Close()
			} else if alive.GetValue() {
				return &client, true
			}
		}
		dial, e := net.DialTimeout("tcp", client.Address, time.Millisecond*200)
		if e != nil {
			return nil, false
		}
		conn, e := muxconn.New(dial)
		if e == nil {
			client.DRPCIglooClient = pb.NewDRPCIglooClient(conn)
			cluster[name] = client
			return &client, true
		}
	}
	return nil, false
}

func (client *IglooClient) Ping(ctx context.Context) (specs *pb.InstanceSpecification, responseTime float64) {
	start := time.Now()
	specs, e := client.Specification(ctx, nil)
	if e != nil {
		return nil, -1
	}
	responseTime = float64(time.Since(start).Nanoseconds()) / 1e6
	return
}
