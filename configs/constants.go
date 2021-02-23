package configs

import "time"

var (
	// Map info
	MapMinX float32 = -493.2
	MapMaxX float32 = 1507
	MapMinY float32 = -361
	MapMaxY float32 = 1118

	// Prop max count in map
	MaxPropCountInMap int = 50

	// GlobalInfoNotify interval
	GlobalInfoNotifyInterval time.Duration = 5000 * time.Millisecond

	// Server addr
	ServerAddr = "127.0.0.1:8888"

	// debug mode
	Debug bool = false
)
