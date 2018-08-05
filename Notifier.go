package event

import (
	"sync"
	"github.com/ezaurum/cthulthu/generators"
	"github.com/ezaurum/cthulthu/generators/snowflake"
)

type NotificationCallback func(event AppEvent) bool

type Notifier struct {
	mu        sync.Mutex
	Observer  map[int64]NotificationCallback
	Ch        chan AppEvent
	generator generators.IDGenerator
}

func (n *Notifier) Start() chan AppEvent {
	ch := make(chan AppEvent)
	n.Ch = ch
	go func() {
		for {
			ev := <-ch
			for _, ob := range n.Observer {
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
	n.Observer[key] = target
	n.mu.Unlock()
	return n.Ch, key
}

func (n *Notifier) Unsubscribe(key int64) {
	if _, b := n.Observer[key]; b {
		delete(n.Observer, key)
	}
}

func NewNotifier(nodeNumber int64) *Notifier {
	return &Notifier{
		mu:        sync.Mutex{},
		generator: snowflake.New(nodeNumber),
		Observer:  make(map[int64]NotificationCallback),
	}
}
