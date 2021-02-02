package model

import "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"

type Hero struct {
	ID          int32
	Status      int32
	Size        float32
	Speed       float32
	HeroDirection Coordinate
	HeroPosition Coordinate
}

func (h *Hero)ToEvent() info.HeroInfo {
	return info.HeroInfo{
		ID: h.ID,
		Status: h.Status,
		Speed: h.Speed,
		Size: h.Size,
		HeroDirection: h.HeroDirection.ToEvent(),
		HeroPosition: h.HeroPosition.ToEvent(),
	}
}

func (h *Hero)FromEvent(info info.HeroInfo)  {
	h.ID = info.ID
	h.Status = info.Status
	h.Speed = info.Speed
	h.Size = info.Size
	h.HeroDirection = Coordinate{}
	h.HeroDirection.FromEvent(info.HeroDirection)
	h.HeroPosition = Coordinate{}
	h.HeroPosition.FromEvent(info.HeroPosition)
}