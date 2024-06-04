package test

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	broadcast2 "github.com/senyu-up/toolbox/tool/broadcast"
	"github.com/senyu-up/toolbox/tool/broadcast/adapter"
	config2 "github.com/senyu-up/toolbox/tool/config"
	"testing"
	"time"
)

func getRedis() redis.UniversalClient {
	type AppConf struct {
		Redis config2.RedisConfig
	}

	var app = AppConf{
		Redis: config2.RedisConfig{
			Addrs: []string{
				"127.0.0.1:6379",
			},
		},
	}
	return cache.InitRedisByConf(&app.Redis)
}

func TestNewBroadCast(t *testing.T) {
	broadcast := broadcast2.NewBroadCast(&adapter.RedisAdapter{
		Client: getRedis(),
	}, 100)
	broadcast.RegisterHandler("xx", func(msg *broadcast2.Message) {
		fmt.Printf("%v\n", msg)
	})
	broadcast.RegisterHandler("yy", func(msg *broadcast2.Message) {
		fmt.Printf("%v\n", msg)
	})
	go func() {
		for true {
			time.Sleep(time.Millisecond * 245)
			if err := broadcast.Publish(&broadcast2.Message{
				Topic: "xx",
				Data:  []byte("xx"),
			}); err != nil {
				fmt.Printf("%s\n", err.Error())
			}
		}
	}()
	go func() {
		for true {
			time.Sleep(time.Millisecond * 168)
			if err := broadcast.Publish(&broadcast2.Message{
				Topic: "yy",
				Data:  []byte("yy"),
			}); err != nil {
				fmt.Printf("%s\n", err.Error())
			}
		}
	}()
	go func() {
		for true {
			time.Sleep(time.Millisecond * 168)
			if err := broadcast.Publish(&broadcast2.Message{
				Topic: "zz",
				Data:  []byte("zz"),
			}); err != nil {
				fmt.Printf("%s\n", err.Error())
			}
		}
	}()
	if err := broadcast.Subscribe(context.TODO()); err != nil {

	}
}
