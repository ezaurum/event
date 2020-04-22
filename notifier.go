package event

import (
	"github.com/ezaurum/cthulthu/generators"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"sync"
)

type NotificationCallback func(event AppEvent) bool

type Notifier struct {
	mu        sync.Mutex
	observer  map[int64]NotificationCallback
	Ch        chan AppEvent
	generator generators.IDGenerator
	buffer    int
}

func (n *Notifier) Start() chan AppEvent {
	ch := make(chan AppEvent, n.buffer)
	n.Ch = ch
	go func() {
		for {
			ev := <-ch
			for _, ob := range n.observer {
				if !ob(ev) {
					break
				}
			}
		}
	}()
	return ch
}

func (n *Notifier) Subscribe(target NotificationCallback) (chan AppEvent, int64) {
	n.mu.Lock()
	key := n.generator.GenerateInt64()
	//fmt.Printf("subscribe %d %p\n",key, target )
	n.observer[key] = target
	n.mu.Unlock()
	return n.Ch, key
}

func (n *Notifier) Unsubscribe(key int64) {
	n.mu.Lock()
	if _, b := n.observer[key]; b {
		delete(n.observer, key)
	}
	n.mu.Unlock()
}

func NewNotifier(nodeNumber int64, buffer int) *Notifier {
	if buffer < 0 {
		// 자동 설정
		buffer = 10
	}
	return &Notifier{
		generator: snowflake.New(nodeNumber),
		observer:  make(map[int64]NotificationCallback),
		buffer:    buffer,
	}
}
