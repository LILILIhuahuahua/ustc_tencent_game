package game

import (
	"fmt"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
)

func (g *GameRoom) DeleteUnavailableSession() error {
	var needDelete []*framework.BaseSession
	// 将不能正常通信的session存储到needDelete中
	g.sessions.Range(func(_, obj interface{}) bool {
		sess := obj.(*framework.BaseSession)
		if !sess.IsAvailable() {
			needDelete = append(needDelete, sess)
		}
		return true
	})
	//fmt.Println(needDelete)
	for _, session := range needDelete {
		session.ChangeStatus(configs.SessionStatusDead)
		err := session.CloseKcpSession()
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *GameRoom) DeleteOfflinePlayer() error {
	towers := g.GetTowers()
	var needDelete []*framework.BaseSession
	g.sessions.Range(func(_, obj interface{}) bool {
		sess := obj.(*framework.BaseSession)
		if sess.IsDeprecated() {
			needDelete = append(needDelete, sess)
		}
		return true
	})

	for _, session := range needDelete {
		deletedObj, ok := g.SessionHeroMap.Load(session.Id)
		if !ok {
			continue
			//return errors.New("玩家不存在")
		}
		hero := deletedObj.(*model.Hero)
		hero.ChangeHeroStatus(configs.Dead)
		towers[hero.TowerId].HeroLeave(hero)
		session.ChangOfflineStatus(true)
		fmt.Println("我调用了玩家删除函数")
	}
	return nil
}
