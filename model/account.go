package model

import "time"

// TODO make tencent doc consistent with code
type Account struct {
	Id int64
	LoginName string
	LoginPassword string
	AccountAvatar string
	Level int
	Skin string
	Deleted bool
	Region string
	CreateAt time.Time
	UpdateAt time.Time
}
