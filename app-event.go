package event

/**
프로그램 내 이벤트 전달에 사용할 클래스
 */
type AppEvent struct {
	Name    string
	Data    interface{}
}
