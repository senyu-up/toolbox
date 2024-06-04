package broadcast

import (
	"context"
	"github.com/senyu-up/toolbox/tool/runtime"
)

type Topic string
type Handler func(msg *Message)

const (
	Default Topic = "__center__broadcast__"
)

type Broadcast struct {
	adapter  IAdapterInterface
	handlers map[Topic]Handler
	cMsg     chan *Message
}

type Message struct {
	Topic Topic
	Data  []byte
}

// NewBroadCast 获取一个广播对象
func NewBroadCast(adapter IAdapterInterface, bufferLen int) *Broadcast {
	return &Broadcast{
		adapter:  adapter,
		handlers: map[Topic]Handler{},
		cMsg:     make(chan *Message, bufferLen),
	}
}

// Subscribe 开启订阅
func (b *Broadcast) Subscribe(ctx context.Context) error {
	b.handleTask(ctx)
	if err := b.adapter.Subscribe(ctx, Default, b.cMsg); err != nil {
		return err
	}
	return nil
}

// RegisterHandler 注册触发器
func (b *Broadcast) RegisterHandler(topic Topic, handler Handler) {
	b.handlers[topic] = handler
}

// Publish 推送消息
func (b *Broadcast) Publish(msg *Message) error {
	return b.adapter.Publish(msg)
}

func (b *Broadcast) handleTask(ctx context.Context) {
	runtime.GOSafe(ctx, string(Default), func() {
		for true {
			select {
			case msg := <-b.cMsg:
				if handler, ok := b.handlers[msg.Topic]; ok {
					handler(msg)
				} else {
					//TODO log
				}
			case <-ctx.Done():
				return
			}
		}
	})
}
