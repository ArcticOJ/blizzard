package utils

import (
	"net"
	"time"
)

func Ping(addr string) bool {
	dial, e := net.DialTimeout("tcp", addr, 1*time.Second)
	defer dial.Close()
	if e != nil {
		return false
	}
	return true
}
