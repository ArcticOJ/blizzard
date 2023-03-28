package judge

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/pb"
	"context"
	"go.arsenm.dev/drpc/muxconn"
	"net"
	"time"
)

var Igloo IglooCluster

func makeClients(addrs map[string]string) (cluster IglooCluster) {
	cluster = make(IglooCluster)
	for name, addr := range addrs {
		cluster[name] = IglooClient{
			DRPCIglooClient: nil,
			Address:         addr,
		}
		dial, e := net.DialTimeout("tcp", addr, time.Second*3)
		if e != nil {
			continue
		}
		conn, _ := muxconn.New(dial)
		if client, ok := cluster[name]; ok {
			client.DRPCIglooClient = pb.NewDRPCIglooClient(conn)
			cluster[name] = client
		}
	}
	return
}

func TrySelectClient(ctx context.Context) (ok bool, client *IglooClient) {
	for name := range Igloo {
		if client, ok := Renew(ctx, Igloo, name); ok {
			return true, client
		}
	}
	return false, nil
}

func Init() {
	Igloo = makeClients(config.Config.Judges)
}
