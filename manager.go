package event

type Manager struct {
	ObserverMap map[string]*Notifier
}

func New() *Manager {
	em := Manager{
		ObserverMap: make(map[string]*Notifier),
	}
	return &em
}
func (em *Manager) Subscribe(eventName string, target NotificationCallback) (chan AppEvent, int64) {
	notifier := em.getNotifierInstance(eventName)
	return notifier.Subscribe(target)
}

func (em *Manager) Notify(eventName string, eventData interface{}) {
	event := AppEvent{
		Name: eventName,
		Data: eventData,
	}

	notifier := em.getNotifierInstance(eventName)
	notifier.Ch <- event
}

func (em *Manager) getNotifierInstance(eventName string) *Notifier {
	notifier, b := em.ObserverMap[eventName]
	if !b {
		notifier = NewNotifier(int64(len(eventName)), 100)
		notifier.Start()
		em.ObserverMap[eventName] = notifier
	}
	return notifier
}

func (em *Manager) Unsubscribe(eventName string, subID int64) {
	notifier := em.getNotifierInstance(eventName)
	notifier.Unsubscribe(subID)
}
