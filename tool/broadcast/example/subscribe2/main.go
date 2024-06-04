package main

import (
	"context"
	"github.com/senyu-up/toolbox/tool/broadcast/adapter"
	"github.com/senyu-up/toolbox/tool/config"
	"log"
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

	b.BroadcastBase.RegisterHandler(func(ctx context.Context, msg []byte) {
		log.Printf("handler2 received msg: %s", string(msg))
	})
	b.BroadcastBase.RegisterHandler(func(ctx context.Context, msg []byte) {
		log.Printf("handler2-2 received msg: %s", string(msg))
	})

	b.Subscribe(ctx)
}
