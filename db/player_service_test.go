package db

import (
	"log"
	"testing"
)

func TestPlayerUpdateHighestScoreByPlayerId(t *testing.T) {
	InitConnection("localhost:8890")
	err := PlayerUpdateHighestScoreByPlayerId(560, 2333)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("成功了")
}
