package configs

import (
	"math"
	"time"
)

var (
	// Map info
	MapMinX float32 = -493.2
	MapMaxX float32 = 1507
	MapMinY float32 = -361
	MapMaxY float32 = 1118

	// Aoi coordinate
	TileSize      float32 = 1                                                              //地图瓦片大小，一个瓦片对应一个地图坐标
	TowerRadius   float32 = 200                                                            // 灯塔AOI半径
	TowerDiameter float32 = TowerRadius * 2                                                //灯塔AOI直径
	PlayerRange   float32 = 100                                                            // 玩家视野半径
	TowerCols     int32   = int32(math.Ceil(float64((MapMaxX - MapMinX) / TowerDiameter))) // 整个地图中有多少列Tower 从1开始
	TowerRows     int32   = int32(math.Ceil(float64((MapMaxY - MapMinY) / TowerDiameter))) // 整个地图中有多少行Tower 从1开始

	// Collision-Check info
	MaxObjectNum int32 = 5
	MaxLevelNum  int32 = 5

	// Hero info
	HeroSizeGrowthStep float32 = 5.0   // 英雄吃道具以后size增长步长
	HeroSpeedSlowStep  float32 = 0.5   // 英雄吃道具以后速度减缓步长
	HeroSizeUpLimit    float32 = 200.0 // 英雄size上限
	HeroSpeedDownLimit float32 = 10.0  // 英雄速度下限
	HeroEatItemBonus   int32   = 10
	HeroEatEnemyBonus  int32   = 50
	HeroRankListLength int32   = 10

	// Game info
	GameWinLiminationScore int32 = 300 //对局优胜分数值
	MinMatchingBatchSessionNum int32 = 1 //最小批量匹配会话数量
	MatchingWaitOverTime int64 = 10 //匹配等待超时时间（单位：s ）
	GameAliveHeroLimit int32 = 10 //房间最大人数限制

	// Prop max count in map
	MaxPropCountInMap int = 100

	// GlobalInfoNotify interval
	GlobalInfoNotifyInterval time.Duration = 5000 * time.Millisecond
	GlobalScheduleConfig     time.Duration = 20 * time.Millisecond

	// Server info
	ServerAddr              = "0.0.0.0:8888"
	MaxEventQueueSize int32 = 500

	// RemoteCLB is LoadBalancer address used for dgs
	RemoteCLB = "150.158.216.120:32001"
	//PodIP = "1.116.109.211:8888"

	// debug mode
	Debug bool = false

	// Hero msg
	HeroInitSize       float32 = 45
	SizeBig            float32 = 5
	HeroMoveSpeed      float32 = 100
	HeroInitPositionX  float32 = 0
	HeroInitPositionY  float32 = 0
	HeroInitDirectionX float32 = 1
	HeroInitDirectionY float32 = 0

	// props
	PropRadius     float32 = 20 // Initial prop radius
	PropStatusLive int32   = 0
	PropStatusDead int32   = 1

	// 道具效果
	PropInvincibleTimeMax = int64(time.Second * 5) // 最长无敌时间
	PropSpeedUpTimeMax = int64(time.Second * 5) // 最长加速时间

	// mongodb
	MongoURI      string = ""
	MongoPoolSize uint64 = 100
	DBName        string = "happyball"

	// DBProxy configuration
	DBProxyAddr string = ""
)