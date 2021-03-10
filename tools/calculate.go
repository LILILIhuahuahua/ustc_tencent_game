package tools

import (
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"math"
)

// x^2 + y^2 = dis^2
//(0.3 + 1) y^2 = dis^2
//y^2 = dis^2 / (0.3 + 1)
//通过球移动的距离来计算横纵坐标的变化
func CalXY(dis float64, coordX, coordY float32) (x, y float32) {
	if coordX == float32(0) {
		if coordY >= 0 {
			return float32(0), float32(dis)
		}
		return float32(0), float32(-dis)
	}
	if coordY == float32(0) {
		if coordX >= 0 {
			return float32(dis), float32(0)
		}
		return float32(-dis), float32(0)
	}
	xDivY := coordX / coordY //x相当于多少个y
	ySquare := dis * dis / float64(1+ xDivY * xDivY)
	y = float32(math.Sqrt(ySquare))
	x = y * xDivY
	if coordX < 0 {
		x = -getAb(x)
	}
	if coordY < 0 {
		y = -getAb(y)
	}
	return x, y
}

func getAb(x float32) float32 {
	if x >= 0 {
		return x
	} else {
		return -x
	}
}

// 根据游戏坐标计算towerId
func CalTowerId(coordX, coordY float32) int32 { //通过游戏中的横纵坐标来计算出Tower的ID
	col := int32((coordX - configs.MapMinX) / configs.TowerDiameter)
	row := int32((coordY - configs.MapMinY) / configs.TowerDiameter)
	return row*configs.TowerCols + col
}

// 根据towerId计算周围towerId
func GetOtherTowers(towerId int32) []int32 {
	var towers []int32
	towerRow := towerId / configs.TowerCols
	towerCol := towerId % configs.TowerCols
	if towerRow > 0 &&
		towerRow < configs.TowerRows - 1 && //减1的原因是因为TowerRows是从1开始的
		towerCol > 0 &&
		towerCol < configs.TowerCols - 1 {
		towers = append(towers,
			towerId - 1,
			towerId + 1,
			towerId + configs.TowerCols,
			towerId - configs.TowerCols,
			towerId + configs.TowerCols + 1,
			towerId + configs.TowerCols - 1,
			towerId - configs.TowerCols + 1,
			towerId - configs.TowerCols - 1,
			)
		return towers
	}

	if towerRow == 0 {
		switch towerCol {
		case 0:
			towers = append(
				towers,
				towerId + 1,
				towerId + configs.TowerCols,
				towerId + configs.TowerCols + 1,
				)
			return towers
		case configs.TowerCols - 1:
			towers = append(
				towers,
				towerId - 1,
				towerId + configs.TowerCols,
				towerId + configs.TowerCols - 1,
				)
			return towers
		default:
			towers = append(
				towers,
				towerId + 1,
				towerId - 1,
				towerId + configs.TowerCols,
				towerId + configs.TowerCols + 1,
				towerId + configs.TowerCols - 1,
				)
			return towers
		}
	} else if towerRow == configs.TowerRows - 1 {
		switch towerCol {
		case 0:
			towers = append(
				towers,
				towerId + 1,
				towerId - configs.TowerCols,
				towerId - configs.TowerCols + 1,
				)
			return towers
		case configs.TowerCols - 1:
			towers = append(
				towers,
				towerId - 1,
				towerId - configs.TowerCols,
				towerId - configs.TowerCols - 1,
				)
			return towers
		default:
			towers = append(
				towers,
				towerId + 1,
				towerId - 1,
				towerId - configs.TowerCols,
				towerId - configs.TowerCols + 1,
				towerId - configs.TowerCols - 1,
				)
			return towers
		}
	} else {
		switch towerCol {
		case 0:
			towers = append(
				towers,
				towerId + 1,
				towerId + configs.TowerCols,
				towerId - configs.TowerCols,
				towerId + configs.TowerCols + 1,
				towerId - configs.TowerCols + 1,
				)
			return towers
		case configs.TowerCols - 1:
			towers = append(
				towers,
				towerId + configs.TowerCols,
				towerId - configs.TowerCols,
				towerId + configs.TowerCols - 1,
				towerId - configs.TowerCols - 1,
				towerId - 1,
				)
			return towers
		}
	}
	return nil
}

func GetMax(x, y float32) float32 {
	if x > y {
		return x
	} else {
		return y
	}
}

func GetMin(x, y float32) float32 {
	if x < y {
		return x
	} else {
		return y
	}
}