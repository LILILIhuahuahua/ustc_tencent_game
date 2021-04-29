package test

import (
	"fmt"
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/golang/protobuf/proto"
	"github.com/xtaci/kcp-go"
	"log"
	"testing"
	"time"
)

var (
	seqID     int32
	localAddr = "127.0.0.1:8888"
)

func PingServer(reqCount int32, sessionID int32, portNumber int32) (map[int32]int64, map[int32]int64, int) {
	buf := make([]byte, 4096)
	sent, recv := make(map[int32]int64), make(map[int32]int64)
	sess, err := kcp.DialWithOptions(configs.RemoteCLB, nil, 0, 0)
	startTime := time.Now().Second()

	// send enterGameReq to server continuously
	if err == nil {
		for seqID < reqCount {
			enterGameReq := getEnterGameReq2(sessionID, &pb.ConnectMsg{
				Ip:   "127.0.0.1",
				Port: portNumber,
			})

			now := time.Now().UnixNano()
			enterGameReq.SeqId = seqID
			enterGameReq.SendTime = now
			sent[seqID] = now

			data, err := proto.Marshal(enterGameReq)
			if err != nil {
				log.Println("fail to marshal enterGameReq")
				continue
			}

			sess.Write(data)
			ch := make(chan bool)
			go func() {
				var msg pb.GMessage
				flag := true
				for flag {
					n, err := sess.Read(buf)
					if err != nil {
						log.Println("fail to read data from server")
					}

					proto.Unmarshal(buf[:n], &msg)
					if msg.MsgType == pb.MSG_TYPE_RESPONSE && msg.MsgCode == pb.GAME_MSG_CODE_ENTER_GAME_RESPONSE && msg.SeqId == seqID {
						recv[seqID] = time.Now().UnixNano()
						flag = false
					}
				}

				ch <- true
			}()

			<-ch
			seqID++
			sessionID++
			portNumber++
		}

	} else {
		log.Fatal(err)
	}

	totalTime := time.Now().Second() - startTime
	return sent, recv, totalTime
}

// TestLatency is used to test rtt time from client to remote loadBalancer
func TestLatency(t *testing.T) {
	var reqCount int32
	var sessionID int32
	var portNumber int32 = 10000

	for i := 0; i < 5; i++ {
		sent, recv, totalTime := PingServer(reqCount, sessionID, portNumber)

		var count int64
		for seq, t := range recv {
			elapse := (t - sent[seq]) / 1000000
			fmt.Printf("sedID %d, send time %d , recv time %d, rtt %v ms\n", seq, sent[seq], recv[seq], elapse)
			count += elapse
		}

		n := int64(len(recv))
		if n != 0 {
			avg := count / n

			fmt.Println("####################################################################")
			fmt.Printf("#### Round %d, %d req, average rtt %v, totalTime %d s      #####\n", i, n, avg, totalTime)
			fmt.Println("####################################################################")
		}
		reqCount += 50
		sessionID += 50
		portNumber += 50
	}

}

func getEnterGameReq2(sessionID int32, client *pb.ConnectMsg) *pb.GMessage {
	playerID := sessionID
	enterGameReq := &pb.GMessage{
		MsgType: pb.MSG_TYPE_REQUEST,
		MsgCode: pb.GAME_MSG_CODE_ENTER_GAME_REQUEST,
		Request: &pb.Request{
			EnterGameRequest: &pb.EnterGameRequest{
				PlayerId:         playerID,
				ClientConnectMsg: client,
			},
		},
	}

	return enterGameReq
}
