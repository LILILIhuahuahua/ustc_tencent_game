package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type accountCollection struct {}

type Account struct { // 里面的字段名一定要大写开头
	Name string
	LoginPassword string // 登录密码
	AccountAvatar string // 头像
	Level int64 // 当前等级
	Delete bool // 当前账号是否住校
	Region string // 用户的地区
	Phone string // 电话
	CreateAt int64
	UpdateAt int64
}

var AccountCollection = &accountCollection{}

func (this *accountCollection) getCollection() *mongo.Collection {
	return Mc.GetCollection("Account", nil)
}

func (this *accountCollection) InsertAccount(account *Account) (string, error) {
	collection := this.getCollection()
	insertResult, err := collection.InsertOne(context.TODO(), account)
	if err != nil {
		return "", err
	}
	return insertResult.InsertedID.(primitive.ObjectID).String(), nil
}

func (this *accountCollection) FindAccount(accountId string) (*Account, error) {
	collection := this.getCollection()
	accIdObject, err := primitive.ObjectIDFromHex(accountId)
	if err != nil {
		return nil, err
	}
	result := collection.FindOne(context.TODO(), bson.M{"_id": accIdObject})
	if result.Err() != nil {
		return nil, result.Err()
	}
	findAcc := &Account{}
	err = result.Decode(findAcc)
	if err != nil {
		return nil, err
	}
	return findAcc, nil
}

func (this *accountCollection) UpdateAccount(accountId string, account *Account) {
	collection := this.getCollection()
	accIdObject, err := primitive.ObjectIDFromHex(accountId)
	update := bson.M{
		"$set": bson.M{
			"phone": account.Phone,
		},
	}

	if err != nil {
		return
	}
	collection.UpdateByID(context.TODO(), accIdObject, update)
}