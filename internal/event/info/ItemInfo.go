package info

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
)

type ItemInfo struct {
	framework.BaseEvent
	ID           int32
	Type         int32
	Status       int32
	ItemPosition *CoordinateXYInfo
	ItemRadius float32
}

func NewItemInfo(item *model.Prop) *ItemInfo{
	return &ItemInfo{
		ID: item.Id,
		Type: item.PropType,
		Status: item.Status,
		ItemPosition: NewCoordinateInfo(item.Pos.X, item.Pos.Y),
		ItemRadius: 0,
	}
}

func (item *ItemInfo) FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.ItemMsg)
	item.ID = pbMsg.GetItemId()
	item.Type = int32(pbMsg.GetItemType())
	item.Status = int32(pbMsg.GetItemStatus())
	coordi := &CoordinateXYInfo{}
	coordi.FromMessage(pbMsg.GetItemPosition())
	item.ItemPosition = coordi
}

func (item *ItemInfo) CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.ItemMsg)
	coordi := CoordinateXYInfo{}
	coordi.FromMessage(pbMsg.GetItemPosition())
	return &ItemInfo{
		ID:           pbMsg.GetItemId(),
		Type:         int32(pbMsg.GetItemType()),
		Status:       int32(pbMsg.GetItemStatus()),
		ItemPosition: &coordi,
	}
}

func (item *ItemInfo) ToMessage() interface{} {
	pbMsg := &pb.ItemMsg{
		ItemId: item.ID,
		ItemType: pb.ENTITY_TYPE(item.Type),
		ItemStatus: pb.ITEM_STATUS(item.Status),
		ItemPosition: item.ItemPosition.ToMessage().(*pb.CoordinateXY),
	}
	return pbMsg
}
