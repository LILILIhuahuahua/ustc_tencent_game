package response

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
)
type EntityInfoChangeResponse struct {
	changeResult bool
}

func (this *EntityInfoChangeResponse)ToMessage() interface{} {
	return pb.EntityInfoChangeResponse{
		ChangeResult: this.changeResult,
	}
}

