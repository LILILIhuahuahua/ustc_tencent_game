package tools

import "time"

var TIME_UTIL = &TimeUtil{}

type TimeUtil struct {
}

func (t *TimeUtil) NowSecs() int64 {
	return time.Now().Unix()
}

func (t *TimeUtil) NowMillis() int64 {
	return time.Now().UnixNano() / 1e6
}

func (t *TimeUtil) NowNanos() int64 {
	return time.Now().UnixNano()
}
