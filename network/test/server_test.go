package test

import (
	"github.com/LILILIhuahuahua/ustc_tencent_game/network"
	"testing"
)

var (
	addr = "127.0.0.1:12345"
)

// TestNewKcpServer create a server for test
func TestNewKcpServer(t *testing.T) {
	b, _ := network.NewBroadcaster(addr)
	b.Serv()
}
