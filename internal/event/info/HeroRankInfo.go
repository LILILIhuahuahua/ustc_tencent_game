package info

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
)

type HeroRankInfo struct {
	framework.BaseEvent
	HeroID    int32
	HeroName  string
	HeroRank  int32
	HeroScore int32
}

func NewHeroRankInfo(hero *model.Hero) *HeroRankInfo{
	return &HeroRankInfo{
		HeroID:    hero.ID,
		HeroName:  hero.Name,
		HeroRank:  hero.Rank,
		HeroScore: hero.Score,
	}
}

func (h *HeroRankInfo) FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.HeroRankMsg)
	h.HeroID = pbMsg.GetHeroId()
	h.HeroName = pbMsg.GetHeroName()
	h.HeroRank = pbMsg.GetHeroRank()
	h.HeroScore = pbMsg.GetHeroScore()
}

func (h *HeroRankInfo) CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.HeroRankMsg)
	return &HeroRankInfo{
		HeroID:   pbMsg.GetHeroId(),
		HeroName: pbMsg.GetHeroName(),
		HeroRank: pbMsg.GetHeroRank(),
		HeroScore: pbMsg.GetHeroScore(),
	}
}

func (h *HeroRankInfo) ToMessage() interface{} {
	pbMsg := &pb.HeroRankMsg{
		HeroId: h.HeroID,
		HeroName: h.HeroName,
		HeroRank: h.HeroRank,
		HeroScore: h.HeroScore,
	}
	return pbMsg
}

