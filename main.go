package main

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/db"
	"blizzard/blizzard/judge"
	"blizzard/blizzard/logger"
	"blizzard/blizzard/oauth"
	"blizzard/blizzard/server"
	"blizzard/blizzard/server/utils"
)

func init() {
	logger.Init()
}

func init() {
	config.Load()
}

func init() {
	db.Init()
	judge.Init()
	oauth.Init()
}

func main() {
	utils.Listen(server.CreateServer())
}
