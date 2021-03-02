package aoi

import (
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"sync"
)

type Tower struct {
	id int32
	heros sync.Map
	props sync.Map
}

func (this *Tower) HeroEnter(hero *model.Hero) {
	this.heros.Store(hero.ID, hero)
}

func (this *Tower) HeroLeave(hero *model.Hero) {
	this.heros.Delete(hero.ID)
}

