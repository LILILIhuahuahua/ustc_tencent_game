package model

import (
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"log"
	"sync"
	"time"
)

type Hero struct {
	ID                  int32
	Name                string
	Status              int32
	Size                float32
	Speed               float32
	Score               int32
	Rank                int32
	Invincible          bool
	SpeedUp             bool
	SpeedDown           bool
	InvincibleStartTime int64
	SpeedUpStartTime    int64
	SpeedDownStartTime  int64
	HeroDirection       Coordinate
	HeroPosition        Coordinate
	CreateTime          int64
	UpdateTime          int64                  // ns
	TowerId             int32                  // 所属的towerId
	OtherTowers         sync.Map               // 九宫格内其他的tower TowerId *aoi.Tower
	Session             *framework.BaseSession //该hero对应的session
}

func NewHero(name string, sess *framework.BaseSession) *Hero {
	h := &Hero{}
	//初始化英雄数据
	h.Init(name, sess)
	log.Printf("[Hero]初始化新英雄！hero：%v \n", h)
	return h
}

func (h *Hero) Init(name string, sess *framework.BaseSession) {
	dcit := Coordinate{
		X: configs.HeroInitDirectionX,
		Y: configs.HeroInitDirectionY,
	}
	pos := Coordinate{
		X: configs.HeroInitPositionX,
		Y: configs.HeroInitPositionY,
	}
	nowTime := time.Now().UnixNano()
	h.ID = tools.UUID_UTIL.GenerateInt32UUID()
	h.Name = name
	h.Rank = 0
	h.Score = 0
	h.Status = configs.HeroStatusLive
	h.Size = configs.HeroInitSize
	h.Speed = configs.HeroMoveSpeed
	h.HeroDirection = dcit
	h.HeroPosition = pos
	h.CreateTime = nowTime
	h.UpdateTime = nowTime
	h.Session = sess
}

//func (h *Hero) ToEvent() info.HeroInfo {
//	return info.HeroInfo{
//		ID:            h.ID,
//		Status:        h.Status,
//		Speed:         h.Speed,
//		Size:          h.Size,
//		HeroDirection: h.HeroDirection.ToEvent(),
//		HeroPosition:  h.HeroPosition.ToEvent(),
//	}
//}
//
//func (h *Hero) FromEvent(info info.HeroInfo) {
//	h.ID = info.ID
//	h.Status = info.Status
//	h.Speed = info.Speed
//	h.Size = info.Size
//	h.HeroDirection = Coordinate{}
//	h.HeroDirection.FromEvent(info.HeroDirection)
//	h.HeroPosition = Coordinate{}
//	h.HeroPosition.FromEvent(info.HeroPosition)
//}

func (h *Hero) ChangeHeroStatus(status int32) {
	h.Status = status
}
