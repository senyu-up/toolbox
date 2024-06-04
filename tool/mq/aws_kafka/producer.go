package aws_kafka

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	sigv4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/aws_msk_iam"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/marshaler"
	"github.com/senyu-up/toolbox/tool/trace"
	"log"
	"strings"
	"time"
)

type Producer struct {
	Producer *kafka.Writer

	e         HandleErrorFunc
	s         HandleSucceedFunc
	m         marshaler.Marshaler
	closeChan chan struct{} // 协程关闭通知

	traceOn bool
}

func NewProducer(conf *config.AwsKafkaConfig, opts ...KafkaOption) *Producer {
	var kafkaOpt = &KafkaOpt{}
	var p = &Producer{traceOn: conf.TraceOn, m: marshaler.JsonMarshaler{}, closeChan: make(chan struct{})}
	for _, opt := range opts {
		opt(kafkaOpt)
	}
	if kafkaOpt.m != nil {
		p.m = kafkaOpt.m
	}

	p.Producer = &kafka.Writer{
		Addr:     kafka.TCP(conf.Brokers...),
		Balancer: &kafka.Hash{},
		Async:    kafkaOpt.Async,

		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				if p.e != nil {
					p.e(err)
				} else {
					log.Printf("kafka producer error: %v", err)
				}
			} else {
				// success
				if p.s != nil {
					for _, msg := range messages {
						p.s(msg)
					}
				}
			}
		},
	}
	if conf.SASL.Enable {
		p.Producer.Transport = getSASLTransport(conf)
	}
	if conf.Timeout > 0 {
		p.Producer.WriteTimeout = time.Duration(conf.Timeout) * time.Second
		p.Producer.ReadTimeout = time.Duration(conf.Timeout) * time.Second
	}
	if conf.BackoffMin > 0 {
		p.Producer.WriteBackoffMin = time.Duration(conf.BackoffMin) * time.Millisecond
	}
	if conf.BackoffMax > 0 {
		p.Producer.WriteBackoffMax = time.Duration(conf.BackoffMax) * time.Millisecond
	}
	if conf.MaxAttempts > 0 {
		p.Producer.MaxAttempts = conf.MaxAttempts
	}
	if conf.Balancer != "" {
		p.Producer.Balancer = getKafkaBalancerByStr(conf.Balancer)
	}
	if conf.RequiredAcks > 0 {
		p.Producer.RequiredAcks = kafka.RequiredAcks(conf.RequiredAcks)
	}
	if conf.AllowAutoTopicCreation == true {
		p.Producer.AllowAutoTopicCreation = true
	}
	return p
}

func getSASLTransport(conf *config.AwsKafkaConfig) *kafka.Transport {
	if conf.SASL.Enable {
		switch strings.ToUpper(conf.SASL.Mechanism) {
		case "AWS_MSK_IAM":
			return getAwsCredit(conf.SASL.Region, conf.SASL.AccessId, conf.SASL.SecretKey)
		case "PLAIN":
			mechanism := plain.Mechanism{Username: conf.SASL.UserName, Password: conf.SASL.Password}
			return &kafka.Transport{SASL: mechanism, TLS: &tls.Config{ /*InsecureSkipVerify: true*/ }}
		case "SCRAM-SHA-256":
			mechanism, err := scram.Mechanism(scram.SHA256, conf.SASL.UserName, conf.SASL.Password)
			if err != nil {
				log.Fatalf("scram mechanism 256 error: %v", err)
			}
			return &kafka.Transport{SASL: mechanism, TLS: &tls.Config{}}
		case "SCRAM-SHA-512":
			mechanism, err := scram.Mechanism(scram.SHA512, conf.SASL.UserName, conf.SASL.Password)
			if err != nil {
				log.Fatalf("scram mechanism 512 error: %v", err)
			}
			return &kafka.Transport{SASL: mechanism, TLS: &tls.Config{}}
		default:
			log.Fatalf("not support mechanism: %s", conf.SASL.Mechanism)
		}
	}
	return nil
}

func getAwsCredit(region, accessId, secretKey string) *kafka.Transport {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessId, secretKey, ""),
	})
	if err != nil {
		log.Printf("new aws session error: %v", err)
		return nil
	}

	sharedTransport := &kafka.Transport{
		SASL: &aws_msk_iam.Mechanism{
			Signer: sigv4.NewSigner(sess.Config.Credentials),
			Region: region,
		},
		TLS: &tls.Config{
			//InsecureSkipVerify: true,
		},
	}
	return sharedTransport
}
func (p *Producer) HandleError(e HandleErrorFunc) {
	p.e = e
}

func (p *Producer) HandleSucceed(s HandleSucceedFunc) {
	p.s = s
}

// 等待关闭
func (p *Producer) WaitClose() {
	<-p.closeChan
	return
}

func (p *Producer) Close() error {
	if p.Producer != nil {
		close(p.closeChan)
		return p.Producer.Close()
	}
	return nil
}

func (p *Producer) PushMsgs(ctx context.Context, msgs []kafka.Message) (err error) {
	if p.Producer == nil {
		return ErrKafkaClientIsNil
	}
	if msgs == nil {
		return
	}

	if p.traceOn && 0 < len(msgs) {
		var span opentracing.Span = nil
		ctx, span = kafkaProductSpan(ctx, msgs)
		defer func() {
			if span != nil {
				if err != nil {
					span.SetTag("error", err.Error())
				}
				span.Finish()
			}
		}()
	}
	for i, _ := range msgs {
		msgFill(ctx, &msgs[i])
	}
	return p.Producer.WriteMessages(ctx, msgs...)
}

func (p *Producer) PushObj(ctx context.Context, topic string, event interface{}) (err error) {
	if p.Producer == nil {
		return ErrKafkaClientIsNil
	}

	msg, err := p.m.Marshal(event)
	if err != nil {
		return err
	}
	p.PushMsgs(ctx, []kafka.Message{kafka.Message{
		Topic: topic,
		Value: msg,
	}})
	return nil
}

// msgFilter
//
//	@Description: 消息填充, 检测消息中是否存在链路信息, 没有则进行补充
//	@param ctx  body any true "-"
//	@param msg  body any true "-"
func msgFill(ctx context.Context, msg *kafka.Message) {
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
		msg.Headers = append(msg.Headers, kafka.Header{
			Key:   enum.RequestId,
			Value: []byte(traceId),
		})
	}
	if !hasSpan {
		msg.Headers = append(msg.Headers, kafka.Header{
			Key:   enum.SpanId,
			Value: []byte(spanId),
		})
	}
}

func kafkaProductSpan(ctx context.Context, msgs []kafka.Message) (context.Context, opentracing.Span) {
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
				value = string(msg.Value)
			}
			if msg.Key != nil {
				key = string(msg.Key)
			}
			var tags = map[string]interface{}{"topic": msg.Topic, "key": key, "value": value, "msg_timestamp": msg.Time,
				"offset": msg.Offset, "partition": msg.Partition}
			ctx, reqId, spanId, nextSpanId = trace.ParseOrGenContext(ctx)
			return ctx, trace.NewJaegerSpan(opName, reqId, nextSpanId, spanId, tags, nil)
		}
	}
	return ctx, nil
}
