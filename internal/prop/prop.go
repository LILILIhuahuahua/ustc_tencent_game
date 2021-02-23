package prop

import (
	"errors"
	"github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
	guuid "github.com/google/uuid"
	"math/rand"
	"sync"
	"time"
)

var (
	ErrNilProp        = errors.New("nil prop value")
	ErrPropNotExist   = errors.New("prop not exist in propManger")
	ErrPropDuplicate  = errors.New("prop duplicate in propManger")
	ErrNilPropManager = errors.New("nil propManger")
)

type prop struct {
	id     int32
	status int32
	pos    info.CoordinateXYInfo
}

type PropsManger struct {
	mu    *sync.RWMutex
	props map[int32]*prop
}

// ID returns the id of prop
func (p *prop) ID() int32 {
	return p.id
}

// Status returns the status (alive or dead) of prop
func (p *prop) Status() int32 {
	return p.status
}

// GetX returns the x coordinate of prop
func (p *prop) GetX() float32 {
	return p.pos.CoordinateX
}

// GetY returns the y coordinate of prop
func (p *prop) GetY() float32 {
	return p.pos.CoordinateY
}

// New return an instance of propsManger, which contains many props
func New() *PropsManger {
	return &PropsManger{
		mu:    &sync.RWMutex{},
		props: newProps(configs.MapMinX, configs.MapMaxX, configs.MapMinY, configs.MapMaxY, configs.MaxPropCountInMap),
	}
}

// GetProps return all props in propManager
func (p *PropsManger) GetProps() ([]prop, error) {
	if p == nil {
		return nil, ErrNilPropManager
	}

	var props []prop
	for _, v := range p.props {
		props = append(props, *v)
	}
	return props, nil
}

// AddProp add prop to propManger,if prop is nil or prop has existed on propManager, it will return error.
func (p *PropsManger) AddProp(pr *prop) error {
	if pr == nil {
		return ErrNilProp
	}
	if _, ok := p.props[pr.id]; ok {
		return ErrPropDuplicate
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.props[pr.id] = pr

	return nil
}

// RemoveProp remove prop according to prop id. If prop does not exist in propManger, it will return error
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
func newProps(minX float32, maxX float32, minY float32, maxY float32, count int) map[int32]*prop {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	m := make(map[int32]*prop, count)
	for i := 0; i < count; i++ {
		x := minX + r.Float32()*(maxX-minX)
		y := minY + r.Float32()*(maxY-minY)
		pid := int32(guuid.New().ID())

		m[int32(i)] = &prop{
			id:     pid,
			status: int32(proto.ITEM_STATUS_ITEM_LIVE),
			pos: info.CoordinateXYInfo{
				BaseEvent:   framework.BaseEvent{},
				CoordinateX: x,
				CoordinateY: y,
			},
		}
	}

	return m
}
