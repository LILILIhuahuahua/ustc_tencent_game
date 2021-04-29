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
		DBUser string
		DBPassword string
		DBHost string
		DBPort string
	)

	flag.StringVar(&DBUser, "DBUser","","User name of database" )
	flag.StringVar(&DBPassword, "DBPassword", "","Password of database user")
	flag.StringVar(&DBHost,"Host","","IP address of database")
	flag.StringVar(&DBPort, "Port","","Port to connect to database")

	flag.Parse()
	configs.MongoURI = "mongodb://"+ DBUser + ":"+DBPassword + "@"+ DBHost + ":" + DBPort + "/" + configs.DBName
}

func main() {
	// Initialize mongodb
	initDB()
	if configs.MongoURI == "" {
		log.Fatalln("MongoURI is nil")
	}
	log.Println("Initialize MongoURI to ",configs.MongoURI)

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
