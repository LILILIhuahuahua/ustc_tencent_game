package event

import (
	"fmt"
	"github.com/pkg/errors"
	"math"
	"sync"
)

type EventRingQueue struct {
	rear     int32           //头指针，指向队列开头元素
	head     int32           //尾指针，指向下一次待插入元素的位置
	size     int32           //队列大小
	eventMap map[int32]Event //消息集合，key为下标，value为消息
	lock     sync.Mutex
}

func NewEventRingQueue(maxSize int32) *EventRingQueue {
	return &EventRingQueue{
		eventMap: make(map[int32]Event, maxSize+1),
		head:     0,
		rear:     0,
		size:     maxSize + 1,
	}
}

func (queue *EventRingQueue) IsEmpty() bool {
	res := queue.rear == queue.head
	return res
}

func (queue *EventRingQueue) IsFull() bool {
	res := (queue.rear+1)%(queue.size) == queue.head
	return res
}

func (queue *EventRingQueue) Pop() (Event, error) {
	queue.lock.Lock()
	if queue.IsEmpty() {
		queue.lock.Unlock()
		return nil, errors.New("[EventRingQueue]消息队列为空，无法读出新的消息！")
	}
	event := queue.eventMap[queue.head]
	queue.eventMap[queue.head] = nil
	queue.head++
	fmt.Printf("[EventRingQueue]取出消息成功，当前容量：%v \n", queue.capacity())
	queue.lock.Unlock()
	return event, nil
}

func (queue *EventRingQueue) Push(event Event) error {
	queue.lock.Lock()
	if queue.IsFull() {
		queue.lock.Unlock()
		return errors.New("[EventRingQueue]消息队列已满，无法继续添加消息！")
	}
	queue.eventMap[queue.rear] = event
	queue.rear++
	fmt.Printf("[EventRingQueue]插入消息成功，当前容量：%v \n", queue.capacity())
	queue.lock.Unlock()
	return nil
}

func (queue *EventRingQueue) capacity() int32 {
	res := int32(math.Abs(float64(queue.rear - queue.head)))
	return res
}

//func (queue *EventRingQueue) Push(event Event)  {
//	if nil != event {
//		queue.eventChan <- queue.index
//		queue.eventMap[queue.index] = event
//		queue.index++
//		if queue.index >= queue.maxSize {
//			queue.index = 0
//		}
//	}
//}
//
//func (queue *EventRingQueue) Pop() Event{
//	idx := <- queue.eventChan
//	return queue.eventMap[idx]
//}
