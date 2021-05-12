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

func (t *Tower) GetId() int32 {
	return t.Id
}

func (t *Tower) HeroEnter(hero *model.Hero) {
	t.Heros.Store(hero.ID, hero)
}

func (t *Tower) PropEnter(prop *model.Prop) {
	//fmt.Printf("存入的prop ID: %d \n", prop.ID())
	t.Props.Store(prop.Id, prop)
}

func (t *Tower) HeroLeave(hero *model.Hero) { // 后期做优化
	t.Heros.Delete(hero.ID)
}

func (t *Tower) PropLeave(prop *model.Prop) {
	t.Props.Delete(prop.Id)
}

func (t *Tower) GetHeros() []*model.Hero {
	var heros []*model.Hero
	t.Heros.Range(func(k, v interface{}) bool {
		heros = append(heros, v.(*model.Hero))
		return true
	})
	return heros
}

func (t *Tower) GetProps() []*model.Prop {
	var props []*model.Prop
	t.Props.Range(func(k, v interface{}) bool {
		prop := v.(*model.Prop)
		//if prop.Status() == configs.PropStatusLive {
		props = append(props, prop)
		//}
		return true
	})
	return props
}

func (t *Tower) NotifyHeroPropMsg(callback func([]int32, *framework.BaseSession)) { // 通知灯塔内所有玩家附近的道具和玩家信息
	var needToNotify []*model.Hero
	t.Heros.Range(func(k, v interface{}) bool {
		needToNotify = append(needToNotify, v.(*model.Hero))
		return true
	})
	towerIds := tools.GetOtherTowers(t.Id)
	towerIds = append(towerIds, t.Id)
	for _, hero := range needToNotify {
		callback(towerIds, hero.Session)
	}
}
