package info

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
)

type ItemInfo struct {
	framework.BaseEvent
	ID int32
	Type int32
	Status int32
	ItemPosition CoordinateXYInfo
}

func (item *ItemInfo)FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.ItemMsg)
	item.ID = pbMsg.GetItemId()
	item.Type = int32(pbMsg.GetItemType())
	item.Status = int32(pbMsg.GetItemStatus())
		coordi := CoordinateXYInfo{}
		coordi.FromMessage(pbMsg.GetItemPosition())
	item.ItemPosition = coordi
}

func (item *ItemInfo)CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.ItemMsg)
	coordi := CoordinateXYInfo{}
	coordi.FromMessage(pbMsg.GetItemPosition())
	return &ItemInfo{
		ID:   pbMsg.GetItemId(),
		Type: int32(pbMsg.GetItemType()),
		Status: int32(pbMsg.GetItemStatus()),
		ItemPosition: coordi,
	}
}

func (item *ItemInfo)ToMessage() interface{} {
	pbMsg:= &pb.ItemMsg{
		ItemId: item.ID,
		//ItemType: (item.Type),
		//ItemStatus: item.Status,
		ItemPosition: item.ItemPosition.ToMessage().(*pb.CoordinateXY),
	}
	//todo:补齐类型和状态
	return pbMsg
}