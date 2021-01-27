package main

import (
	"fmt"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/game"
)

func main() {
	fmt.Println("ustc_tencent_game_server started!")
	addr := "127.0.0.1:8888"
	s := game.NewGameStarter(addr)
	s.Boot()
	//b, err := framework.NewGameRoom(addr)
	//if err != nil {
	//
	//}
	//err = b.Serv()
	//if err != nil {
	//
	//}
}
