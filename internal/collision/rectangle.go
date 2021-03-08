package collision

type Rectangle struct {
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
		X: (maxX-minX)/2,
		Y: (maxY-minY)/2,
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

func (rec *Rectangle)GetX() float32 {
	return rec.X
}

func (rec *Rectangle)GetY() float32 {
	return rec.Y
}
