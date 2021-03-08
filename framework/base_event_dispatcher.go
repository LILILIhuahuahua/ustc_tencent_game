package framework

import "github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"

type BaseEventDispatcher struct {
}

func (b BaseEventDispatcher) FireEvent(e event.Event) {
	//todo:改造为线程池 而非无脑开协程
	EVENT_HANDLER.onEvent(e)
}

func (b BaseEventDispatcher) FireEventToSession(e event.Event, s event.Session) {

}

func (b BaseEventDispatcher) Close() {

}
