package scheduler

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync/atomic"
	"time"
)

const (
	messageQueueBacklog = 1 << 10
	sessionCloseBacklog = 1 << 8
	taskSize            = 1 << 8
)

//LocalScheduler schedules task to a customized goroutine
type LocalScheduler interface {
	Schedule(Task)
}

type Task func()
type Hook func()

var (
	chDie   = make(chan struct{})
	chExit  = make(chan struct{})
	chTasks = make(chan Task, taskSize)
	started int32
	closed  int32
)

func try(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("Handle messssage panic: %+v\n%s", err, debug.Stack()))
		}
	}()

	fn()
}

//Sched schedule task every duration interval. duration must be greater than zero,; if not, Sched will panic.
func Sched(duration time.Duration) {
	if duration <= 0 {
		panic("scheduler: non-positive interval for scheduler")

	}
	if atomic.AddInt32(&started, 1) != 1 {
		return
	}

	ticker := time.NewTicker(duration)
	defer func() {
		ticker.Stop()
		close(chExit)
	}()

	for {
		select {
		case <-ticker.C:
			cron()
		case f := <-chTasks:
			try(f)
		case <-chDie:
			return
		}
	}
}

func Close() {
	if atomic.AddInt32(&closed, 1) != 1 {
		return
	}
	close(chDie)
	<-chExit
	log.Println("Scheduler stopped")
}

func PushTask(task Task) {
	chTasks <- task
}
