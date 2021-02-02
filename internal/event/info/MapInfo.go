package info

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
)

type MapInfo struct {
	framework.BaseEvent
	XMin float32
	XMax float32
	YMin float32
	YMax float32
}

func (m *MapInfo)FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.MapMsg)
	m.XMax = pbMsg.XMax
	m.XMin = pbMsg.XMin
	m.YMax = pbMsg.YMax
	m.YMin = pbMsg.YMin
}

func (m *MapInfo)CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.MapMsg)
	return &MapInfo{
		XMax: pbMsg.XMax,
		XMin: pbMsg.XMin,
		YMax: pbMsg.YMax,
		YMin: pbMsg.YMin,
	}
}

func (m *MapInfo)ToMessage() interface{} {
	return &pb.MapMsg{
		XMax: m.XMax,
		XMin: m.XMin,
		YMax: m.YMax,
		YMin: m.YMin,
	}
}
