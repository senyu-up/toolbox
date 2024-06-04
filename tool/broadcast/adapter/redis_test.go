package adapter

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/senyu-up/toolbox/tool/config"
	"testing"
	"time"
)

type mockRedisClient struct {
	pubSub *redis.PubSub
}

func (m *mockRedisClient) Publish(channel string, message interface{}) *redis.IntCmd {
	return redis.NewIntResult(1, nil)
}

func (m *mockRedisClient) Subscribe(channels ...string) *redis.PubSub {
	return m.pubSub
}

func (m *mockRedisClient) Close() error {
	return nil
}

var redClient redis.UniversalClient

func setUp() {
	var redConf = config.RedisConfig{
		Addrs: []string{"127.0.0.1:6379"},
	}
	redClient = cache.InitRedisByConf(&redConf)
}

func TestRedisPublish(t *testing.T) {
	setUp()

	var b = RedisBroadcast{}
	err := b.Init(config.Broadcast{
		Redis: redClient,
		Topic: "test_topic",
	})
	if err != nil {
		t.Errorf("unexpected init error: %v", err)
		return
	}
	b.BroadcastBase.RegisterHandler(func(ctx context.Context, msg []byte) {
		t.Logf("handler1 received msg: %s", string(msg))
	})
	b.BroadcastBase.RegisterHandler(func(ctx context.Context, msg []byte) {
		t.Logf("handler2 received msg: %s", string(msg))
	})

	go b.Subscribe(context.Background())

	for i := 0; i < 10; i++ {
		err = b.Publish(context.Background(), []byte("hello"+fmt.Sprint(i)))
	}

	var b2 = RedisBroadcast{}
	err = b2.Init(config.Broadcast{
		Redis: redClient,
		Topic: "test_topic2",
	})
	if err != nil {
		t.Errorf("unexpected init error: %v", err)
		return
	}

	b2.BroadcastBase.RegisterHandler(func(ctx context.Context, msg []byte) {
		t.Logf("b2 handler1 received msg: %s", string(msg))
	})
	b2.BroadcastBase.RegisterHandler(func(ctx context.Context, msg []byte) {
		t.Logf("b2 handler2 received msg: %s", string(msg))
	})

	go b2.Subscribe(context.Background())

	for i := 0; i < 10; i++ {
		err = b2.Publish(context.Background(), []byte("b2 hello"+fmt.Sprint(i)))
	}

	time.Sleep(time.Second)
}

func TestRedisBroadcast(t *testing.T) {
	//mockPubSub := &redis.PubSub{}
	//mockRedis := &mockRedisClient{pubSub: mockPubSub}
	tests := []struct {
		name        string
		initErr     error
		subscribeFn func(*RedisBroadcast, context.Context) error
		publishFn   func(*RedisBroadcast, context.Context) error
		wantErr     bool
	}{
		{
			name:    "init with nil redis",
			initErr: errors.New("redis is nil"),
			wantErr: true,
		},
		{
			name:    "init with empty topic",
			initErr: errors.New("topic is empty"),
			wantErr: true,
		},
		{
			name: "happy path subscribe",
			subscribeFn: func(rb *RedisBroadcast, ctx context.Context) error {
				go func() {
					time.Sleep(time.Millisecond * 500)
					rb.redis.Publish(rb.topic, "hello")
				}()
				return rb.Subscribe(ctx)
			},
			wantErr: false,
		},
		{
			name: "happy path publish",
			publishFn: func(rb *RedisBroadcast, ctx context.Context) error {
				return rb.Publish(ctx, []byte("hello"))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &RedisBroadcast{
				redis: nil,
				topic: "test_topic",
			}
			if tt.initErr != nil {
				rb.redis = nil
			}
			err := rb.Init(config.Broadcast{
				Redis: nil,
				Topic: "test_topic",
			})
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected init error: %v", err)
				return
			} else if err == nil && tt.wantErr {
				t.Errorf("expected an error but got nil")
				return
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tt.subscribeFn != nil {
				err = tt.subscribeFn(rb, ctx)
				if err != nil && !tt.wantErr {
					t.Errorf("unexpected subscribe error: %v", err)
					return
				} else if err == nil && tt.wantErr {
					t.Errorf("expected an error but got nil")
					return
				}
			}

			if tt.publishFn != nil {
				err = tt.publishFn(rb, ctx)
				if err != nil && !tt.wantErr {
					t.Errorf("unexpected publish error: %v", err)
				} else if err == nil && tt.wantErr {
					t.Errorf("expected an error but got nil")
				}
			}
		})
	}
}
