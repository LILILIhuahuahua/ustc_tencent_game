package test

import (
	"github.com/xtaci/kcp-go"
	"log"
	"testing"
	"time"
)

// TestConnect check whether the client could connect to server correctly
func TestConnect(t *testing.T) {
	// wait for server to become ready
	time.Sleep(time.Second)

	// dial to the echo server
	if sess, err := kcp.DialWithOptions(addr, nil, 10, 3); err == nil {
		for {
			data := time.Now().String()
			//buf := make([]byte, len(data))
			log.Println("[Client SENT]: ", data)
			if _, err := sess.Write([]byte(data)); err != nil {
				log.Println(err.Error())
			}
			time.Sleep(time.Second)
		}
	} else {
		log.Fatal(err)
	}
}
