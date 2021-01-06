package main

import (
	"fmt"
	"github.com/LILILIhuahuahua/ustc_tencent_game/network"
)

func main() {
	fmt.Println("ustc_tencent_game_server started!")
	addr := "127.0.0.1:8888"
	b, err := network.NewBroadcaster(addr)
	if err != nil {

	}
	err = b.Serv()
	if err != nil {

	}
}
