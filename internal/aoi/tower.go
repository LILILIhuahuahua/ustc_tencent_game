package aoi

import (
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"sync"
)

type Tower struct {
	id    int32
	heros sync.Map
	props sync.Map
}

func InitTower(id int32) *Tower {
	return &Tower{id: id}
}

func (this *Tower) GetId() int32 {
	return this.id
}

func (this *Tower) HeroEnter(hero *model.Hero, callback func([]int32, *framework.BaseSession)) {
	this.heros.Store(hero.ID, hero)
	towerIds := tools.GetOtherTowers(this.id)
	towerIds = append(towerIds, this.id)
	callback(towerIds, hero.Session)
}

func (this *Tower) HeroLeave(hero *model.Hero) {
	this.heros.Delete(hero.ID)
}

func (this *Tower) GetHeros() []*model.Hero {
	var heros []*model.Hero
	this.heros.Range(func(k, v interface{}) bool {
		heros = append(heros, v.(*model.Hero))
		return true
	})
	return heros
}

func (this *Tower) NotifyHeroMsg(
	changeHero *model.Hero,
	notifyType int32,
	callback func(changeHero *model.Hero, hero *model.Hero, notifyType int32)) {

	var needToNotify []*model.Hero
	this.heros.Range(func(k, v interface{}) bool {
		needToNotify = append(needToNotify, v.(*model.Hero))
		return true
	})
	for _, hero := range needToNotify {
		callback(changeHero, hero, notifyType)
	}
}
