package main

import (
	"context"
	"github.com/senyu-up/toolbox/tool/broadcast/adapter"
	"github.com/senyu-up/toolbox/tool/config"
	"log"
	"time"
)

func main() {
	var redConf = config.RedisConfig{
		Addrs: []string{"127.0.0.1:6379"},
	}
	var redClient = cache.InitRedisByConf(&redConf)

	var b = adapter.RedisBroadcast{}
	err := b.Init(config.Broadcast{
		Redis: redClient,
		Topic: "test_topic",
	})
	if err != nil {
		log.Printf("unexpected init error: %v", err)
		return
	}

	var ctx = context.TODO()
	for {
		b.Publish(ctx, []byte("hello 2223"))
		log.Printf("published msg")
		time.Sleep(time.Second * 3)
	}
}
