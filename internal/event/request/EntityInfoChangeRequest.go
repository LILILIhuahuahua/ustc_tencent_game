package request

import pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"

type EntityInfoChangeRequest struct {
	eventType string
	heroId int32
	linkedId int32
	linkedType string
	//heroMsg
}

func (this *EntityInfoChangeRequest)FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.EntityInfoChangeRequest)
	this.eventType = pbMsg.EventType.String()
	this.heroId = pbMsg.HeroId
	this.linkedId = pbMsg.LinkedId
	this.linkedType = pbMsg.LinkedType.String()
}
