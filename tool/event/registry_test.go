package event

import (
	"fmt"
	"testing"
	"time"
)

type listenerTest struct{}

func (l *listenerTest) Handle(val interface{}) error {
	fmt.Println(fmt.Sprintf("事件1执行了 val: %s", val))
	return nil
}

type listenerTest2 struct{}

func (l *listenerTest2) Handle(val interface{}) error {
	fmt.Println(fmt.Sprintf("事件2执行了 val: %s", val))
	return nil
}

type eventTest struct{}

func (eventTest) Name() string {
	return "test"
}

func TestRegistry_Async(t *testing.T) {
	reg := NewRegistry(&Option{
		Async: true,
		Retry: []*RetryOption{
			{T: time.Second},      //第一次重试间隔
			{T: time.Second * 3},  //第二次重试间隔
			{T: time.Second * 5},  //第三次重试间隔
			{T: time.Second * 10}, //第四次重试间隔
			{T: time.Second * 10}, //第五次重试间隔
		},
	})
	reg.RegisterListener(eventTest{}, &listenerTest{})
	reg.RegisterListener(eventTest{}, &listenerTest2{})
	reg.TriggerEvent(eventTest{}, "a")
	time.Sleep(time.Millisecond * 50)
}

func TestRegistry_Sync(t *testing.T) {
	reg := NewRegistry(&Option{
		Async: false,
		Retry: []*RetryOption{
			{T: time.Second},      //第一次重试间隔
			{T: time.Second * 3},  //第二次重试间隔
			{T: time.Second * 5},  //第三次重试间隔
			{T: time.Second * 10}, //第四次重试间隔
			{T: time.Second * 10}, //第五次重试间隔
		},
	})
	reg.RegisterListener(eventTest{}, &listenerTest{})
	reg.RegisterListener(eventTest{}, &listenerTest2{})
	reg.TriggerEvent(eventTest{}, "a")
}
