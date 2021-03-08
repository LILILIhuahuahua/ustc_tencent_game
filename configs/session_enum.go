package configs

const (
	// session的状态
	SessionStatusCreated int32 = 1
	SessionStatusDead    int32 = 0

	// session的网络类型
	SessionTypeTcp int32 = 1
	SessionTypeUdp int32 = 2
)
