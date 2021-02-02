package model

import "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"

type Coordinate struct {
	X float32
	Y float32
}

func (c *Coordinate)ToEvent() info.CoordinateXYInfo{
	return info.CoordinateXYInfo{
		CoordinateX: c.X,
		CoordinateY: c.Y,
	}
}

func (c *Coordinate)FromEvent(info info.CoordinateXYInfo){
	c.X = info.CoordinateX
	c.Y = info.CoordinateY
}