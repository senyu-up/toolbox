package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/trace"
	"sync"
)

type Kafka struct {
	producer  *Producer
	consumers map[string]*Consumer
	pc        sarama.PartitionerConstructor

	consRwLock sync.RWMutex
	conf       config.KafkaConfig
}

func New(cfg *config.KafkaConfig, opts ...KafkaOption) (*Kafka, error) {
	var (
		err error
		ka  = new(Kafka)
	)

	// producer
	ka.producer, err = NewProducer(cfg, opts...)
	if err != nil {
		return nil, err
	}

	// consumer
	ka.consumers = make(map[string]*Consumer)
	if cfg != nil {
		ka.conf = *cfg
	}
	return ka, nil
	if len(cfg.Consumers) == 0 {
		return ka, nil
	}

	ka.consumers = make(map[string]*Consumer)
	for _, consumerCfg := range cfg.Consumers {
		consumerCfg.Brokers = cfg.Brokers
		consumer, err := NewConsumer(consumerCfg)
		if err != nil {
			return nil, err
		}
		ka.consumers[consumerCfg.Topic] = consumer
	}
	return ka, nil
}

// GetVersion
//
//	@Description: 获取kafka版本
//	@receiver k
//	@return string
func (k *Kafka) GetVersion() string {
	return k.producer.config.Version.String()
}

func (k *Kafka) Consumer(topic string) *Consumer {
	k.consRwLock.RLock()
	defer k.consRwLock.RUnlock()
	return k.consumers[topic]
}

func (k *Kafka) Consumers() map[string]*Consumer {
	k.consRwLock.RLock()
	defer k.consRwLock.RUnlock()
	return k.consumers
}

func (k *Kafka) Producer() *Producer {
	return k.producer
}

// PushAsyncRaw
//
//	@Description: 异步发送一条 复杂、自定义参数多的 sarama消息
//	@receiver k
//	@param msg  body any true "-"
func (k *Kafka) PushAsyncRaw(ctx context.Context, msg *sarama.ProducerMessage) {
	k.producer.PushAsyncRaw(ctx, msg)
}

// PushSyncRawMsgs
//
//	@Description: 同步发送多条消息
//	@receiver k
//	@param ctx  body any true "-"
//	@param msg  body any true "-"
func (k *Kafka) PushSyncRawMsgs(ctx context.Context, msgs []*sarama.ProducerMessage) error {
	return k.producer.PushSyncRawMsgs(ctx, msgs)
}

// PushAsync
//
//	@Description: 异步发送消息，设置了marshal则用它序列化消息体，如果没设置则使用默认的json序列化
//	@receiver k
//	@param topic  body any true "-"
//	@param data  body any true "-"
func (k *Kafka) PushAsync(ctx context.Context, topic string, data interface{}) error {
	return k.producer.PushAsync(ctx, topic, data)
}

// PushSyncRaw
//
//	@Description: 同步发送一条 复杂、自定义参数多的 sarama消息
//	@receiver k
//	@param ctx  body any true "-"
//	@param msg  body any true "-"
//	@return partition
//	@return offset
//	@return err
func (k *Kafka) PushSyncRaw(ctx context.Context, msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return k.producer.PushSyncRaw(ctx, msg)
}

// PushSync
//
//	@Description: 同步发送消息，设置了marshal则用它序列化消息体，如果没设置则使用默认的json序列化
//	@receiver k
//	@param ctx  body any true "-"
//	@param topic  body any true "-"
//	@param data  body any true "-"
//	@return partition
//	@return offset
//	@return err
func (k *Kafka) PushSync(ctx context.Context, topic string, data interface{}) (partition int32, offset int64, err error) {
	return k.producer.PushSync(ctx, topic, data)
}

// RegisterConsumerHandler
//
//	@Description: 注册消费者，handler 入参带有 trace context
//	@receiver k
//	@param topic  body any true "-"
//	@param group  body any true "-"
//	@param handler  body any true "-"
//	@return error
func (k *Kafka) RegisterConsumerHandler(topic string, group string, handler HandleConsumerMsgCtxFunc, errHandler HandleErrorFunc) error {
	if _, ok := k.consumers[topic]; ok { // 重复注册、订阅检查
		return ErrConsumerReRegister
	}

	var c = config.KafkaConsumerConfig{
		Brokers: k.conf.Brokers,
		Version: k.conf.Version,
		Workers: k.conf.Workers,
		Oldest:  k.conf.Oldest,
		SASL:    k.conf.SASL,

		Topic: topic,
		Group: group,
	}
	if newCon, err := newConsumerWithConfig(&c, nil); err != nil {
		return err
	} else {
		// warp handler
		var warpHandler = func(msg *sarama.ConsumerMessage) {
			traceId, spanId := ExtractTraceIdSpanId(msg)
			var ctx = trace.NewContextWithRequestIdAndSpanId(context.Background(), traceId, spanId)
			var err error
			if k.conf.TraceOn {
				var tags = map[string]interface{}{"topic": topic, "group": group, "partition": msg.Partition,
					"offset": msg.Offset, "key": string(msg.Key), "msg_timestamp": msg.Timestamp}
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
		newCon.HandleMsg(warpHandler)
		newCon.e = errHandler
		k.consumers[topic] = newCon
	}
	return nil
}

// ConsumeTopicCtx
//
//	@Description: 消费topic，handler 入参带有 trace context, 该方法废弃，请使用：RegisterConsumerHandler
//	@receiver k
//	@param topic  body any true "-"
//	@param handler  body any true "-"
//	@return error
//
// deprecated
func (k *Kafka) ConsumeTopicCtx(topic string, handler HandleConsumerMsgCtxFunc) error {
	consumer := k.Consumer(topic)
	if consumer == nil {
		return ErrConsumerNotFound
	}
	// warp handler
	var warpHandler = func(msg *sarama.ConsumerMessage) {
		traceId, spanId := ExtractTraceIdSpanId(msg)
		var ctx = trace.NewContextWithRequestIdAndSpanId(context.Background(), traceId, spanId)
		var span = trace.NewJaegerSpan("kafka:topic:"+consumer.topic, traceId,
			trace.NewSpanID(), spanId, nil, nil)
		defer span.Finish()
		var err = handler(ctx, msg)
		if err != nil {
			span.SetTag("error", err.Error())
			logger.Ctx(ctx).SetErr(err).Error("kafka consume topic [" + consumer.topic + "] error")
		}
	}
	consumer.HandleMsg(warpHandler)
	return nil
}

// StartConsume
//
//	@Description:
//	@receiver k
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
	err := k.producer.Close()
	for _, consumer := range k.consumers {
		_ = consumer.Close()
	}
	return err
}

// 从 kafka head 提取 traceId, spanId
// 如果提取不到，就自己生成
func ExtractTraceIdSpanId(msg *sarama.ConsumerMessage) (traceId, spanId string) {
	traceId = trace.NewTraceID()
	spanId = trace.NewSpanID()
	if msg == nil || msg.Headers == nil {
		return
	}
	for _, h := range msg.Headers {
		if string(h.Key) == enum.RequestId {
			traceId = string(h.Value)
		}
		if string(h.Key) == enum.SpanId {
			spanId = string(h.Value)
		}
	}
	return
}
