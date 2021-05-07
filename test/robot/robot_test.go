package robot

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	event2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/response"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/game"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"github.com/golang/protobuf/proto"
	"github.com/xtaci/kcp-go"
	"log"
	"testing"
	"time"
)

type Robot struct {
	recvQueue *event.EventRingQueue
	sessionId int32
	hero      *model.Hero
	session   *kcp.UDPSession
}

func NewTestRobot() *Robot {
	return &Robot{
		recvQueue: event.NewEventRingQueue(300),
		sessionId: tools.UUID_UTIL.GenerateInt32UUID(),
	}
}

func (robot *Robot) boot()  {
	go robot.handle()
	robot.accept()
}

func (robot *Robot) accept() {
	if sess, err := kcp.DialWithOptions("127.0.0.1:8888", nil, 0, 0); err == nil {
		//sess调优
		sess.SetNoDelay(1, 10, 2, 1)
		sess.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(2)))
		sess.SetACKNoDelay(true)
		robot.session = sess
		//开启进入世界流程
		data := request.NewEnterGameRequest(robot.sessionId, *info.NewConnectInfo("0.0.0.0", -1), "").ToGMessageBytes()
		sess.Write(data)
		buf := make([]byte, 4096)
		for  {
			num, _ := sess.Read(buf)
			if num > 0 {
				pbGMsg := &pb.GMessage{}
				proto.Unmarshal(buf, pbGMsg)
				msg := event2.GMessage{}
				m := msg.CopyFromMessage(pbGMsg)
				robot.recvQueue.Push(m)
				//buf清零
				for i := range buf {
					buf[i] = 0
				}
			}
		}
	}
}

func (robot *Robot) handle()  {
	for  {
		e, err := robot.recvQueue.Pop()
		if nil == e { //todo
			continue
		}
		if nil != err {
			log.Println(err)
			continue
		}
		msg := e.(*event2.GMessage)
		robot.dispatchGMessage(msg)
	}
}

func (robot *Robot) dispatchGMessage(msg *event2.GMessage)  {
	// 二级解码
	data := msg.Data

	switch data.GetCode() {
	case int32(pb.GAME_MSG_CODE_ENTER_GAME_RESPONSE):
		robot.onEnterGame(data.(*response.EnterGameResponse))

	}
}

func (robot *Robot) onEnterGame(resp *response.EnterGameResponse) {
	robot.hero = model.NewHero("", nil)
	robot.hero.ID = resp.HeroId
}

func TestRobot(t *testing.T) {
	// 初始化framework包组件
	g := &game.GameStarter{
	}
	g.Init()
	// 启动机器人
	robot := NewTestRobot()
	robot.boot()
}
