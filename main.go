package main

import (
	"fmt"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/game"
)

func main() {
	fmt.Println("ustc_tencent_game_server started!")
	s := game.NewGameStarter(configs.ServerAddr)
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
