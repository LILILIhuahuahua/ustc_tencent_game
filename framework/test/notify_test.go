package test

import (
	"fmt"
	"github.com/xtaci/kcp-go"
	"io"
	"log"
	"testing"
	"time"
)

// TestNotify check if the client could be notified by server
func TestNotify(t *testing.T) {
	// wait for server to become ready
	time.Sleep(time.Second)

	// dial to the echo server
	if sess, err := kcp.DialWithOptions(addr, nil, 10, 3); err == nil {
		for {
			data := time.Now().String()
			buf := make([]byte, len(data))

			_, err := io.ReadFull(sess, buf)
			if err != nil {
				log.Println("Failed to read data")
				sess.Close()
				return
			}

			fmt.Println(string(buf))
		}
	}
}
