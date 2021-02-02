package tools

import (
	"encoding/binary"
	"github.com/satori/go.uuid"
	guuid "github.com/google/uuid"
)

var UUID_UTIL = &UUIDUtil{}

type UUIDUtil struct {
}

func (util *UUIDUtil)GenerateInt64UUID() int64 {
 	u:=uuid.NewV4()
 	rlt := binary.BigEndian.Uint64(u[0:8])
 	return int64(rlt)
}

func (util *UUIDUtil)GenerateInt32UUID() int32 {

	return int32(guuid.New().ID())
}
