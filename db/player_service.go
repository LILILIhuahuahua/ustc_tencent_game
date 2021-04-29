package db

import (
	"context"
	"github.com/LILILIhuahuahua/ustc_tencent_game/db/databaseGrpc"
	"time"
)

var playerService databaseGrpc.PlayerServiceClient

func getService() (databaseGrpc.PlayerServiceClient, error) {
	if playerService == nil {
		service, err := GetPlayerService()
		if err != nil {
			return nil, err
		}
		playerService = service
		return service, nil
	}
	return playerService, nil
}

func PlayerUpdateHighestScoreByPlayerId(playerId, newScore int32) error {
	service, err  := getService()
	if err != nil {
		return err
	}
	ctx, cancel:= context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	requestMsg := &databaseGrpc.PlayerUpdateHighestScoreByPlayerIdRequest{
		PlayerId:     playerId,
		HighestScore: newScore,
	}
	_, err = service.PlayerUpdateHighestScoreByPlayerId(ctx, requestMsg)
	if err != nil {
		return err
	}
	return nil
}