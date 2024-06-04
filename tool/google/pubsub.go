package google

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/encrypt"
	"github.com/senyu-up/toolbox/tool/logger"
	"google.golang.org/api/option"
	"time"
)

type PubSub struct {
	*pubsub.Client
}

func NewPubSub(ctx context.Context, cnf config.GcPubSub) (cli *PubSub, err error) {
	cli = &PubSub{}
	byteData, err := encrypt.Base64Decode([]byte(cnf.CredentialsJson))
	if err != nil {
		return nil, err
	}
	cli.Client, err = pubsub.NewClient(ctx, cnf.ProjectId, option.WithCredentialsJSON(byteData))

	return
}

type Consumer struct {
	inst *pubsub.Subscription
}

func (c *Consumer) Consume(ctx context.Context, handler func(ctx context.Context, msg *pubsub.Message) error) error {
	err := c.inst.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		err := handler(ctx, msg)

		if err != nil {
			logger.Ctx(ctx).SetExtra(logger.E().String("id", msg.ID).String("data", string(msg.Data))).
				Error("pubSub consume", err)
		}

		msg.Ack()
	})

	return err
}

func (p *PubSub) NewConsumer(ctx context.Context, cfg ConsumerConf) (cst *Consumer, err error) {
	topic := p.Client.Topic(cfg.Topic)

	recConf := pubsub.SubscriptionConfig{
		Topic:                         topic,
		PushConfig:                    pubsub.PushConfig{},
		BigQueryConfig:                pubsub.BigQueryConfig{},
		AckDeadline:                   0,
		RetainAckedMessages:           false,
		RetentionDuration:             0,
		ExpirationPolicy:              nil,
		Labels:                        nil,
		EnableMessageOrdering:         false,
		DeadLetterPolicy:              nil,
		Filter:                        "",
		RetryPolicy:                   nil,
		Detached:                      false,
		TopicMessageRetentionDuration: 0,
		EnableExactlyOnceDelivery:     false,
		State:                         0,
	}
	cst = &Consumer{}
	cst.inst, err = p.Client.CreateSubscription(ctx, cfg.Group, recConf)

	return
}

type Producer struct {
	inst *pubsub.Topic
}

func (p *Producer) Send(ctx context.Context, msg *pubsub.Message) error {
	rs := p.inst.Publish(ctx, msg)
	_, err := rs.Get(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) Flush() {
	p.inst.Flush()
}

func (p *Producer) BatchSend(ctx context.Context, msgList []*pubsub.Message, forceFlush ...bool) error {
	for _, msg := range msgList {
		rs := p.inst.Publish(ctx, msg)
		_, err := rs.Get(ctx)
		if err != nil {
			return err
		}
	}

	if len(forceFlush) > 0 && forceFlush[0] {
		p.Flush()
	}

	return nil
}

type ConsumerConf struct {
	Topic string
	Group string
}

type ProducerConf struct {
	Topic             string
	DelayThresholdSec int
	CountThreshold    int
	ByteThreshold     int
	TimeoutSec        int
}

func (p *PubSub) NewProducer(ctx context.Context, cfg ProducerConf) (producer *Producer, err error) {
	t, err := p.Client.CreateTopic(ctx, cfg.Topic)
	if err != nil {
		return nil, err
	}
	t.PublishSettings = pubsub.DefaultPublishSettings
	if cfg.DelayThresholdSec > 0 {
		t.PublishSettings.DelayThreshold = time.Duration(cfg.DelayThresholdSec) * time.Second
	}
	if cfg.CountThreshold > 0 {
		t.PublishSettings.CountThreshold = cfg.CountThreshold
	}

	if cfg.ByteThreshold > 0 {
		t.PublishSettings.ByteThreshold = cfg.ByteThreshold
	}

	if cfg.TimeoutSec > 0 {
		t.PublishSettings.Timeout = time.Duration(cfg.TimeoutSec) * time.Second
	}

	producer = &Producer{
		inst: t,
	}

	return producer, err
}
