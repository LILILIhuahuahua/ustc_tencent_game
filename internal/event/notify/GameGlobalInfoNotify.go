package notify

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
)

type GameGlobalInfoNotify struct {
	framework.BaseEvent
	HeroNumber int32
	Time       int64
	HeroInfos  []info.HeroInfo
	ItemInfos  []info.ItemInfo
	MapInfo    info.MapInfo
}

func (notify *GameGlobalInfoNotify) FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.GameGlobalInfoNotify)
	notify.SetCode(int32(pb.GAME_MSG_CODE_GAME_GLOBAL_INFO_NOTIFY))
	notify.HeroNumber = pbMsg.GetHeroNumber()
	notify.Time = pbMsg.GetTime()
	if nil == notify.HeroInfos {
		notify.HeroInfos = make([]info.HeroInfo, 0)
	}
	for _, heroMsg := range pbMsg.GetHeroMsg() {
		heroInfo := info.HeroInfo{}
		heroInfo.FromMessage(heroMsg)
		notify.HeroInfos = append(notify.HeroInfos, heroInfo)
	}

	if nil == notify.ItemInfos {
		notify.ItemInfos = make([]info.ItemInfo, 0)
	}
	for _, itemMsg := range pbMsg.GetItemMsg() {
		itemInfo := info.ItemInfo{}
		itemInfo.FromMessage(itemMsg)
		notify.ItemInfos = append(notify.ItemInfos, itemInfo)
	}
	mapInfo := info.MapInfo{}
	mapInfo.FromMessage(pbMsg.GetMapMsg())
	notify.MapInfo = mapInfo
}

func (notify *GameGlobalInfoNotify) CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.Notify).GameGlobalInfoNotify
	n := &GameGlobalInfoNotify{}
	n.FromMessage(pbMsg)
	return n
}

func (notify *GameGlobalInfoNotify) ToMessage() interface{} {
	pbMsg := &pb.GameGlobalInfoNotify{
		HeroNumber: notify.HeroNumber,
		Time:       notify.Time,
	}
	for _, heroInfo := range notify.HeroInfos {
		heroPbMsg := heroInfo.ToMessage().(*pb.HeroMsg)
		pbMsg.HeroMsg = append(pbMsg.HeroMsg, heroPbMsg)
	}
	for _, itemInfo := range notify.ItemInfos {
		itemPbMsg := itemInfo.ToMessage().(*pb.ItemMsg)
		pbMsg.ItemMsg = append(pbMsg.ItemMsg, itemPbMsg)
	}
	pbMsg.MapMsg = notify.MapInfo.ToMessage().(*pb.MapMsg)
	return pbMsg
}
