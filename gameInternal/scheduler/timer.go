package scheduler

import (
	"fmt"
	"log"
	"math"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

const (
	infinite = -1
)

var (
	//timerManager manage for all timers
	timerManager = &struct {
		incrementID int64 // auto increment id
		timers      map[int64]*Timer

		muClosingTimer sync.RWMutex
		closingTimer   []int64
		muCreatedTimer sync.RWMutex
		createdTimer   []*Timer
	}{}
)

type (
	//TimerFunc represents a function which wil be called periodically in main logic goroutine
	TimerFunc func()

	//TimerCondition represents a checker that returns true when cron jon needs to execute
	TimerCondition interface {
		Check(now time.Time) bool
	}

	//Time represents a cron job
	Timer struct {
		id        int64          // time id
		fn        TimerFunc      // function that execute
		createAt  int64          // timer create time
		interval  time.Duration  // execution interval
		condition TimerCondition // condition to cron job execution
		elapse    int64          // total elapse time
		closed    int32          // is timer closed
		counter   int            // counter
	}
)

func init() {
	timerManager.timers = map[int64]*Timer{}
}

//ID returns id of current timer
func (t *Timer) Id() int64 {
	return t.id
}

//Stop turns off a timer. After Stop, fn will not be called forever
func (t *Timer) Stop() {
	if atomic.AddInt32(&t.closed, 1) != 1 {
		return
	}

	t.counter = 0
}

//execute jon function with protection
func safeCall(id int64, fn TimerFunc) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("Handle timer panic: %+v\n%s", err, debug.Stack()))
		}
	}()

	fn()
}

func cron() {
	if len(timerManager.createdTimer) > 0 {
		timerManager.muCreatedTimer.Lock()
		for _, t := range timerManager.createdTimer {
			timerManager.timers[t.id] = t
		}
		timerManager.createdTimer = timerManager.createdTimer[:0]
		timerManager.muCreatedTimer.Unlock()
	}

	if len(timerManager.timers) < 1 {
		return
	}

	now := time.Now()
	unn := now.UnixNano()
	for id, t := range timerManager.timers {
		if t.counter == infinite || t.counter > 0 {
			// condition timer
			if t.condition != nil {
				if t.condition.Check(now) {
					safeCall(id, t.fn)
				}
				continue
			}

			// execute job
			if t.createAt+t.elapse <= unn {
				safeCall(id, t.fn)
				t.elapse += int64(t.interval)

				// update timer counter
				if t.counter != infinite && t.counter > 0 {
					t.counter--
				}
			}
		}

		if t.counter == 0 {
			timerManager.muClosingTimer.Lock()
			timerManager.closingTimer = append(timerManager.closingTimer, id)
			timerManager.muClosingTimer.Unlock()
			continue
		}
	}

	if len(timerManager.closingTimer) > 0 {
		timerManager.muClosingTimer.Lock()
		for _, id := range timerManager.closingTimer {
			delete(timerManager.timers, id)
		}
		timerManager.closingTimer = timerManager.closingTimer[:0]
		timerManager.muClosingTimer.Unlock()
	}
}

//NewCountTimer returns a new Timer containing a function that wil be called with a period  specified by the interval
//argument. It adjusts the intervals foe slow receivers.
//The duration must be greater than zero; if not, NewTimer will panic. Stop the timer to release associated resources
func NewTimer(interval time.Duration, fn TimerFunc) *Timer {
	return NewCountTimer(interval, infinite, fn)
}

//NewCountTimer returns a new Timer containing a function that will be called with a period specified by the
//interval argument. After count times, timer will  be stopped automatically,
//it adjusts the intervals for slow receivers. The interval must be greater than zero; if not,
//NewCountTimer will panic. Stop the timer to release associated resources.
func NewCountTimer(interval time.Duration, count int, fn TimerFunc) *Timer {
	if fn == nil {
		panic("timer: nil timer function")
	}
	if interval <= 0 {
		panic("timer: non-positive interval for NewTimer")
	}
	t := &Timer{
		id:       atomic.AddInt64(&timerManager.incrementID, 1),
		fn:       fn,
		createAt: time.Now().UnixNano(),
		interval: interval,
		elapse:   int64(interval), // first execution wil lbe after interval
		counter:  count,
	}

	timerManager.muCreatedTimer.Lock()
	timerManager.createdTimer = append(timerManager.createdTimer, t)
	timerManager.muCreatedTimer.Unlock()
	return t
}

//NewAfterTimer returns a new Timer containing a function that will be called after duration that specified by
//duration argument. The duration must be greater than zero; if not,
//NewAfterTimer will panic. Stop the timer to release associated resources.
func NewAfterTimer(duration time.Duration, fn TimerFunc) *Timer {
	return NewCountTimer(duration, 1, fn)
}

//NewCondTimer returns a new Timer containing a function that will be called when condition satisfied that specified
//by the condition argument. The duration d must be greater than zero; if not,
//NewCondTimer will panic. Stop the timer to release associated resources.
func NewCondTimer(condition TimerCondition, fn TimerFunc) *Timer {
	if condition == nil {
		panic("timer: nil condition")
	}

	t := NewCountTimer(time.Duration(math.MaxInt64), infinite, fn)
	t.condition = condition

	return t
}
