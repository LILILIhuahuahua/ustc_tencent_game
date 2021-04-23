package prop

import (
	"errors"
	proto "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	guuid "github.com/google/uuid"
	"math/rand"
	"sync"
	"time"
)

type PropsManger struct {
	mu    *sync.RWMutex
	props map[int32]*model.Prop
}

var (
	ErrNilProp        = errors.New("nil Prop value")
	ErrPropNotExist   = errors.New("Prop not exist in propManger")
	ErrPropDuplicate  = errors.New("Prop duplicate in propManger")
	ErrNilPropManager = errors.New("nil propManger")
)

// New return an instance of propsManger, which contains many props
func New() *PropsManger {
	return &PropsManger{
		mu:    &sync.RWMutex{},
		props: newProps(configs.MaxPropCountInMap),
	}
}

// GetProps return all props in propManager
func (p *PropsManger) GetProps() ([]*model.Prop, error) {
	if p == nil {
		return nil, ErrNilPropManager
	}

	var props []*model.Prop
	for _, v := range p.props {
		props = append(props, v)
	}
	return props, nil
}

func (p *PropsManger) GetProp(id int32) (*model.Prop, error) {
	if p == nil {
		return nil, ErrNilPropManager
	}
	var prop *model.Prop
	prop = p.props[id]
	return prop, nil
}

// AddProp add Prop to propManger,if Prop is nil or Prop has existed on propManager, it will return error.
func (p *PropsManger) AddProp(pr *model.Prop) error {
	if pr == nil {
		return ErrNilProp
	}
	if _, ok := p.props[pr.Id]; ok {
		return ErrPropDuplicate
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.props[pr.Id] = pr

	return nil
}

// RemoveProp remove Prop according to Prop id. If Prop does not exist in propManger, it will return error
func (p *PropsManger) RemoveProp(id int32) error {
	if _, ok := p.props[id]; !ok {
		return ErrPropNotExist
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.props, id)

	return nil
}

// newProps generate a bunch of props randomly
func newProps(count int) map[int32]*model.Prop {
	minX := configs.MapMinX
	minY := configs.MapMinY
	maxX := configs.MapMaxX
	maxY := configs.MapMaxY
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	m := make(map[int32]*model.Prop, count)
	for i := 0; i < count; i++ {
		x := minX + r.Float32()*(maxX-minX)
		y := minY + r.Float32()*(maxY-minY)
		pid := int32(guuid.New().ID())
		z := r.Intn(10)
		var propType int32
		if z <= 2 {
			propType = configs.PropTypeInvincible // 无敌
		} else if z <= 4 {
			propType = configs.PropTypeJump // 跃迁道具
		} else {
			propType = configs.PropTypeFood
		}

		m[pid] = &model.Prop{
			Id:     pid,
			Status: int32(proto.ITEM_STATUS_ITEM_LIVE),
			Pos: model.Coordinate{
				X: x,
				Y: y,
			},
			PropType: propType,
		}
	}
	return m
}