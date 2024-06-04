package aws_kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/trace"
	"sync"
)

type (
	HandleErrorFunc   func(error)
	HandleSucceedFunc func(message kafka.Message)

	HandleConsumerMsgFunc    func(message kafka.Message)
	HandleConsumerMsgCtxFunc func(context.Context, kafka.Message) error
)

var (
	ErrKafkaClientIsNil   = fmt.Errorf("Kafka client is nil ")
	ErrKafkaConfigIsNil   = fmt.Errorf("Kafka config is nil ")
	ErrTopicNotFound      = fmt.Errorf("Topic not found ")
	ErrConsumerReRegister = fmt.Errorf("Already register consumer topic - group exist ")
)

type Kafka struct {
	consRwLock sync.RWMutex
	conf       *config.AwsKafkaConfig
	opts       []KafkaOption
	Conn       *kafka.Conn // admin conn

	syncProducer  *Producer // 同步发送
	asyncProducer *Producer // 异步发送

	consumers map[string]Consumer // 消费者们， topicName-> reader
}

func New(ctx context.Context, cfg *config.AwsKafkaConfig, opts ...KafkaOption) (k *Kafka, err error) {
	k = &Kafka{conf: cfg, opts: opts}
	if cfg.SASL.Enable {
		var d = getAwsDialer(cfg.SASL.Region, cfg.SASL.AccessId, cfg.SASL.SecretKey)
		k.Conn, err = d.DialContext(ctx, "tcp", cfg.Brokers[0])
	} else {
		k.Conn, err = kafka.DialContext(ctx, "tcp", cfg.Brokers[0])
	}

	// producer
	k.syncProducer = NewProducer(cfg, append(opts, KafkaOptWithAsync(false))...)
	// async producer
	k.asyncProducer = NewProducer(cfg, append(opts, KafkaOptWithAsync(true))...)

	// consumer
	k.consumers = make(map[string]Consumer, 0)
	return k, err
}

// GetVersion
//
//	@Description: 获取kafka版本
//	@receiver k
//	@return string
func (k *Kafka) GetVersion() string {
	return "unknown"
}

func (k *Kafka) Consumer(topic string) Consumer {
	k.consRwLock.RLock()
	defer k.consRwLock.RUnlock()
	return k.consumers[topic]
}

func (k *Kafka) Consumers() map[string]Consumer {
	k.consRwLock.RLock()
	defer k.consRwLock.RUnlock()
	return k.consumers
}

func (k *Kafka) SyncProducer() *Producer {
	return k.syncProducer
}

func (k *Kafka) AsyncProducer() *Producer {
	return k.asyncProducer
}

func (k *Kafka) PushSyncRawMsgs(ctx context.Context, msgs []kafka.Message) error {
	return k.syncProducer.PushMsgs(ctx, msgs)
}

func (k *Kafka) PushAsyncRawMsgs(ctx context.Context, msgs []kafka.Message) error {
	return k.asyncProducer.PushMsgs(ctx, msgs)
}

func (k *Kafka) PushSyncMsgs(ctx context.Context, topic string, data interface{}) error {
	return k.syncProducer.PushObj(ctx, topic, data)
}

func (k *Kafka) PushAsyncMsgs(ctx context.Context, topic string, data interface{}) error {
	return k.asyncProducer.PushObj(ctx, topic, data)
}

func (k *Kafka) RegisterConsumerHandler(topic string, group string, handler HandleConsumerMsgCtxFunc, errHandler HandleErrorFunc) error {
	if k.conf == nil {
		return ErrKafkaConfigIsNil
	}
	if _, ok := k.consumers[topic]; ok { // 重复注册、订阅检查
		return ErrConsumerReRegister
	}

	var consumer = NewConsumer(k.conf, append(k.opts, KafkaOptWithGroupId(group), KafkaOptWithTopic(topic))...)
	consumer.e = errHandler
	// warp handler
	var warpHandler = func(msg kafka.Message) {
		traceId, spanId := ExtractTraceIdSpanId(msg)
		var ctx = trace.NewContextWithRequestIdAndSpanId(context.Background(), traceId, spanId)
		var err error
		if k.conf.TraceOn {
			var tags = map[string]interface{}{"topic": topic, "group": group, "partition": msg.Partition,
				"offset": msg.Offset, "key": string(msg.Key), "msg_time": msg.Time}
			var span = trace.NewJaegerSpan("kafka:topic:"+topic+":group:"+group, traceId,
				trace.NewSpanID(), spanId, tags, nil)
			defer func() {
				if err != nil {
					span.SetTag("error", err.Error())
				}
				span.Finish()
			}()
		}
		err = handler(ctx, msg)
		if err != nil {
			logger.Ctx(ctx).SetErr(err).Error("kafka consume topic [" + topic + "] error")
		}
	}
	consumer.HandleMsg(warpHandler)
	k.consumers[topic] = consumer
	return nil
}

func (k *Kafka) StartConsume() {
	k.consRwLock.RLock()
	defer k.consRwLock.RUnlock()
	for _, consumer := range k.consumers {
		go consumer.Start()
	}
}

// Close
//
//	@Description: 关闭 生产者、消费者
//	@receiver k
//	@return error
func (k *Kafka) Close() error {
	err := k.Conn.Close() // admin conn close
	if err != nil {
		logger.Warn("kafka conn close err: ", err)
	}
	err = k.asyncProducer.Close()
	if err != nil {
		logger.Warn("kafka async producer close err: ", err)
	}
	err = k.syncProducer.Close()
	if err != nil {
		logger.Warn("kafka sync producer close err: ", err)
	}
	for topic, consumer := range k.consumers {
		err = consumer.Close()
		if err != nil {
			logger.Warn("kafka consumer topic ["+topic+"] close err: ", err)
		}
	}
	return err
}
