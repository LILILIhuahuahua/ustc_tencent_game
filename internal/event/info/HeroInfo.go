package info

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
)

type HeroInfo struct {
	framework.BaseEvent
	ID int64
	Status int32
	Speed float32
	Size float32
	HeroPosition CoordinateXYInfo
	HeroDirection CoordinateXYInfo
}

func (h *HeroInfo)FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.HeroMsg)
	h.ID = pbMsg.GetHeroId()
	h.Status = int32(pbMsg.GetHeroStatus())
	h.Speed = pbMsg.GetHeroSpeed()
	h.Size = pbMsg.HeroSize
	  pos := CoordinateXYInfo{}
	  pos.FromMessage(pbMsg.GetHeroPosition())
	h.HeroPosition = pos
	  dict := CoordinateXYInfo{}
	  dict.FromMessage(pbMsg.GetHeroDirection())
	h.HeroDirection = dict
}

func (h *HeroInfo)CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.HeroMsg)
	pos := CoordinateXYInfo{}
	dict := CoordinateXYInfo{}
	dict.FromMessage(pbMsg.GetHeroDirection())
	pos.FromMessage(pbMsg.GetHeroPosition())
	return &HeroInfo{
		ID: pbMsg.GetHeroId(),
		Status: int32(pbMsg.GetHeroStatus()),
		Speed: pbMsg.GetHeroSpeed(),
		Size: pbMsg.GetHeroSize(),
		HeroPosition: pos,
		HeroDirection: dict,
	}
}

func (h *HeroInfo)ToMessage() interface{} {
	pbMsg := &pb.HeroMsg{
		HeroId: h.ID,
		//HeroStatus: ,
		HeroSpeed: h.Speed,
		HeroSize: h.Size,
		HeroPosition: h.HeroPosition.ToMessage().(*pb.CoordinateXY),
		HeroDirection: h.HeroDirection.ToMessage().(*pb.CoordinateXY),
	}
	//todo:
	return pbMsg
}
