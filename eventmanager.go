package event

import "sync"

type NotificationCallback func(event AppEvent) bool

type Notifier struct {
	mu sync.Mutex
	Observer []NotificationCallback
	Ch       chan AppEvent
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

func (n *Notifier) Subscribe(target NotificationCallback) chan AppEvent {
	n.mu.Lock()
	n.Observer = append(n.Observer, target)
	n.mu.Unlock()
	return n.Ch
}

//TODO unsubscribe
/*
func (n *Notifier) Unsubscribe(target NotificationCallback) {

		vsf := make([]NotificationCallback, 0)
		for _, v := range n.Observer {
			if target == v {
				continue
			}
		}
		if f(v) {
		vsf = append(vsf, v)
	}
		return vsf

}*/

type EventManger struct {
	ObserverMap map[string]*Notifier
}

func NewEventManager() *EventManger {
	em := EventManger{
		ObserverMap: make(map[string]*Notifier),
	}
	return &em
}
func (em *EventManger) Subscribe(eventName string, target NotificationCallback) chan AppEvent {
	notifier := em.getNotifierInstance(eventName)
	return notifier.Subscribe(target)
}

func (em *EventManger) Notify(eventName string, eventData interface{}) {
	event := AppEvent{
		Name: eventName,
		Data: eventData,
	}

	notifier := em.getNotifierInstance(eventName)
	notifier.Ch <- event
}

func (em *EventManger) getNotifierInstance(eventName string) *Notifier {
	notifier, b := em.ObserverMap[eventName]
	if !b {
		notifier = &Notifier{
			mu:sync.Mutex{},
		}
		notifier.Start()
		em.ObserverMap[eventName] = notifier
	}
	return notifier
}
