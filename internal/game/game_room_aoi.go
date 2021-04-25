package game

import (
	"fmt"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/aoi"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
)

func (gameRoom *GameRoom) NotifyHeroView(changeHero *model.Hero, notifyType int32, tower *aoi.Tower) {
	var needToNotify []*model.Hero
	tower.Heros.Range(func(k, v interface{}) bool {
		needToNotify = append(needToNotify, v.(*model.Hero))
		return true
	})
	for _, hero := range needToNotify {
		if notifyType == configs.Enter {
			fmt.Printf("%d向%d发送了进入视野的信息  ---enter", changeHero.Session.Id, hero.Session.Id)
		}
		if notifyType == configs.Leave {
			fmt.Printf("%d向%d发送了liking视野的信息 ---leave", changeHero.Session.Id, hero.Session.Id)
		}
		gameRoom.SendHeroViewNotify(changeHero, hero, notifyType)
	}
}

func (gameRoom *GameRoom) NotifyHeroPropMsg(enterHero *model.Hero) {
	towersId := tools.GetOtherTowers(enterHero.TowerId)
	towersId = append(towersId, enterHero.TowerId)
	gameRoom.SendHeroPropGlobalInfoNotify(towersId, enterHero.Session)
}
