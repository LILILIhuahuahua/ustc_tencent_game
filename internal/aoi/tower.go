package aoi

import (
	"fmt"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
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

func (this *Tower) HeroEnter(hero *model.Hero) {
	this.heros.Store(hero.ID, hero)
	fmt.Printf("hero加入了新的灯塔\n")
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
