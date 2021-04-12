package db

import (
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	database *mongo.Database
	client *mongo.Client
}

var Mc *MongoClient

func InitClient() {
	cli := GetMgoCli()
	Mc = &MongoClient{
		database: cli.Database(configs.MongoDatabase),
		client: cli,
	}
}

func (this *MongoClient) GetCollection(collectionName string, options *options.CollectionOptions) *mongo.Collection {
	return this.database.Collection(collectionName, options)
}
