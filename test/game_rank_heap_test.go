package test

import (
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/game"
	"testing"
)

func TestGameRankHeap(t *testing.T) {
	heap := game.NewGameRankHeap(5)
	heroRank1 := &info.HeroRankInfo{
		HeroID:    1,
		HeroScore: 1,
	}
	heroRank2 := &info.HeroRankInfo{
		HeroID:    2,
		HeroScore: 2,
	}
	heroRank3 := &info.HeroRankInfo{
		HeroID:    3,
		HeroScore: 3,
	}
	heroRank4 := &info.HeroRankInfo{
		HeroID:    4,
		HeroScore: 4,
	}
	heroRank5 := &info.HeroRankInfo{
		HeroID:    5,
		HeroScore: 5,
	}
	//heroRank6 := &info.HeroRankInfo{
	//	HeroID: 6,
	//	HeroScore: 6,
	//}
	heap.ChallengeRank(heroRank3)
	heap.ChallengeRank(heroRank4)
	heap.ChallengeRank(heroRank1)
	heap.ChallengeRank(heroRank2)
	heap.ChallengeRank(heroRank5)
	heroRank5.HeroScore = 1
	heap.ChallengeRank(heroRank5)
	res := heap.GetSortedHeroRankInfos()
	t.Logf("%+v", res)
}
