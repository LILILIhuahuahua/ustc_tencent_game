package game

import "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"

type GameRankHeap struct {
	Size     int32
	Capacity int32
	rankHeap []*info.HeroRankInfo
}

func NewGameRankHeap(size int32) *GameRankHeap {
	return &GameRankHeap{
		Size:     size,
		Capacity: 0,
		rankHeap: make([]*info.HeroRankInfo, size+1),
	}
}

func CopyFromGameRankHeap(heap *GameRankHeap) *GameRankHeap {
	copyHeap := &GameRankHeap{
		Size:     heap.Size,
		Capacity: heap.Capacity,
		rankHeap: make([]*info.HeroRankInfo, heap.Size+1),
	}
	for i := 1; i <= int(heap.Capacity); i++ {
		copyHeap.rankHeap[i] = heap.rankHeap[i]
	}
	return copyHeap
}

func (h *GameRankHeap) isEmpty() bool {
	return 0 == h.Capacity
}

func (h *GameRankHeap) isFull() bool {
	return h.Size == h.Capacity
}

func (h *GameRankHeap) getNextIndex() int32 {
	if !h.isFull() {
		return h.Capacity + 1
	}
	return -1
}

func (h *GameRankHeap) ContainsRank(rank *info.HeroRankInfo) bool {
	return !(0 == h.GetRankIndex(rank))
}

func (h *GameRankHeap) GetRankIndex(target *info.HeroRankInfo) int32 {
	idx := 0
	if h.isEmpty() {
		return int32(idx)
	}
	var l int = int(h.Capacity)
	for i := 1; i <= l; i++ {
		if target.HeroID == h.rankHeap[i].HeroID {
			idx = i
			break
		}
	}
	return int32(idx)
}

func (h *GameRankHeap) InsertRank(rank *info.HeroRankInfo) {
	newIdx := h.getNextIndex()
	h.rankHeap[newIdx] = rank
	h.Capacity++
	h.upAdjust(newIdx)
}

func (h *GameRankHeap) ChallengeRank(challenger *info.HeroRankInfo) bool {
	if h.isEmpty() { //case1：堆为空，直接插入挑战者
		h.InsertRank(challenger)
		return true
	} else {
		if h.ContainsRank(challenger) { //case2：挑战者在堆中，更新挑战者的值重新维护堆
			idx := h.GetRankIndex(challenger)
			h.rankHeap[idx] = challenger
			return h.Adjust()
		} else {
			if h.isFull() { //case3：堆满且挑战者不在堆中，比较挑战者与堆顶者，若挑战者分数大则胜出
				if challenger.HeroScore > h.rankHeap[1].HeroScore {
					h.rankHeap[1] = challenger
					h.Adjust()
					return true
				}
			} else { //case4：堆不满且挑战者不在堆中，插入挑战者
				h.InsertRank(challenger)
				return true
			}
		}
	}
	return false
}

func (h *GameRankHeap) GetSortedHeroRankInfos() []info.HeroRankInfo {
	copyHeap := CopyFromGameRankHeap(h)
	for i := copyHeap.Capacity; i > 1; i-- {
		temp := copyHeap.rankHeap[1]
		copyHeap.rankHeap[1] = copyHeap.rankHeap[i]
		copyHeap.rankHeap[i] = temp
		copyHeap.rankHeap[i].HeroRank = i
		copyHeap.downAdjust(1, i-1)
	}
	if !copyHeap.isEmpty() {
		copyHeap.rankHeap[1].HeroRank = 1
	}
	res := make([]info.HeroRankInfo, copyHeap.Capacity)
	for i := 1; i <= int(copyHeap.Capacity); i++ {
		res[i-1] = *copyHeap.rankHeap[i]
	}
	return res
}

func (h *GameRankHeap) downAdjust(i int32, max int32) bool {
	var isSwap bool = false
	childIdx := 2 * i
	for childIdx <= max {
		if childIdx+1 <= max && h.rankHeap[childIdx+1].HeroScore < h.rankHeap[childIdx].HeroScore {
			childIdx = childIdx + 1
		}
		if h.rankHeap[childIdx].HeroScore < h.rankHeap[i].HeroScore {
			temp := h.rankHeap[i]
			h.rankHeap[i] = h.rankHeap[childIdx]
			h.rankHeap[childIdx] = temp
			i = childIdx
			childIdx = 2 * i
			isSwap = true
		} else {
			break
		}
	}
	return isSwap
}

func (h *GameRankHeap) upAdjust(i int32) bool {
	var isSwap bool = false
	parentIdx := i / 2
	for parentIdx >= 1 {
		if h.rankHeap[i].HeroScore < h.rankHeap[parentIdx].HeroScore {
			temp := h.rankHeap[i]
			h.rankHeap[i] = h.rankHeap[parentIdx]
			h.rankHeap[parentIdx] = temp
			i = parentIdx
			parentIdx = i / 2
			isSwap = true
		} else {
			break
		}
	}
	return isSwap
}

func (h *GameRankHeap) Adjust() bool {
	isSwap := false
	for i := h.Capacity / 2; i >= 1; i-- {
		res := h.downAdjust(i, h.Capacity)
		if res {
			isSwap = res
		}
	}
	return isSwap
}
