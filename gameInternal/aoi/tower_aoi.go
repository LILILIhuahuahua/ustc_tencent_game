package aoi

import "github.com/LILILIhuahuahua/ustc_tencent_game/configs"

func InitTowers() []*Tower {
	var towers []*Tower
	for i := int32(0); i < configs.TowerRows*configs.TowerCols; i++ {
		towers = append(towers, InitTower(i))
	}
	return towers
}
