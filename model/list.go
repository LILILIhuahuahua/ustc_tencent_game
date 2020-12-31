package model

// BList represents UDP address of client for broadcasting
type BList struct {
	Id int64
	PlayerId int64
	ClientAddr string
}

// GList represents the dgs address of a game
type GList struct {
	Id int64
	GameId int64
	DgsAddr string
}

// MList represents the network information.
// It maps client address with dgs address.
type MList struct {
	Id int64
	BId int64
	GId int64
}
