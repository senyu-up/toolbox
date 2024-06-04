package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cast"
	"testing"
	"time"
)

func TestPublish(t *testing.T) {
	ci, err := InitNats([]string{"nats://127.0.0.1:4222"})
	if err != nil {
		t.Log(err)
	}
	for true {
		err = ci.client.Publish("test", []byte("the world."))
		if err != nil {
			t.Log(err)
		}
		fmt.Println("push down.")
		time.Sleep(time.Millisecond)
	}
}

func TestSubscribe(t *testing.T) {
	ci, err := InitNats([]string{"nats://127.0.0.1:4222"})
	if err != nil {
		t.Log(err)
	}
	_, err = ci.client.Subscribe("test", func(msg *nats.Msg) {
		fmt.Println(string(msg.Data))
	})
	if err != nil {
		t.Log(err)
	}
	select {}
}

func TestSubscribe2(t *testing.T) {
	ci, err := InitNats([]string{"nats://127.0.0.1:4222"})
	if err != nil {
		t.Log(err)
	}

	_, err = ci.client.Subscribe("test", func(msg *nats.Msg) {
		fmt.Println(string(msg.Data))
	})
	if err != nil {
		t.Log(err)
	}
	select {}
}

func TestNatsHandler(t *testing.T) {
	ci, err := InitNats([]string{"nats://127.0.0.1:4222"})
	if err != nil {
		t.Log(err)
	}
	err = ci.Subscribe("test", func(data *nats.Msg) {
		fmt.Println(string(data.Data))
	})
	if err != nil {
		t.Log(err)
	}
	select {}
}

func TestQueue1(t *testing.T) {
	ci, err := InitNats([]string{"nats://127.0.0.1:4222"})
	if err != nil {
		t.Log(err)
	}
	err = ci.Pop("test", "queue1", func(data *nats.Msg) {
		fmt.Println("queue1" + string(data.Data))
	})
	if err != nil {
		t.Log(err)
	}
	select {}
}

func TestQueue2(t *testing.T) {
	ci, err := InitNats([]string{"nats://127.0.0.1:4222"})
	if err != nil {
		t.Log(err)
	}
	err = ci.Pop("test", "queue1", func(data *nats.Msg) {
		fmt.Println("queue1" + string(data.Data))
	})
	if err != nil {
		t.Log(err)
	}
	select {}
}

func TestQueue3(t *testing.T) {
	ci, err := InitNats([]string{"nats://127.0.0.1:4222"})
	if err != nil {
		t.Log(err)
	}
	err = ci.Pop("test", "queue2", func(data *nats.Msg) {
		fmt.Println("queue2" + string(data.Data))
	})
	if err != nil {
		t.Log(err)
	}
	select {}
}

func TestRequestWithTimeout(t *testing.T) {
	ci, err := InitNats([]string{"nats://127.0.0.1:4222"})
	if err != nil {
		t.Log(err)
	}
	err = ci.RequestWithTimeout("test.request", []byte("hello!"), func(data *nats.Msg) {
		fmt.Println("get test.request response:", string(data.Data))
	}, time.Second*3)
	if err != nil {
		t.Log(err)
	}
}

func TestRequest(t *testing.T) {
	ci, err := InitNats([]string{"nats://127.0.0.1:4222"})
	if err != nil {
		t.Log(err)
	}
	var i int
	for {
		go func(j int) {
			ch, err := ci.Request("test", []byte("hello!"+cast.ToString(j)))
			if err != nil {
				t.Log(err)
			}
			msg := <-ch
			fmt.Println(string(msg.Data))
		}(i)
		i++
		time.Sleep(time.Millisecond * 100)
	}
}

func TestResponse(t *testing.T) {
	ci, err := InitNats([]string{"nats://127.0.0.1:4222"})
	if err != nil {
		t.Log(err)
	}
	err = ci.Response("test.request", func(request *nats.Msg) []byte {
		fmt.Println("get request..." + string(request.Data))
		time.Sleep(time.Second * 1)
		return append([]byte("response1"), request.Data...)
	})
	if err != nil {
		t.Log(err)
	}
	select {}
}

func TestResponse2(t *testing.T) {
	ci, err := InitNats([]string{"nats://127.0.0.1:4222"})
	if err != nil {
		t.Log(err)
	}
	err = ci.Response("test.request", func(request *nats.Msg) []byte {
		fmt.Println("get request..." + string(request.Data))
		time.Sleep(time.Second * 2)
		return append([]byte("response2"), request.Data...)
	})
	if err != nil {
		t.Log(err)
	}
	select {}
}
