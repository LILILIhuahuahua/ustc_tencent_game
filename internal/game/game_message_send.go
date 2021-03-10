package game

import (
	"fmt"
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	event2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/notify"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"github.com/golang/protobuf/proto"
	"log"
)

func (r *GameRoom) SendHeroViewNotify(changeHero *model.Hero, notifyHero *model.Hero, notifyType int32) {
	heroViewNotify := &notify.HeroViewNotify{
		HeroId:    changeHero.ID,
		ViewType:  notifyType,
		HeroMsg:   changeHero.ToEvent(),
	}
	notifyMsg := event2.GMessage{
		MsgType:     configs.MsgTypeNotify,
		GameMsgCode: configs.HeroViewNotify,
		Data:        heroViewNotify,
	}
	pbMsg := notifyMsg.ToMessage().(*pb.GMessage)
	notifySession := notifyHero.Session
	out, err := proto.Marshal(pbMsg)
	if err != nil {
		fmt.Println("调用SendHeroViewNotify时发生了错误")
	}
	err = r.Unicast(out, notifySession)
	if err != nil {
		fmt.Println("调用Unicast时发生了错误")
	}
}

func (r *GameRoom) SendHeroPropGlobalInfoNotify(towers []int32, session *framework.BaseSession) {
	ts := r.GetTowers()
	var heroMsg []*model.Hero
	var heroEvent []info.HeroInfo
	//后面加上道具
	for _, id := range towers {
		hs := ts[id].GetHeros()
		heroMsg = append(heroMsg, hs...)
	}
	for _, h := range heroMsg {
		heroEvent = append(heroEvent, h.ToEvent())
		fmt.Printf("向%d，发送%d的位置信息 -------刚刚进入灯塔\n", session.Id, h.Session.Id)
	}
	notify := notify.GameGlobalInfoNotify{
		HeroNumber: int32(len(heroEvent)),
		HeroInfos:  heroEvent,
		//ItemInfos:  nil,
		//MapInfo:    info.MapInfo{},
	}
	msg := event2.GMessage{
		MsgType:     configs.MsgTypeNotify,
		GameMsgCode: configs.GameGlobalInfoNotify,
		Data:        &notify,
	}

	pbMsg := msg.ToMessage().(*pb.GMessage)
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		log.Printf("获取tower中hero信息的时候解析出错")
	}
	r.Unicast(data, session)
}
