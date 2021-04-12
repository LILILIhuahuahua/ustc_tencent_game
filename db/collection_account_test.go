package db

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestInsertAccount(t *testing.T) {
	InitClient()

	id, err := AccountCollection.InsertAccount(&Account{
		Name:          "song",
		LoginPassword: "ttt",
		AccountAvatar: "www.baidu.com",
		Level:         1,
		Delete:        false,
		Region:        "China",
		Phone:         "17376515082",
		CreateAt:      time.Now().UnixNano(),
		UpdateAt:      0,
	})
	if err != nil {
		log.Println("插入时发生了错误", err)
		return
	}
	log.Println("获取的ID为", id)
}

func TestFindAccount(t *testing.T) {
	InitClient()

	acc, err := AccountCollection.FindAccount("605b7267be255a7618e38d6a")
	if err != nil {
		t.Error("查找时发生了错误", err)
	}
	fmt.Println("获取的acc为", acc)
}
