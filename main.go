package main

import (
	"backend/blizzard"
	"backend/blizzard/core"
)

func main() {
	c := core.ReadConfig()
	blizzard.CreateServer(c).Listen()
}
