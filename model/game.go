package model

import "time"

type Game struct {
	Id        int64
	Players   []Player
	Count     int
	CreateAt  time.Time
	CountDown time.Time
	Props     []Prop
	Ranks     map[int64]int // TODO
	Snakes    []Hero
}
