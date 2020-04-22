package owlbear

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

/**
프로그램 내 이벤트 전달에 사용할 클래스
*/
func TestEventManager_Start(t *testing.T) {

	em := Notifier{}
	ch := em.Start()

	if nil == ch {
		t.Fail()
	}
}

type ttt struct {
	Name   string
	Events []Event
	wg     *sync.WaitGroup
}

func (t *ttt) OnNotified(event Event) bool {
	t.Events = append(t.Events, event)
	if t.wg != nil  {
		t.wg.Done()
	}
	return true
}

func TestEventManager_WaitEvent(t *testing.T) {

	em := Notifier{}
	ch := em.Start()

	if nil == ch {
		t.Fail()
	}

	ch <- Event{
		Name: "Test",
		Data: "WTF",
	}
}

func TestEventManager_Observe(t *testing.T) {

	eventName := "Test"
	em := NewNotifier(0, 1)
	ch := em.Start()

	if nil == ch {
		t.Fail()
	}

	i := &ttt{
		Name: "ttdfs",
	}
	em.Subscribe(i.OnNotified)
	ch <- Event{
		Name: eventName,
		Data: "WTF",
	}

	time.Sleep(time.Second)
}

func TestEventManger_Observe(t *testing.T) {
	em := New()
	var wg sync.WaitGroup
	test1 := ttt{
		Name: "receive Test1",
		wg: &wg,
	}
	test2 := ttt{
		Name: "receive Test2",
		wg: &wg,
	}
	test1Ch, _ := em.Subscribe("Test1", test1.OnNotified)
	_, _ = em.Subscribe("Test1", test2.OnNotified)

	test2Ch, _ := em.Subscribe("Test2", test2.OnNotified)

	test3Ch, _ := em.Subscribe("Test3", test2.OnNotified)

	wg.Add(4)
	test1Ch <- Event{
		Name: "Test1",
		Data: "WTF",
	}
	test2Ch <- Event{
		Name: "Test2",
		Data: "WTF",
	}
	test3Ch <- Event{
		Name: "Test3",
		Data: "WTF",
	}

	wg.Wait()

	assert.Equal(t, 1, len(test1.Events))
	assert.Equal(t, 3, len(test2.Events))

}

func TestEventManager_Unsubscribe(t *testing.T) {
	em := New()
	var wg sync.WaitGroup
	eventName := "Test1"
	test1 := ttt{
		Name: "receiver1",
		wg:   &wg,
	}
	test2 := ttt{
		Name: "receiver2",
		wg:   &wg,
	}
	test1Ch, subID := em.Subscribe(eventName, test1.OnNotified)
	_, _ = em.Subscribe(eventName, test2.OnNotified)
	wg.Add(2)
	test1Ch <- Event{
		Name: eventName,
		Data: "WTF",
	}

	wg.Wait()
	assert.Equal(t, 1, len(test1.Events))
	assert.Equal(t, 1, len(test2.Events))

	wg.Add(1)
	em.Unsubscribe(eventName, subID)
	test1Ch <- Event{
		Name: eventName,
		Data: "WTF2",
	}
	wg.Wait()

	assert.Equal(t, 1, len(test1.Events))
	assert.Equal(t, 2, len(test2.Events))
}
