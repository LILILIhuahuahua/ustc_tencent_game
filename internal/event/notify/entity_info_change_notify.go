package notify

import (
	//pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	e "github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"github.com/golang/protobuf/proto"
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
	HeroMsg    *info.HeroInfo
	ItemMsg    *info.ItemInfo
}

func NewEntityInfoChangeNotify(entityType int32, entityId int32, heroInfo *info.HeroInfo, itemInfo *info.ItemInfo) *EntityInfoChangeNotify {
	return &EntityInfoChangeNotify{
		EntityType: entityType,
		EntityId:   entityId,
		HeroMsg:    heroInfo,
		ItemMsg:    itemInfo,
	}
}

func (notify *EntityInfoChangeNotify) FromMessage(obj interface{}) {
	//不需要做,因为这个消息只会由服务端发送给客户端，不涉及到解析
}

func (notify *EntityInfoChangeNotify) CopyFromMessage(obj interface{}) e.Event {
	//不需要做,因为这个消息只会由服务端发送给客户端，不涉及到解析
	return &EntityInfoChangeNotify{}
}

func (notify *EntityInfoChangeNotify) ToMessage() interface{} {
	pbMsg := &pb.EntityInfoChangeNotify{
		EntityType: pb.ENTITY_TYPE(notify.EntityType),
		EntityId:   notify.EntityId,
		//HeroMsg:    notify.HeroMsg.ToMessage().(*pb.HeroMsg),
		//ItemMsg:    notify.ItemMsg.ToMessage().(*pb.ItemMsg),
	}
	if nil != notify.HeroMsg {
		pbMsg.HeroMsg = notify.HeroMsg.ToMessage().(*pb.HeroMsg)
	}
	if nil != notify.ItemMsg {
		pbMsg.ItemMsg = notify.ItemMsg.ToMessage().(*pb.ItemMsg)
	}
	return pbMsg
}

func (notify *EntityInfoChangeNotify) ToGMessageBytes() []byte {
	n := &pb.Notify{
		EntityInfoChangeNotify: notify.ToMessage().(*pb.EntityInfoChangeNotify),
	}
	msg := pb.GMessage{
		MsgType:  pb.MSG_TYPE_NOTIFY,
		MsgCode:  pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_NOTIFY,
		Notify:   n,
		SendTime: tools.TIME_UTIL.NowMillis(),
	}
	out, _ := proto.Marshal(&msg)
	return out
}
