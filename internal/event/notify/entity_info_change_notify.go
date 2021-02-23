package notify

import (
	//pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"

	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	e "github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
)

//message EntityInfoChangeNotify {
//ENTITY_TYPE entityType = 1; //实体的类型
//int32 entityId = 2; //实体的Id
//HeroMsg HeroMsg = 3; //玩家的信息
//ItemMsg itemMsg = 4; //物体的信息
//}

type EntityInfoChangeNotify struct {
	framework.BaseEvent
	EntityType int32
	EntityId   int32
	HeroMsg    info.HeroInfo
	ItemMsg    info.ItemInfo
}

func (notify *EntityInfoChangeNotify) FromMessage(obj interface{}) {
	// 暂时不需要做
}

func (notify *EntityInfoChangeNotify) CopyFromMessage(obj interface{}) e.Event {
	// 暂时不需要做
	return &EntityInfoChangeNotify{}
}

func (notify *EntityInfoChangeNotify) ToMessage() interface{} {
	pbMsg := &pb.EntityInfoChangeNotify{
		EntityType: pb.ENTITY_TYPE(notify.EntityType),
		EntityId:   notify.EntityId,
		HeroMsg:    notify.HeroMsg.ToMessage().(*pb.HeroMsg),
		ItemMsg:    notify.ItemMsg.ToMessage().(*pb.ItemMsg),
	}

	return pbMsg
}
