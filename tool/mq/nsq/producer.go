package nsq

import (
	"context"
	"github.com/nsqio/go-nsq"
)

func InitProducer(ctx context.Context, addr string, topic string, pusher chan []byte) error {
	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(addr, cfg)
	if err != nil {
		return err
	}
	var data []byte
	var open bool
	for {
		select {
		case data, open = <-pusher:
			if !open {
				break
			}
			err = producer.Publish(topic, data)
			if err != nil {
				logger.ERR("NSQ push msg to ", topic, " error:", err)
			}
		case <-ctx.Done():
			break
		}
	}
}
