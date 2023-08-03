package main

import (
	"blizzard/blizzard/cache"
	"blizzard/blizzard/config"
	"blizzard/blizzard/db"
	"blizzard/blizzard/judge"
	"blizzard/blizzard/logger"
	"blizzard/blizzard/oauth"
	"blizzard/blizzard/server/http"
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
	cache.Init()
}

func main() {
	go judge.ResultWatcher.Watch()
	utils.Listen(http.CreateServer())
}
