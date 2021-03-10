package aoi

import (
	"fmt"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/prop"
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
	//fmt.Printf("hero加入了新的灯塔\n")
}

func (this *Tower) PropEnter(prop *prop.Prop) {
	//fmt.Printf("存入的prop ID: %d \n", prop.ID())
	this.props.Store(prop.ID(), prop)
}

func (this *Tower) HeroLeave(hero *model.Hero) { // 后期做优化
	this.heros.Delete(hero.ID)
}

func (this *Tower) PropLeave(prop *model.Prop) {
	this.props.Delete(prop.Id)
}

func (this *Tower) GetHeros() []*model.Hero {
	var heros []*model.Hero
	this.heros.Range(func(k, v interface{}) bool {
		heros = append(heros, v.(*model.Hero))
		return true
	})
	return heros
}

func (this *Tower) GetProps() []*prop.Prop {
	var props []*prop.Prop
	this.props.Range(func(k, v interface{}) bool {
		prop := v.(*prop.Prop)
		//if prop.Status() == configs.PropStatusLive {
			props = append(props, prop)
		//}
		return true
	})
	return props
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
		if notifyType == configs.Enter {
			fmt.Printf("%d向%d发送了进入视野的信息  ---enter", changeHero.Session.Id, hero.Session.Id)
		}
		if notifyType == configs.Leave {
			fmt.Printf("%d向%d发送了liking视野的信息 ---leave", changeHero.Session.Id, hero.Session.Id)
		}
		callback(changeHero, hero, notifyType)
	}
}

func (this *Tower) NotifyHeroPropMsg(callback func([]int32, *framework.BaseSession)) { // 通知灯塔内所有玩家附近的道具和玩家信息
	var needToNotify []*model.Hero
	this.heros.Range(func(k, v interface{}) bool {
		needToNotify = append(needToNotify, v.(*model.Hero))
		return true
	})
	towerIds := tools.GetOtherTowers(this.id)
	towerIds = append(towerIds, this.id)
	for _, hero := range needToNotify {
		callback(towerIds, hero.Session)
	}
}
