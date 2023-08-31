package main

import (
	"blizzard/blizzard/judge"
	"blizzard/blizzard/server/http"
	"blizzard/blizzard/server/utils"
)

func main() {
	go judge.ResponseObserver.Work()
	utils.Listen(http.CreateServer())
}
