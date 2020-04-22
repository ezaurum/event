package event

import "sync"

type NotifierMap struct {
	mu        sync.Mutex
	notifierMap map[string]*Notifier
}

func New() *NotifierMap {
	em := NotifierMap{
		notifierMap: make(map[string]*Notifier),
	}
	return &em
}
func (em *NotifierMap) Subscribe(eventName string, target NotificationCallback) (chan AppEvent, int64) {
	notifier := em.getNotifierInstance(eventName)
	return notifier.Subscribe(target)
}

func (em *NotifierMap) Notify(eventName string, eventData interface{}) {
	event := AppEvent{
		Name: eventName,
		Data: eventData,
	}

	notifier := em.getNotifierInstance(eventName)
	notifier.Ch <- event
}

func (em *NotifierMap) getNotifierInstance(eventName string) *Notifier {
	em.mu.Lock()
	notifier, b := em.notifierMap[eventName]
	if !b {
		notifier = NewNotifier(int64(len(eventName)), 0)
		notifier.Start()
		em.notifierMap[eventName] = notifier
	}
	em.mu.Unlock()
	return notifier
}

func (em *NotifierMap) Unsubscribe(eventName string, subID int64) {
	notifier := em.getNotifierInstance(eventName)
	notifier.Unsubscribe(subID)
}
