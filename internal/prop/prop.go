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
	ErrNilProp        = errors.New("nil Prop value")
	ErrPropNotExist   = errors.New("Prop not exist in propManger")
	ErrPropDuplicate  = errors.New("Prop duplicate in propManger")
	ErrNilPropManager = errors.New("nil propManger")
)

type Prop struct {
	id     int32
	status int32
	pos    info.CoordinateXYInfo
	towerId int32
	//radius float32
}

type PropsManger struct {
	mu    *sync.RWMutex
	props map[int32]*Prop
}

// ID returns the id of Prop
func (p *Prop) ID() int32 {
	return p.id
}

// Status returns the status (alive or dead) of Prop
func (p *Prop) Status() int32 {
	return p.status
}

func (p *Prop) SetStatus(status int32) {
	p.status = status
}

// GetX returns the x coordinate of Prop
func (p *Prop) GetX() float32 {
	return p.pos.CoordinateX
}

// GetY returns the y coordinate of Prop
func (p *Prop) GetY() float32 {
	return p.pos.CoordinateY
}

func (p *Prop) SetTowerId(towerId int32) {
	p.towerId = towerId
}

func (p *Prop) GetTowerId() int32 {
	return p.towerId
}

// New return an instance of propsManger, which contains many props
func New() *PropsManger {
	return &PropsManger{
		mu:    &sync.RWMutex{},
		props: newProps(configs.MapMinX, configs.MapMaxX, configs.MapMinY, configs.MapMaxY,
			configs.MaxPropCountInMap),
	}
}

// GetProps return all props in propManager
func (p *PropsManger) GetProps() ([]*Prop, error) {
	if p == nil {
		return nil, ErrNilPropManager
	}

	var props []*Prop
	for _, v := range p.props {
		props = append(props, v)
	}
	return props, nil
}

func (p *PropsManger) GetProp(id int32) (*Prop, error) {
	if p == nil {
		return nil, ErrNilPropManager
	}
	var prop *Prop
	prop = p.props[id]
	return prop, nil
}

// AddProp add Prop to propManger,if Prop is nil or Prop has existed on propManager, it will return error.
func (p *PropsManger) AddProp(pr *Prop) error {
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
func newProps(minX float32, maxX float32, minY float32, maxY float32, count int) map[int32]*Prop {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	m := make(map[int32]*Prop, count)
	for i := 0; i < count; i++ {
		x := minX + r.Float32()*(maxX-minX)
		y := minY + r.Float32()*(maxY-minY)
		pid := int32(guuid.New().ID())

		m[pid] = &Prop{
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
