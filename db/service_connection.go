package db

import (
	"context"
	"errors"
	"github.com/LILILIhuahuahua/ustc_tencent_game/db/databaseGrpc"
	"google.golang.org/grpc"
	"log"
	"time"
)

type ServiceConnection struct {
	conn           *grpc.ClientConn
	accountService databaseGrpc.AccountServiceClient
	playerService  databaseGrpc.PlayerServiceClient
}

var serviceConnection *ServiceConnection

func InitConnection(address string) {
	serviceConnection = &ServiceConnection{}
	ctx, _ := context.WithDeadline(context.TODO(), time.Now().Add(5*time.Second))
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln("[dbProxy] dbProxy初始化失败, 失败原因:", err)
	}
	serviceConnection.conn = conn
	serviceConnection.playerService = databaseGrpc.NewPlayerServiceClient(conn)
	serviceConnection.accountService = databaseGrpc.NewAccountServiceClient(conn)
	log.Println("[dbProxy] dbProxy初始化成功")
}

func GetPlayerService() (databaseGrpc.PlayerServiceClient, error) {
	if serviceConnection == nil {
		return nil, errors.New("没有初始化connection")
	}
	return serviceConnection.playerService, nil
}

func GetAccountService() (databaseGrpc.AccountServiceClient, error) {
	if serviceConnection == nil {
		return nil, errors.New("没有初始化connection")
	}
	return serviceConnection.accountService, nil
}
