package model

import "time"

type Prop struct {
	Id       int64
	Name     string
	Type     string
	Level    int
	Duration time.Duration // TODO  is usage of duration here right ?
	Skin     string
}
