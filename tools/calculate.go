package tools

import (
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"math"
)

// x^2 + y^2 = dis^2
//(0.3 + 1) y^2 = dis^2
//y^2 = dis^2 / (0.3 + 1)
//通过球移动的距离来计算横纵坐标的变化
func CalXY(dis float64, dir model.Coordinate) (x, y float32) {
	xDivY := dir.X / dir.Y //x相当于多少个y
	ySquare := dis * dis / float64(1 + xDivY)
	y = float32(math.Sqrt(ySquare))
	x = y * xDivY
	return x, y
}
