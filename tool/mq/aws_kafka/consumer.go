package aws_kafka

import (
	"context"
	"crypto/tls"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	sigv4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/aws_msk_iam"
	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/trace"
)

type Consumer struct {
	run bool

	ctx    context.Context
	cancel context.CancelFunc

	Reader *kafka.Reader
	wg     *sync.WaitGroup

	cfg   *config.AwsKafkaConfig
	topic string

	// 消息消息错误回调
	e HandleErrorFunc
	// 消费组发生重平衡时调用
	beforeReBalance func()
	// 消息消费停止
	closeCall func()
	// 消息消费函数
	consumerMsgHandler HandleConsumerMsgFunc
	handlerLock        sync.RWMutex
}

func NewConsumer(conf *config.AwsKafkaConfig, opts ...KafkaOption) Consumer {
	var c = Consumer{cfg: conf, wg: &sync.WaitGroup{}, handlerLock: sync.RWMutex{}}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	var kafkaOpt = &KafkaOpt{}
	for _, opt := range opts {
		opt(kafkaOpt)
	}

	var readerConf = kafka.ReaderConfig{
		Brokers: conf.Brokers,
		GroupID: kafkaOpt.GroupId,
		Topic:   kafkaOpt.Topic,
	}
	if conf.SASL.Enable {
		readerConf.Dialer = getSASLDialer(conf)
	}
	// 超时
	if conf.Timeout > 0 {
		readerConf.ReadBatchTimeout = time.Duration(conf.Timeout) * time.Second
	}
	// 轮询新消息的最小等待时间
	if conf.BackoffMin > 0 {
		readerConf.ReadBackoffMin = time.Duration(conf.BackoffMin) * time.Millisecond
	}
	// 轮询新消息的最大等待时间
	if conf.BackoffMax > 0 {
		readerConf.ReadBackoffMax = time.Duration(conf.BackoffMax) * time.Millisecond
	}
	// 最大重试次数
	if conf.MaxAttempts > 0 {
		readerConf.MaxAttempts = conf.MaxAttempts
	}
	// 读取消息的最大字节数
	if conf.MaxBytes > 0 {
		readerConf.MaxBytes = int(conf.MaxBytes)
	}
	// 事务隔离级别
	if conf.IsolationLevel != 0 {
		readerConf.IsolationLevel = kafka.IsolationLevel(conf.IsolationLevel)
	}
	c.Reader = kafka.NewReader(readerConf)

	return c
}

func getSASLDialer(conf *config.AwsKafkaConfig) *kafka.Dialer {
	if conf.SASL.Enable {
		switch strings.ToUpper(conf.SASL.Mechanism) {
		case "AWS_MSK_IAM":
			return getAwsDialer(conf.SASL.Region, conf.SASL.AccessId, conf.SASL.SecretKey)
		case "PLAIN":
			mechanism := plain.Mechanism{Username: conf.SASL.UserName, Password: conf.SASL.Password}
			return &kafka.Dialer{
				Timeout:       time.Duration(conf.Timeout) * time.Second,
				DualStack:     true,
				SASLMechanism: mechanism,
			}
		case "SCRAM-SHA-256":
			mechanism, err := scram.Mechanism(scram.SHA256, conf.SASL.UserName, conf.SASL.Password)
			if err != nil {
				log.Fatalf("scram mechanism 256 error: %v", err)
			}
			return &kafka.Dialer{
				Timeout:       time.Duration(conf.Timeout) * time.Second,
				DualStack:     true,
				SASLMechanism: mechanism,
			}
		case "SCRAM-SHA-512":
			mechanism, err := scram.Mechanism(scram.SHA512, conf.SASL.UserName, conf.SASL.Password)
			if err != nil {
				log.Fatalf("scram mechanism 512 error: %v", err)
			}
			return &kafka.Dialer{
				Timeout:       time.Duration(conf.Timeout) * time.Second,
				DualStack:     true,
				SASLMechanism: mechanism,
			}
		default:
			log.Fatalf("not support mechanism: %s", conf.SASL.Mechanism)
		}
	}
	return nil
}

func getAwsDialer(region, accessId, secretKey string) *kafka.Dialer {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessId, secretKey, ""),
	})
	if err != nil {
		log.Printf("new aws session error: %v", err)
		return nil
	}

	return &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		SASLMechanism: &aws_msk_iam.Mechanism{
			Signer: sigv4.NewSigner(sess.Config.Credentials),
			Region: region,
		},
		TLS: &tls.Config{
			//InsecureSkipVerify: true,
		},
	}
}

func (c *Consumer) Start() {
	c.handlerLock.Lock()
	defer c.handlerLock.Unlock()
	if c.run {
		return
	}
	for i := 0; i < c.cfg.Workers; i++ {
		go c.worker()
	}
	c.run = true
}

func (c *Consumer) Close() error {
	c.cancel()
	c.wg.Wait()
	var err = c.Reader.Close()
	if c.closeCall != nil {
		c.closeCall()
	}
	c.run = false
	return err
}

func (c *Consumer) worker() {
	defer func() {
		log.Printf("consumer worker end\n")
	}()
	c.wg.Add(1)
	defer c.wg.Done()

	for {
		msg, err := c.Reader.ReadMessage(c.ctx)
		if err != nil {
			if c.e != nil {
				c.e(err)
			} else {
				logger.SetErr(err).Error("aws kafka consumer error")
			}
			continue
		} else {
			c.consumerMsgHandler(msg)
		}
		if c.ctx.Err() != nil {
			return
		}
	}
}

func (c *Consumer) HandleMsg(handler HandleConsumerMsgFunc) {
	c.handlerLock.Lock()
	defer c.handlerLock.Unlock()
	c.consumerMsgHandler = handler
}

// 从 kafka head 提取 traceId, spanId
// 如果提取不到，就自己生成
func ExtractTraceIdSpanId(msg kafka.Message) (traceId, spanId string) {
	if msg.Headers == nil {
		traceId = trace.NewTraceID()
		spanId = trace.NewSpanID()
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
	if len(traceId) <= 0 {
		traceId = trace.NewTraceID()
	}
	if len(spanId) <= 0 {
		spanId = trace.NewSpanID()
	}
	return
}
