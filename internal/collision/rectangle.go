package collision

import "math"

//碰撞检测功能基本单元
type Rectangle struct {
	ID int32
	Type int32
	Size float32
	MinX float32	//最小边界点横坐标（左下角）
	MinY float32	//最小边界点纵坐标（左下角）
	MaxX float32	//最大边界点横坐标（右上角）
	MaxY float32	//最大边界点纵坐标（右上角）
	X float32	//中心点横坐标
	Y float32	//中心点纵坐标
}

func NewRectangleByBounds(minX float32,minY float32,maxX float32,maxY float32) *Rectangle {
	return &Rectangle{
		MinX: minX,
		MinY: minY,
		MaxX: maxX,
		MaxY: maxY,
		X: minX + (maxX-minX) / 2,
		Y: minY + (maxY-minY) / 2,
	}
}

func NewRectangleByXY(x float32,y float32) *Rectangle {
	return &Rectangle{
		MinX: x,
		MinY: y,
		MaxX: x,
		MaxY: y,
		X: x,
		Y: y,
	}
}

func NewRectangleByObj(id int32, t int32, size float32, x float32,y float32) *Rectangle{
	return &Rectangle{
		ID: id,
		Type: t,
		Size: size,
		X: x,
		Y: y,
	}
}

func (rec *Rectangle)GetX() float32 {
	return rec.X
}

func (rec *Rectangle)GetY() float32 {
	return rec.Y
}

func CheckCollision(obj1 *Rectangle, obj2 *Rectangle) bool{
	desX := float64(obj1.GetX() - obj2.GetX())
	desY := float64(obj1.GetY() - obj2.GetY())
	des := math.Sqrt(math.Pow(desX,2) + math.Pow(desY,2))
	max := float64(obj1.Size)
	if obj2.Size > obj1.Size {max = float64(obj2.Size)}
	if des < max {
		return true
	}
	return false
}
