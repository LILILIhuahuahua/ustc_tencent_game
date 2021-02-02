package info

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
)

type ConnectInfo struct {
	framework.BaseEvent  	//基础消息类作为父类
	Ip                  string
	Port                int32
}

func (c *ConnectInfo)FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.ConnectMsg)
	c.Ip = pbMsg.Ip
	c.Port = pbMsg.Port
}

func (c *ConnectInfo)CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.ConnectMsg)
	return &ConnectInfo{
		Ip:   pbMsg.Ip,
		Port: pbMsg.Port,
	}
}

func (c *ConnectInfo)ToMessage() interface{} {
	return &pb.ConnectMsg{
		Ip: c.Ip,
		Port: c.Port,
	}
}