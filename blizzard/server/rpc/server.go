package rpc

import (
	"blizzard/blizzard/logger"
	"blizzard/blizzard/pb/blizzard"
	"context"
	"go.arsenm.dev/drpc/muxserver"
	"log"
	"net"
	"storj.io/drpc/drpcmux"
)

func StartRpc() {
	mux := drpcmux.New()
	e := blizzard.DRPCRegisterBlizzard(mux, new(Blizzard))
	if e != nil {
		panic(e)
	}
	lis, err := net.Listen("tcp", ":2345")
	if err != nil {
		panic(err)
	}
	server := muxserver.New(createMiddleware("test2")(mux))
	logger.Logger.Info().Msgf("RPC listening on %s", ":2345")
	log.Fatalln(server.Serve(context.Background(), lis))
}
