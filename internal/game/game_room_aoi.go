package game

import (
	"fmt"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/aoi"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
)

func (gameRoom *GameRoom) NotifyHeroView(changeHero *model.Hero, notifyType int32, tower *aoi.Tower) {
	var heroNeededToNotify []*model.Hero
	tower.Heros.Range(func(k, v interface{}) bool {
		heroNeededToNotify = append(heroNeededToNotify, v.(*model.Hero))
		return true
	})
	for _, hero := range heroNeededToNotify {
		if notifyType == configs.Enter {
			fmt.Printf("%d向%d发送了进入视野的信息  ---enter", changeHero.Session.Id, hero.Session.Id)
		}
		if notifyType == configs.Leave {
			fmt.Printf("%d向%d发送了liking视野的信息 ---leave", changeHero.Session.Id, hero.Session.Id)
		}
		gameRoom.SendHeroViewNotify(changeHero, hero, notifyType)
	}
}

func (gameRoom *GameRoom) NotifyHeroPropMsgToHero(enterHero *model.Hero) {
	towersIds := tools.GetOtherTowers(enterHero.TowerId)
	towersIds = append(towersIds, enterHero.TowerId)
	go gameRoom.SendHeroInTowerNotify(towersIds, enterHero.Session)
	go gameRoom.SendPropInTowerNotify(towersIds, enterHero.Session)
}

// 向灯塔内的所有小球广播信息
func (gameRoom *GameRoom) NotifyHeroPropMsg() {
	towers := gameRoom.GetTowers()
	for _, tower := range towers {
		towersId := tools.GetOtherTowers(tower.Id)
		towersId = append(towersId, tower.Id)
		heros := tower.GetHeros()
		for _, hero := range heros {
			go gameRoom.SendPropInTowerNotify(towersId, hero.Session)
			go gameRoom.SendHeroInTowerNotify(towersId, hero.Session)
		}
	}
}

func (gameRoom *GameRoom) NotifyEntityInfoChange(entityType int32, entityId int32, changedHero *model.Hero, changedProp *model.Prop) {
	var heroInfo *info.HeroInfo
	var itemInfo *info.ItemInfo
	var entityTowerId int32
	if changedHero != nil {
		heroInfo = info.NewHeroInfo(changedHero)
		entityTowerId = changedHero.TowerId
	}
	if changedProp != nil {
		itemInfo = info.NewItemInfo(changedProp)
		entityTowerId = changedProp.TowerId
	}
	towersIds := tools.GetOtherTowers(entityTowerId)
	towersIds = append(towersIds, entityTowerId)
	towers := gameRoom.GetTowers()
	var heroNeededToNotify []*model.Hero
	for _, towerId := range towersIds {
		towerHeros := towers[towerId].GetHeros()
		heroNeededToNotify = append(heroNeededToNotify, towerHeros...)
	}
	for _, hero := range heroNeededToNotify {
		gameRoom.SendEntityInfoChangeNotify(hero, entityType, entityId, heroInfo, itemInfo)
	}
}
