package test

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	"github.com/golang/protobuf/proto"
	"github.com/xtaci/kcp-go"
	"io"
	"log"
	"testing"
	"time"
)

// TestConnect check whether the client could connect to server correctly
func TestConnect(t *testing.T) {
	// wait for server to become ready
	time.Sleep(time.Second)

	msg := event.GMessage{
		MsgType:     configs.MsgTypeRequest,
		GameMsgCode: configs.EntityInfoChangeRequest,
		SessionId:   0,
		SeqId:       0,
		Data:        &request.EntityInfoChangeRequest{
			EventType:  0,
			HeroId:     0,
			LinkedId:   0,
			LinkedType: "",
			HeroMsg:    info.HeroInfo{},
		},
	}

	pbMsg := msg.ToMessage().(*pb.GMessage)
	out, err := proto.Marshal(pbMsg)
	if err != nil {
		t.Fatalf("get %s", err)
	}
	var buf []byte
	// dial to the echo server
	if sess, err := kcp.DialWithOptions(addr, nil, 0, 0); err == nil {
		for {
			data := time.Now().String()
			log.Println("[Client SENT]: ", data)
			if _, err := sess.Write(out); err != nil {
				log.Println(err.Error())
			}
			if _, err := io.ReadFull(sess, buf); err == nil {
				log.Println("get from server", string(buf))
			} else {
				log.Println(err.Error())
			}
			time.Sleep(time.Second)
		}
	} else {
		log.Fatal(err)
	}
}
