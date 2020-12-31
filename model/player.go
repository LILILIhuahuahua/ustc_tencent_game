package model

import "time"

type Player struct {
	Id int64
	Nickname string
	Level int
	Score int
	TimeCount time.Time
	StartAt time.Time
	Status int
	AccountId int64
	RecentAddr string
	HighestRank int
}