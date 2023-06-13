package main

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/db"
	"blizzard/blizzard/judge"
	"blizzard/blizzard/logger"
	"blizzard/blizzard/oauth"
	"blizzard/blizzard/server/http"
	"blizzard/blizzard/server/rpc"
	"blizzard/blizzard/server/utils"
	"blizzard/blizzard/utils/crypto"
)

func init() {
	logger.Init()
}

func init() {
	config.Load()
}

func init() {
	crypto.Init()
	db.Init()
	judge.Init()
	oauth.Init()
}

func main() {
	go rpc.StartRpc()
	utils.Listen(http.CreateServer())
}
