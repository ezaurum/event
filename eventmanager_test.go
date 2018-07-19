package common

import (
	"testing"
	"fmt"
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
	Name string
}

func (t *ttt) Notify(event AppEvent) bool {
	fmt.Printf("%v, get %v\n", t.Name, event)
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

	em := Notifier{}
	ch := em.Start()

	if nil == ch {
		t.Fail()
	}

	i := &ttt{
		Name: "ttdfs",
	}
	em.Subscribe(i.Notify)
	ch <- AppEvent{
		Name: "Test",
		Data: "WTF",
	}

	time.Sleep(time.Second)
}

func TestEventManger_Observe(t *testing.T) {
	em := NewEventManager()
	test1 := ttt{
		"Test1",
	}
	test2 := ttt{
		"Test2",
	}
	test1Ch := em.Subscribe("Test1", test1.Notify)
	test2Ch := em.Subscribe("Test2", test2.Notify)
	test3Ch := em.Subscribe("Test3", test2.Notify)
	em.Subscribe("Test3", test1.Notify)

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
}
