package judge

import (
	"blizzard/blizzard/config"
	models "blizzard/blizzard/models"
	"blizzard/blizzard/pb/igloo"
	"context"
	"go.arsenm.dev/drpc/muxconn"
	"net"
	"storj.io/drpc/drpcmetadata"
	"time"
)

var Igloo IglooCluster

func makeClients(judges map[string]models.Judge) (cluster IglooCluster) {
	cluster = make(IglooCluster)
	for name, judge := range judges {
		cluster[name] = IglooClient{
			DRPCIglooClient: nil,
			Address:         judge.Address,
		}
		dial, e := net.DialTimeout("tcp", judge.Address, time.Second*3)
		if e != nil {
			continue
		}
		conn, _ := muxconn.New(dial)
		if client, ok := cluster[name]; ok {
			client.DRPCIglooClient = igloo.NewDRPCIglooClient(conn)
			cluster[name] = client
		}
	}
	return
}

func PickClient(ctx context.Context) (ok bool, client *IglooClient) {
	for name := range Igloo {
		if client, ok := Renew(KeyContext(ctx), Igloo, name); ok {
			return true, client
		}
	}
	return false, nil
}

func KeyContext(ctx context.Context) context.Context {
	return drpcmetadata.Add(ctx, "key", "test")
}

func Init() {
	Igloo = makeClients(config.Config.Judges)
}
