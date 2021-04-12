package db

import (
	"context"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var mgoCli *mongo.Client

func initEngine() {
	var err error
	clientOptions := options.Client().ApplyURI(configs.MongoUri)
	clientOptions.SetMaxPoolSize(configs.MongoPoolSize) // 连接池配置
	// 连接到MongoDB
	mgoCli, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = mgoCli.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func GetMgoCli() *mongo.Client {
	if mgoCli == nil {
		initEngine()
	}
	return mgoCli
}