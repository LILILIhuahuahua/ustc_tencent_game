package event

type EventRingQueue struct {
	index int32
	maxSize int32
	eventMap map[int32]Event
	eventChan chan int32
}

func NewEventRingQueue(maxSize int32) *EventRingQueue {
	return &EventRingQueue{
		eventMap: make(map[int32]Event, maxSize),
		eventChan: make(chan int32, maxSize-2),
		index: 0,
		maxSize: maxSize,
	}
}

func (queue *EventRingQueue) Push(event Event)  {
	if nil != event {
		queue.eventChan <- queue.index
		queue.eventMap[queue.index] = event
		queue.index++
		if queue.index >= queue.maxSize {
			queue.index = 0
		}
	}
}

func (queue *EventRingQueue) Pop() Event{
	idx := <- queue.eventChan
	return queue.eventMap[idx]
}


