package aoi

import (
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"sync"
)

type Tower struct {
	Id    int32
	Heros sync.Map
	Props sync.Map
}

func InitTower(id int32) *Tower {
	return &Tower{Id: id}
}

func (this *Tower) GetId() int32 {
	return this.Id
}

func (this *Tower) HeroEnter(hero *model.Hero) {
	this.Heros.Store(hero.ID, hero)
}

func (this *Tower) PropEnter(prop *model.Prop) {
	//fmt.Printf("存入的prop ID: %d \n", prop.ID())
	this.Props.Store(prop.Id, prop)
}

func (this *Tower) HeroLeave(hero *model.Hero) { // 后期做优化
	this.Heros.Delete(hero.ID)
}

func (this *Tower) PropLeave(prop *model.Prop) {
	this.Props.Delete(prop.Id)
}

func (this *Tower) GetHeros() []*model.Hero {
	var heros []*model.Hero
	this.Heros.Range(func(k, v interface{}) bool {
		heros = append(heros, v.(*model.Hero))
		return true
	})
	return heros
}

func (this *Tower) GetProps() []*model.Prop {
	var props []*model.Prop
	this.Props.Range(func(k, v interface{}) bool {
		prop := v.(*model.Prop)
		//if prop.Status() == configs.PropStatusLive {
			props = append(props, prop)
		//}
		return true
	})
	return props
}

func (this *Tower) NotifyHeroPropMsg(callback func([]int32, *framework.BaseSession)) { // 通知灯塔内所有玩家附近的道具和玩家信息
	var needToNotify []*model.Hero
	this.Heros.Range(func(k, v interface{}) bool {
		needToNotify = append(needToNotify, v.(*model.Hero))
		return true
	})
	towerIds := tools.GetOtherTowers(this.Id)
	towerIds = append(towerIds, this.Id)
	for _, hero := range needToNotify {
		callback(towerIds, hero.Session)
	}
}
