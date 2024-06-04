package adapter

import (
	"context"
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"github.com/senyu-up/toolbox/tool/broadcast"
)

/*
RedisAdapter
预定义的消息中间件,基于 redis Pub-Sub实现发布订阅
*/

type RedisAdapter struct {
	Client redis.UniversalClient
}

func (r *RedisAdapter) Subscribe(ctx context.Context, topic broadcast.Topic, c chan<- *broadcast.Message) error {
	sub := r.Client.Subscribe(string(topic))
	for {
		select {
		case msg := <-sub.Channel():
			data := &broadcast.Message{}
			_ = jsoniter.UnmarshalFromString(msg.Payload, &data)
			c <- data
		case <-ctx.Done():
			return nil
		}
	}
}

func (r *RedisAdapter) Publish(msg *broadcast.Message) error {
	data, _ := jsoniter.MarshalToString(msg)
	_, err := r.Client.Publish(string(broadcast.Default), data).Result()
	return err
}
