package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/senyu-up/toolbox/tool/str"
)

var qclient *RedisQ

type ExampleFamily struct {
	Username string `json:"username"`
	Age      int    `json:"age"`
}

// 工作方法 Process 的实现
func work(topic string, data []byte) error {
	var result ExampleFamily
	err := json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

var ctx, cancle = context.WithCancel(context.Background())
var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
		DB:      2,
	})

	qclient = NewRedisMQClient(ctx, client, time.Second*10, IQueueTopic{
		QueueScheme: "gift:start",
		HashValue:   16,
	}, work)
}

func TestSub(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
		DB:      2,
	})

	qclient = NewRedisMQClient(context.Background(), client, time.Second*10, IQueueTopic{
		QueueScheme: "gift:start",
		HashValue:   16,
	}, work)
	//发送请求，并将数据推入异步队列
	http.HandleFunc("/push", func(writer http.ResponseWriter, request *http.Request) {
		b, _ := io.ReadAll(request.Body)

		if err := qclient.Push(request.RequestURI, b); err != nil {
			writer.Write(str.StringToByte(err.Error()))
			return
		}
	})
	//获取队列状态信息
	http.HandleFunc("/queue/status", func(writer http.ResponseWriter, request *http.Request) {
		for _, v := range qclient.QueueClient() {
			w := fmt.Sprintf("%s:%d\n", v.QueueName, v.Length())
			writer.Write([]byte(w))
		}
		return
	})
	http.ListenAndServe(":8081", nil)

	//写入数据
	//go func() {
	//	for i:=0;i<10000;i++{
	//		qclient.Push(fmt.Sprintf("nd98h9823n%d",i),[]byte("3298nijdewihjd98ewh9d8ewd"))
	//	}
	//}()
	//time.Sleep(1*time.Hour)
}

func TestReqPush(t *testing.T) {
	param := ExampleFamily{
		Username: "abc",
		Age:      23,
	}
	b, _ := json.Marshal(param)
	buf := new(bytes.Buffer)
	buf.Write(b)
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8081/push", buf)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	t.Log(resp, err)
}

func TestQueue(t *testing.T) {

	//for i:=0;i<1000000;i++{
	//	qclient.Push(fmt.Sprintf("%d",i),[]byte("ghelloe"))
	//}
	//time.Sleep(13*time.Second)
	//s:=make(chan os.Signal,1)
	//signal.Notify(s,os.Interrupt,os.Kill)
	//<-s
	//cancle()
	//qclient.Stop()
	client.RPush("name", "klm23")
	var w sync.WaitGroup
	w.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(client.LPop("name").Result())
			fmt.Println(time.Now().UnixNano())
			w.Done()
		}()
	}
	w.Wait()

}
