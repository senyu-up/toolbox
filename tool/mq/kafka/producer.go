package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/marshaler"
	"github.com/senyu-up/toolbox/tool/trace"
	"time"
)

func NewProducer(cnf *config.KafkaConfig, opts ...KafkaOption) (*Producer, error) {
	var (
		cfg   = sarama.NewConfig()
		kaOpt = &KafkaOpt{m: marshaler.JsonMarshaler{}, pc: sarama.NewHashPartitioner}
	)
	// 应用options
	for _, opt := range opts {
		opt(kaOpt)
	}

	cfg.Net.KeepAlive = 60 * time.Second
	cfg.Version = defaultVersion

	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Flush.Frequency = time.Second
	cfg.Producer.Flush.MaxMessages = 10
	cfg.Producer.Partitioner = kaOpt.pc
	// 最大重试次数
	cfg.Producer.Retry.Max = 3

	// kafka 版本指定
	if 1 > len(cnf.Version) {
		if ver, err := sarama.ParseKafkaVersion(cnf.Version); err != nil {
			// 制定版本且有效，则设置
			cfg.Version = ver
		}
	}

	// 重试
	if cnf.ProducerMaxRetryTimes > 0 {
		cfg.Producer.Retry.Max = cnf.ProducerMaxRetryTimes
	} else if cnf.ProducerMaxRetryTimes == -1 {
		cfg.Producer.Retry.Max = 0
	}
	// 超时
	cfg.Producer.Timeout = time.Duration(cnf.Timeout) * time.Second
	if cnf.ProducerTimeout > 0 {
		cfg.Producer.Timeout = time.Duration(cnf.ProducerTimeout) * time.Second
	}

	// meta
	if cnf.MetadataTimeout > 0 {
		cfg.Metadata.Timeout = time.Duration(cnf.MetadataTimeout) * time.Second
	} else {
		cfg.Metadata.Timeout = time.Second * 30
	}
	if cnf.MetadataRefreshIntervalSecond > 0 {
		cfg.Metadata.RefreshFrequency = time.Duration(cnf.MetadataRefreshIntervalSecond) * time.Second
	}
	if cnf.MetadataMaxRetryTimes == -1 {
		cfg.Metadata.Retry.Max = 0
	} else if cnf.MetadataMaxRetryTimes > 0 {
		cfg.Metadata.Retry.Max = cnf.MetadataMaxRetryTimes
	}
	cfg.Metadata.Full = cnf.SyncFullMetadata

	// 失败重试间隔
	cfg.Producer.Retry.Backoff = time.Second
	if cnf.Level == 2 {
		cfg.Producer.RequiredAcks = sarama.WaitForAll
	} else if cnf.Level == 1 {
		cfg.Producer.RequiredAcks = sarama.WaitForLocal
	} else {
		cfg.Producer.RequiredAcks = sarama.NoResponse
	}

	// sasl
	cfg.Net.SASL.User = cnf.User
	cfg.Net.SASL.Password = cnf.Password

	//Dial网络时间配置
	if cnf.DialTimeout > 0 {
		cfg.Net.DialTimeout = time.Duration(cnf.DialTimeout) * time.Second
	} else {
		cfg.Net.DialTimeout = time.Second * 30
	}
	if cnf.ReadTimeout > 0 {
		cfg.Net.ReadTimeout = time.Duration(cnf.ReadTimeout) * time.Second
	} else {
		cfg.Net.ReadTimeout = time.Second * 30
	}
	if cnf.WriteTimeout > 0 {
		cfg.Net.WriteTimeout = time.Duration(cnf.WriteTimeout) * time.Second
	} else {
		cfg.Net.WriteTimeout = time.Second * 30
	}

	return newProducerWithCfg(cnf.Brokers, cfg, cnf.TraceOn, kaOpt.m)
}

func newProducerWithCfg(brokers []string, cfg *sarama.Config, traceOn bool, m marshaler.Marshaler) (*Producer, error) {
	producer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	sp, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	p := new(Producer)
	p.asyncProducer = producer
	p.syncProducer = sp
	p.config = cfg
	p.m = m
	p.traceOn = traceOn
	p.closeChan = make(chan struct{}, 0)
	go p.handle() // 异步发送消息 callback
	return p, nil
}

// handle
//
//	@Description: 异步消息发送成功、失败后的消息回调订阅
//	@receiver p
func (p *Producer) handle() {
	defer func() {
		fmt.Print("kafka producer handler end\n")
	}()
	for {
		select {
		case err, ok := <-p.asyncProducer.Errors():
			{
				if ok {
					if p.e != nil {
						p.e(err)
					}
				}
			}
		case msg, ok := <-p.asyncProducer.Successes():
			{
				if ok {
					if p.s != nil {
						p.s(msg)
					}
				}
			}
		case <-p.closeChan:
			return
		}
	}
}

// Close
//
//	@Description: 关闭 producer
//	@receiver p
//	@return error
func (p *Producer) Close() error {
	var sErr = p.syncProducer.Close()
	var aErr = p.asyncProducer.Close()
	p.closeChan <- struct{}{}
	if sErr != nil {
		return sErr
	}
	return aErr
}

func (p *Producer) HandleError(e HandleErrorFunc) {
	p.e = e
}

func (p *Producer) HandleSucceed(s HandleSucceedFunc) {
	p.s = s
}

// PushSyncRaw
//
//	@Description: 同步发送 sarama 消息
//	@receiver p
//	@param ctx  body any true "-"
//	@param msg  body any true "-"
//	@return partition
//	@return offset
//	@return err
func (p *Producer) PushSyncRaw(ctx context.Context, msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	if msg == nil {
		return
	}
	if p.traceOn && msg != nil {
		var span opentracing.Span = nil
		ctx, span = kafkaProductSpan(ctx, []*sarama.ProducerMessage{msg})
		defer func() {
			if span != nil {
				span.SetTag("kafka_partition", partition)
				span.SetTag("kafka_offset", offset)
				if err != nil {
					span.SetTag("error", err.Error())
				}
				span.Finish()
			}
		}()
	}
	msgFill(ctx, msg)
	partition, offset, err = p.syncProducer.SendMessage(msg)
	return partition, offset, err
}

// PushSyncRawMsgs
//
//	@Description: 批量同步发送 sarama 消息列表
//	@receiver p
//	@param ctx  body any true "-"
//	@param msgs  body any true "-"
//	@return err
func (p *Producer) PushSyncRawMsgs(ctx context.Context, msgs []*sarama.ProducerMessage) (err error) {
	if msgs == nil {
		return
	}
	if p.traceOn && msgs != nil {
		var span opentracing.Span = nil
		ctx, span = kafkaProductSpan(ctx, msgs)
		if span != nil {
			defer span.Finish()
		}
	}
	for _, msg := range msgs {
		msgFill(ctx, msg)
	}
	return p.syncProducer.SendMessages(msgs)
}

// PushSync
//
//	@Description: 同步发送 消息 到指定 topic
//	@receiver p
//	@param ctx  body any true "-"
//	@param topic  body any true "-"
//	@param event  body any true "-"
//	@return error
func (p *Producer) PushSync(ctx context.Context, topic string, event interface{}) (partition int32, offset int64, err error) {
	msg, err := p.m.Marshal(event)
	if err != nil {
		return 0, 0, err
	}
	return p.PushSyncRaw(ctx, &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
	})
}

// PushAsyncRaw
//
//	@Description: 异步发送 sarama 消息
//	@receiver p
//	@param ctx  body any true "-"
//	@param msg  body any true "-"
func (p *Producer) PushAsyncRaw(ctx context.Context, msg *sarama.ProducerMessage) {
	if msg == nil {
		return
	}
	if p.traceOn {
		var span opentracing.Span = nil
		ctx, span = kafkaProductSpan(ctx, []*sarama.ProducerMessage{msg})
		if span != nil {
			defer span.Finish()
		}
	}
	msgFill(ctx, msg)
	p.asyncProducer.Input() <- msg
}

func kafkaProductSpan(ctx context.Context, msgs []*sarama.ProducerMessage) (context.Context, opentracing.Span) {
	if msgs != nil {
		var msgLen = len(msgs)
		if msgLen > 0 {
			var msg = msgs[0] // 批量发送，值记录第一条
			var reqId, spanId, nextSpanId, value, key string
			var opName = fmt.Sprintf("kafka:topic:%s", msg.Topic)
			if msgLen > 1 {
				opName = fmt.Sprintf("kafka:topic:%s:len:%d", msg.Topic, msgLen)
			}
			if msg.Value != nil {
				if b, err := msg.Value.Encode(); err == nil {
					value = string(b)
				}
			}
			if msg.Key != nil {
				if b, err := msg.Key.Encode(); err == nil {
					key = string(b)
				}
			}
			var tags = map[string]interface{}{"topic": msg.Topic, "key": key, "value": value, "msg_timestamp": msg.Timestamp,
				"offset": msg.Offset, "partition": msg.Partition}
			ctx, reqId, spanId, nextSpanId = trace.ParseOrGenContext(ctx)
			return ctx, trace.NewJaegerSpan(opName, reqId, nextSpanId, spanId, tags, nil)
		}
	}
	return ctx, nil
}

// PushAsync
//
//	@Description: 异步发送 消息 到指定 topic
//	@receiver p
//	@param ctx  body any true "-"
//	@param msg  body any true "-"
func (p *Producer) PushAsync(ctx context.Context, topic string, event interface{}) error {
	msg, err := p.m.Marshal(event)
	if err != nil {
		return err
	}
	p.PushAsyncRaw(ctx, &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
	})
	return nil
}

// msgFilter
//
//	@Description: 消息填充, 检测消息中是否存在链路信息, 没有则进行补充
//	@param ctx  body any true "-"
//	@param msg  body any true "-"
func msgFill(ctx context.Context, msg *sarama.ProducerMessage) {
	traceId, spanId := trace.ParseCurrentContext(ctx)
	// 检测header中是否存在链路信息, 没有则进行补充
	var hasTrace bool
	var hasSpan bool

	for i, _ := range msg.Headers {
		if string(msg.Headers[i].Key) == enum.RequestId {
			hasTrace = true
		} else if string(msg.Headers[i].Key) == enum.SpanId {
			hasSpan = true
		}
	}
	if !hasTrace {
		msg.Headers = append(msg.Headers, sarama.RecordHeader{
			Key:   []byte(enum.RequestId),
			Value: []byte(traceId),
		})
	}
	if !hasSpan {
		msg.Headers = append(msg.Headers, sarama.RecordHeader{
			Key:   []byte(enum.SpanId),
			Value: []byte(spanId),
		})
	}
}
