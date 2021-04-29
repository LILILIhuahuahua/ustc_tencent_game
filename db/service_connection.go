package db

import (
	"errors"
	"github.com/LILILIhuahuahua/ustc_tencent_game/db/databaseGrpc"
	"google.golang.org/grpc"
)

type ServiceConnection struct {
	conn *grpc.ClientConn
	accountService databaseGrpc.AccountServiceClient
	playerService databaseGrpc.PlayerServiceClient
}

var serviceConnection *ServiceConnection

func InitConnection(address string) error {
	serviceConnection = &ServiceConnection{}
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	serviceConnection.conn = conn
	serviceConnection.playerService = databaseGrpc.NewPlayerServiceClient(conn)
	serviceConnection.accountService = databaseGrpc.NewAccountServiceClient(conn)
	return nil
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