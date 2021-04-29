package main

import (
	"flag"
	"fmt"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/game"
	"log"
)

func initDB() {
	var (
		dbProxyPort string
		dbProxyHost string
	)

	flag.StringVar(&dbProxyHost, "DBProxyHost", "", "Host addr of dbproxy")
	flag.StringVar(&dbProxyPort, "DBProxyPort", "", " Port of dbproxy")
	flag.Parse()
	configs.DBProxyAddr = dbProxyHost + ":" + dbProxyPort
}

func main() {
	initDB()
	if configs.DBProxyAddr == "" {
		log.Fatalln("DBProxy addr is nil")
	}
	log.Println("Initialize MongoURI to ", configs.DBProxyAddr)

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
