package tools

import "github.com/LILILIhuahuahua/ustc_tencent_game/model"

func JudgePosition(newPosX, newPosY float32, heroPos model.Coordinate, heroSpeed float32) bool{
	if getAb(heroPos.X - newPosX) > heroSpeed / 1000 * 5 {
		return false
	}

	if getAb(heroPos.Y - newPosY) > heroSpeed / 1000 * 5 {
		return false
	}
	return true
}
