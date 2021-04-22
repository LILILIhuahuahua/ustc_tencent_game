package model

import (
	"fmt"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
	"testing"
)

func TestNewProps(t *testing.T) {
	ids := make(map[int32]bool)
	coods := make(map[info.CoordinateXYInfo]bool)
	m := New()
	for k, v := range m.props {
		fmt.Printf("index: %d, id: %d, status: %v, coords: %v\n", k, v.id, v.status, v.pos)
		if _, ok := ids[v.id]; ok {
			t.Errorf("id duplicate: %v", v.id)
		} else {
			ids[v.id] = true
		}

		if !(v.pos.CoordinateX > configs.MapMinX && v.pos.CoordinateX < configs.MapMaxX) ||
			!(v.pos.CoordinateY > configs.MapMinY && v.pos.CoordinateY < configs.MapMaxY) {
			t.Errorf("index: %d, id: %d, coordinates out of range", k, v.id)
		}

		if _, ok := coods[v.pos]; ok {
			t.Errorf("duplicate coordinates: %v", v.pos)
		} else {
			coods[v.pos] = true
		}

	}

	t.Logf("total props: %d", len(m.props))
}
