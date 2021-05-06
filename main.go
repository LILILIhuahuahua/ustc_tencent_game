package main

import (
	"flag"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/db"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/game"
	"log"
	"net/http"

	_ "github.com/mkevac/debugcharts"
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
	log.Println("Initialize DBProxyAddr to", configs.DBProxyAddr)
	go db.InitConnection(configs.DBProxyAddr)

	// visual representation of goroutine
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	log.Println("[USTC-Tencent]Game Server Started!")
	s := game.NewGameStarter(configs.ServerAddr)
	s.Boot()
}
