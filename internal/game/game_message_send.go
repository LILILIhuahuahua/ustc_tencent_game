package game

import (
	"fmt"
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	event2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/notify"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"github.com/golang/protobuf/proto"
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
	fmt.Printf("%d将消息发送给%d \n", changeHero.ID, notifyHero.ID)
	if err != nil {
		fmt.Println("调用Unicast时发生了错误")
	}
}

func (r *GameRoom) SendHeroPropGlobalInfoNotify() {

}
