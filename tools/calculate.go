package tools

import (
	"math"
)

// x^2 + y^2 = dis^2
//(0.3 + 1) y^2 = dis^2
//y^2 = dis^2 / (0.3 + 1)
//通过球移动的距离来计算横纵坐标的变化
func CalXY(dis float64, coordX, coordY float32) (x, y float32) {

	xDivY := coordX / coordY //x相当于多少个y
	ySquare := dis * dis / float64(1 + xDivY)
	y = float32(math.Sqrt(ySquare))
	x = y * xDivY
	return x, y
}

func getAb(x float32) float32 {
	if x >= 0 {
		return x
 	} else {
 		return -x
	}
}
