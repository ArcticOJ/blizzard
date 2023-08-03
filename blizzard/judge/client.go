package judge

import (
	"blizzard/blizzard/pb"
	"context"
	"go.arsenm.dev/drpc/muxconn"
	"net"
	"time"
)

type (
	Client struct {
		rpc        pb.DRPCIglooClient
		Name       string                    `json:"name"`
		Specs      *pb.InstanceSpecification `json:"specs"`
		privateKey string
		address    string
	}
)

func checkAlive(client *Client) bool {
	if client.rpc != nil {
		alive, e := client.rpc.Alive(keyedContext(context.Background(), client), nil)
		if e != nil {
			_ = client.rpc.DRPCConn().Close()
		} else if alive.Value {
			return true
		}
	}
	dial, e := net.DialTimeout("tcp", client.address, time.Millisecond*300)
	if e != nil {
		return false
	}
	conn, e := muxconn.New(dial)
	if e == nil {
		client.rpc = pb.NewDRPCIglooClient(conn)
		return true
	}
	return false
}

func (client *Client) Ping(ctx context.Context) (bool, float64) {
	start := time.Now()
	if checkAlive(client) {
		alive, e := client.rpc.Alive(keyedContext(ctx, client), nil)
		if e != nil {
			return false, -1
		}
		return alive.Value, float64(time.Since(start).Nanoseconds()) / 1e6
	}
	return false, -1
}
