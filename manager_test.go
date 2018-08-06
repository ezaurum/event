package event

import (
	"testing"
	"fmt"
	"time"
	"github.com/stretchr/testify/assert"
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
	Events []AppEvent
}

func (t *ttt) Notify(event AppEvent) bool {
	fmt.Printf("%v, get %v\n", t.Name, event)
	t.Events = append(t.Events, event)
	return true
}

func TestEventManager_WaitEvent(t *testing.T) {

	em := Notifier{}
	ch := em.Start()

	if nil == ch {
		t.Fail()
	}

	ch <- AppEvent{
		Name: "Test",
		Data: "WTF",
	}
}

func TestEventManager_Observe(t *testing.T) {

	eventName := "Test"
	em := NewNotifier(0)
	ch := em.Start()

	if nil == ch {
		t.Fail()
	}

	i := &ttt{
		Name: "ttdfs",
	}
	em.Subscribe(i.Notify)
	ch <- AppEvent{
		Name: eventName,
		Data: "WTF",
	}

	time.Sleep(time.Second)
}

func TestEventManger_Observe(t *testing.T) {
	em := New()
	test1 := ttt{
		Name: "Test1",
	}
	test2 := ttt{
		Name: "Test2",
	}
	test1Ch, _ := em.Subscribe("Test1", test1.Notify)
	test2Ch, _ := em.Subscribe("Test2", test2.Notify)
	test3Ch, _ := em.Subscribe("Test3", test2.Notify)

	test1Ch <- AppEvent{
		Name: "Test1",
		Data: "WTF",
	}
	test2Ch <- AppEvent{
		Name: "Test2",
		Data: "WTF",
	}
	test3Ch <- AppEvent{
		Name: "Test3",
		Data: "WTF",
	}
	time.Sleep(time.Second)

	assert.Equal(t, 1, len(test1.Events))
	assert.Equal(t, 2, len(test2.Events))

}

func TestEventManager_Unsubscribe(t *testing.T) {
	em := New()
	test1 := ttt{
		Name: "Test1",
	}
	test2 := ttt{
		Name: "Test1",
	}
	test1Ch, subID := em.Subscribe("Test1", test1.Notify)
	_,_ = em.Subscribe("Test1", test2.Notify)
	test1Ch <- AppEvent{
		Name: "Test1",
		Data: "WTF",
	}
	em.Unsubscribe("Test1", subID)
	test1Ch <- AppEvent{
		Name: "Test1",
		Data: "WTF",
	}

	time.Sleep(time.Second)
	assert.Equal(t, 1, len(test1.Events))
	assert.Equal(t, 2, len(test2.Events))

}
