package model

type coordinate struct {
	X float64
	Y float64
}

type Snake struct {
	Id int64
	Len int
	Direction int
	Kills int
	Statue int
	Speed int
	Coordinates []coordinate
}
