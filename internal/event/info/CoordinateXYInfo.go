package info

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
)

type CoordinateXYInfo struct {
	framework.BaseEvent
	CoordinateX float32
	CoordinateY float32
}

func (c *CoordinateXYInfo) FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.CoordinateXY)
	if nil != pbMsg {
		c.CoordinateX = pbMsg.CoordinateX
		c.CoordinateY = pbMsg.CoordinateY
	}
}

func (c *CoordinateXYInfo) CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.CoordinateXY)
	return &CoordinateXYInfo{
		CoordinateX: pbMsg.CoordinateX,
		CoordinateY: pbMsg.CoordinateY,
	}
}

func (c *CoordinateXYInfo) ToMessage() interface{} {
	return &pb.CoordinateXY{
		CoordinateX: c.CoordinateX,
		CoordinateY: c.CoordinateY,
	}
}
