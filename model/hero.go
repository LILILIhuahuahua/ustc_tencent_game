package model

import (
	"fmt"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"sync"
	"time"
)

type Hero struct {
	ID            int32
	Status        int32
	Size          float32
	Speed         float32
	HeroDirection Coordinate
	HeroPosition  Coordinate
	CreateTime    int64
	UpdateTime    int64  // ns
	TowerId		  int32 // 所属的towerId
	OtherTowers	  sync.Map // 九宫格内其他的tower TowerId *aoi.Tower
	Session       *framework.BaseSession //该hero对应的session
}

func NewHero(sess *framework.BaseSession) *Hero {
	h := &Hero{}
	//初始化英雄数据
	h.Init(sess)
	fmt.Printf("新生成的heroId为%d \n", h.ID)
	return h
}

func (h *Hero) Init(sess *framework.BaseSession) {
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
	h.Status = configs.Live
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
