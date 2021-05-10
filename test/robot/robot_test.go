package robot

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	event2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/notify"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/response"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/game"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"github.com/golang/protobuf/proto"
	"github.com/xtaci/kcp-go"
	"log"
	"math/rand"
	"testing"
	"time"
)

type Robot struct {
	recvQueue *event.EventRingQueue
	sessionId int32
	hero      *model.Hero
	session   *kcp.UDPSession
	done 	  chan struct{}
}

func NewTestRobot() *Robot {
	return &Robot{
		recvQueue: event.NewEventRingQueue(300),
		sessionId: tools.UUID_UTIL.GenerateInt32UUID(),
		done: make(chan struct{}),
	}
}

func (robot *Robot) boot() {
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
		for {
			//fmt.Println("[go accept]读kcp，生产队列")
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

func (robot *Robot) handle() {
	for {
		select {
		case <-robot.done:
			log.Println("[Handler]收到关闭信号，退出handle处理程序，停止从队列中消费GMessage消息！")
			return
		default:
			e, err := robot.recvQueue.Pop()
			if nil == e { //todo
				continue
			}
			if nil != err {
				log.Println(err)
				continue
			}
			msg := e.(*event2.GMessage)
			if nil == msg.Data {
				log.Printf("[未知包无法处理]收到机器人还不能处理的消息包！GMessage：%+v \n", msg)
				continue
			}
			robot.dispatchGMessage(msg)
		}
	}
}

func (robot *Robot) dispatchGMessage(msg *event2.GMessage) {
	// 二级解码
	data := msg.Data
	switch data.GetCode() {

	case int32(pb.GAME_MSG_CODE_ENTER_GAME_RESPONSE):
		robot.onEnterGame(data.(*response.EnterGameResponse))

	case int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_NOTIFY):
		robot.onEntityInfoChange(data.(*notify.EntityInfoChangeNotify))

	case int32(pb.GAME_MSG_CODE_ENTITY_INFO_NOTIFY):
		robot.onEntityInfoChange(data.(*notify.EntityInfoChangeNotify))

	case int32(pb.GAME_MSG_CODE_GAME_FINISH_NOTIFY):
		robot.close()

	case int32(pb.GAME_MSG_CODE_GAME_GLOBAL_INFO_NOTIFY):

	}
}

func (robot *Robot) onEnterGame(resp *response.EnterGameResponse) {
	robot.hero = model.NewHero("", nil)
	robot.hero.ID = resp.HeroId
	// 开启本地更新英雄位置线程
	go robot.updateHeroPos()
	go robot.updateHeroDirt()
	// 随机化英雄出生位置
	// todo
}

func (robot *Robot) onEntityInfoChange(notify *notify.EntityInfoChangeNotify) {
	if notify.EntityType == int32(pb.ENTITY_TYPE_HERO_TYPE) && notify.EntityId == robot.hero.ID { //只处理自己的状态
		log.Printf("[EntityInfoChange]收到本机器人持有英雄状态改变推送！位置：%+v, 方向：%+v, 信息为：%+v",notify.HeroMsg.HeroPosition, notify.HeroMsg.HeroDirection, notify.HeroMsg)
		heroInfo := notify.HeroMsg
		if heroInfo.Status == int32(pb.HERO_STATUS_DEAD) {
			log.Println("[robot]英雄阵亡！")
			robot.close()
		}
		robot.updateHero(heroInfo)
	}
}

func (robot *Robot) updateHero(heroInfo *info.HeroInfo) {
	if nil != heroInfo.HeroPosition {
		robot.hero.HeroPosition.X = heroInfo.HeroPosition.CoordinateX
		robot.hero.HeroPosition.Y = heroInfo.HeroPosition.CoordinateY
	}
	if nil != heroInfo.HeroDirection {
		robot.hero.HeroDirection.X = heroInfo.HeroDirection.CoordinateX
		robot.hero.HeroDirection.Y = heroInfo.HeroDirection.CoordinateY
	}
	robot.hero.Status = heroInfo.Status
	robot.hero.Score = heroInfo.Score
	robot.hero.Size = heroInfo.Size
	robot.hero.Invincible = heroInfo.Invincible
	robot.hero.Speed = heroInfo.Speed
	robot.hero.SpeedDown = heroInfo.SpeedDown
	robot.hero.SpeedUp = heroInfo.SpeedUp
}

func (robot *Robot) updateHeroPos() {
	for {
		//fmt.Println("[go updateHeroPos]更新位置")
		nowTime := time.Now().UnixNano()
		// 更新玩家位置
		timeElapse := nowTime - robot.hero.UpdateTime
		robot.hero.UpdateTime = nowTime
		distance := float64(timeElapse) * float64(robot.hero.Speed) / 1e9
		x, y := tools.CalXY(distance, robot.hero.HeroDirection.X, robot.hero.HeroDirection.Y)
		robot.hero.HeroPosition.X += x
		robot.hero.HeroPosition.X = tools.GetMax(robot.hero.HeroPosition.X, configs.MapMinX)
		robot.hero.HeroPosition.X = tools.GetMin(robot.hero.HeroPosition.X, configs.MapMaxX)
		robot.hero.HeroPosition.Y += y
		robot.hero.HeroPosition.Y = tools.GetMax(robot.hero.HeroPosition.Y, configs.MapMinY)
		robot.hero.HeroPosition.Y = tools.GetMin(robot.hero.HeroPosition.Y, configs.MapMaxY)
		time.Sleep(50 * 1e6) //睡50ms
	}
}

func (robot *Robot) updateHeroDirt() {
	for  {
		//fmt.Println("[go updateHeroDirt]更新方向")
		rand.Seed(time.Now().UnixNano())
		sleepTime := rand.Int63()%5
		randDict := rand.Intn(4)
		switch randDict {
		case 0:
			robot.hero.HeroDirection.X = 1
			robot.hero.HeroDirection.Y = 0
		case 1:
			robot.hero.HeroDirection.X = -1
			robot.hero.HeroDirection.Y = 0
		case 2:
			robot.hero.HeroDirection.X = 0
			robot.hero.HeroDirection.Y = 1
		case 3:
			robot.hero.HeroDirection.X = 0
			robot.hero.HeroDirection.Y = -1
		}
		heroInfo := info.NewHeroInfo(robot.hero)
		entityInfoChangeReq := request.NewEntityInfoChangeRequest(int32(pb.EVENT_TYPE_HERO_MOVE), robot.hero.ID, -1, "", *heroInfo)
		data := entityInfoChangeReq.ToGMessageBytes()
		robot.session.Write(data)
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}

func (robot *Robot) close() {
	//英雄阵亡，处理机器人资源回收
	log.Println("[robot]关闭并回收机器人资源!")
	close(robot.done)
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
