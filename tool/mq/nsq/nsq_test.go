package nsq

import (
	"context"
	"github.com/spf13/cast"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	err := InitConsumer("test", "test", "172.16.49.184:4150", func(c <-chan *Elem) {
		for {
			select {
			case msg, ok := <-c:
				if !ok {
					return
				}
				t.Log("get message :", string(msg.Payload))
			}
		}
	})
	if err != nil {
		t.Log(err)
	}
	select {}
}

func TestProducer(t *testing.T) {
	pusher := make(chan []byte)
	var i int
	go func() {
		for {
			i++
			pusher <- []byte("test message: " + cast.ToString(i))

			t.Log("push message down,times: ", i, "\n")
			time.Sleep(time.Second)
		}
	}()
	err := InitProducer(context.Background(), "172.16.49.184:4150", "test", pusher)
	if err != nil {
		t.Error(err)
	}
	select {}
}

func TestProducer2(t *testing.T) {
	pusher := make(chan []byte)
	var i int
	go func() {
		for {
			i++
			pusher <- []byte("test message: " + cast.ToString(i))
			t.Log("push message down,times: ", i, "\n")
			time.Sleep(time.Second)
		}
	}()
	err := InitProducer(context.Background(), "172.16.49.184:4150", "develop_PUSER_MESSAGE_center_ZGV2ZWxvcD", pusher)
	if err != nil {
		t.Error(err)
	}
	select {}
}

func TestProducer3(t *testing.T) {
	pusher := make(chan []byte)
	var i int
	go func() {
		for {
			i++
			pusher <- []byte("test message: " + cast.ToString(i))
			t.Log("push message down,times: ", i, "\n")
			time.Sleep(time.Second)
		}
	}()
	err := InitProducer(context.Background(), "172.16.49.184:4150", "test4", pusher)
	if err != nil {
		t.Error(err)
	}
	select {}
}

func TestConsumer2(t *testing.T) {
	err := InitConsumer("develop_PUSER_MESSAGE_center_ZGV2ZWxvcD", "develop_PUSER_MESSAGE_center_ZGV2ZWxvcD", "172.16.49.184:4150", func(c <-chan *Elem) {
		for {
			select {
			case msg, ok := <-c:
				if !ok {
					return
				}
				t.Log("get message :", string(msg.Payload))
			}
		}
	})
	if err != nil {
		t.Log(err)
	}
	select {}
}

func TestConsumer3(t *testing.T) {
	err := InitConsumerLookup("develop_PUSER_MESSAGE_center_ZGV2ZWxvcD", "develop_PUSER_MESSAGE_center_ZGV2ZWxvcD", "172.16.49.184:4161", func(c <-chan *Elem) {
		for {
			select {
			case msg, ok := <-c:
				if !ok {
					return
				}
				t.Log("get message :", string(msg.Payload))
			}
		}
	})
	if err != nil {
		t.Log(err)
	}
	select {}
}

func TestConsumer4(t *testing.T) {
	err := InitConsumerLookup("test4", "test", "172.16.49.184:4161", func(c <-chan *Elem) {
		for {
			select {
			case msg, ok := <-c:
				if !ok {
					return
				}
				t.Log("get message :", string(msg.Payload))
			}
		}
	})
	if err != nil {
		t.Log(err)
	}
	select {}
}
