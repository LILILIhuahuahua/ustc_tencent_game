package prop

import (
	"fmt"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"testing"
)

func TestNewProps(t *testing.T) {
	ids := make(map[int32]bool)
	//coods := make(map[info.CoordinateXYInfo]bool)
	m := New()
	fmt.Println(m)
	for k, v := range m.props {
		fmt.Printf("index: %d, id: %d, status: %v, coords: %v\n", k, v.Id, v.Status, v.Pos)
		if _, ok := ids[v.Id]; ok {
			t.Errorf("id duplicate: %v", v.Id)
		} else {
			ids[v.Id] = true
		}

		if !(v.Pos.X > configs.MapMinX && v.Pos.X < configs.MapMaxX) ||
			!(v.Pos.Y > configs.MapMinY && v.Pos.Y < configs.MapMaxY) {
			t.Errorf("index: %d, id: %d, coordinates out of range, x: %f, y: %f", k, v.Id, v.Pos.X, v.Pos.Y)
		}

		//if _, ok := coods[v.pos]; ok {
		//	t.Errorf("duplicate coordinates: %v", v.pos)
		//} else {
		//	coods[v.pos] = true
		//}

	}

	t.Logf("total props: %d", len(m.props))
}
