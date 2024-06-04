package adapter

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"github.com/senyu-up/toolbox/tool/config"
	"time"
)

type RedisBroadcast struct {
	BroadcastBase
	topic string
	redis redis.UniversalClient
}

func (r *RedisBroadcast) Init(cnf config.Broadcast) error {
	if cnf.Redis == nil {
		return errors.New("redis is nil")
	}
	if len(cnf.Topic) == 0 {
		return errors.New("topic is empty")
	}
	r.redis = cnf.Redis
	r.topic = cnf.Topic

	return nil
}

func (r *RedisBroadcast) Subscribe(ctx context.Context) error {
	//TODO implement me
	sub := r.redis.Subscribe(r.topic)
	for {
		select {
		case msg := <-sub.Channel():
			r.broadcast(ctx, []byte(msg.Payload))
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (r *RedisBroadcast) Publish(ctx context.Context, msg []byte) error {
	//TODO implement me
	rs := r.redis.Publish(r.topic, msg)
	if rs.Val() == 0 {
		return nil
	} else {
		return rs.Err()
	}
}
