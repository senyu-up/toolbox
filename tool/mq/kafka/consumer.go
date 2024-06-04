package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/config"
	"log"
	"sync"
	"time"
)

func NewConsumer(cfg *config.KafkaConsumerConfig) (*Consumer, error) {
	return newConsumerWithConfig(cfg, nil)
}

func NewConsumerWithHandler(cfg *config.KafkaConsumerConfig, handler HandleConsumerMsgFunc) (*Consumer, error) {
	if consumer, err := newConsumerWithConfig(cfg, nil); err != nil {
		return consumer, err
	} else {
		consumer.HandleMsg(handler)
		go consumer.Start()
		return consumer, err
	}
}

func newConsumerWithConfig(cfg1 *config.KafkaConsumerConfig, cfg *sarama.Config) (*Consumer, error) {
	if cfg == nil {
		cfg = sarama.NewConfig()
		cfg.Consumer.Return.Errors = true
		cfg.Version = GetDefaultVersion()
		if cfg1.Oldest {
			cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
		} else {
			cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
		}
		if cfg1.SASL.Enable {
			cfg.Net.SASL.Enable = true
			cfg.Net.SASL.User = cfg1.SASL.User
			cfg.Net.SASL.Password = cfg1.SASL.Password
		}
	}
	if cfg1.Workers == 0 {
		cfg1.Workers = 5 // 如果为0要设置为默认值5
	}
	// kafka 版本指定
	if 1 > len(cfg1.Version) {
		if ver, err := sarama.ParseKafkaVersion(cfg1.Version); err != nil {
			// 制定版本且有效，则设置
			cfg.Version = ver
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	// Start with a client
	client, err := sarama.NewClient(cfg1.Brokers, cfg)
	if err != nil {
		return nil, err
	}
	// Start a new consumer group
	consumerGroup, err := sarama.NewConsumerGroupFromClient(cfg1.Group, client)
	if err != nil {
		return nil, err
	}
	c := new(Consumer)
	c.cfg = cfg1
	c.config = cfg
	c.group = consumerGroup
	c.wg = &sync.WaitGroup{}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.handlerLock = sync.RWMutex{}

	return c, nil
}

type consumerGroupHandler struct {
	consumer *Consumer
}

func (h consumerGroupHandler) getTopic() string {
	return h.consumer.GetTopic()
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			// 如果定义了需要原生接受的方法
			if h.consumer.consumerMsgHandler != nil {
				h.consumer.consumerMsgHandler(msg)
			}
			sess.MarkMessage(msg, msgCommit)
		}
	}
}

func (c *Consumer) HandleError(e HandleErrorFunc) {
	c.e = e
}

func (c *Consumer) HandleReBalance(reBalance func()) {
	c.beforeReBalance = reBalance
}

func (c *Consumer) HandleClose(closeCall func()) {
	c.closeCall = closeCall
}

func (c *Consumer) GetTopic() string {
	return c.cfg.Topic
}

func (c *Consumer) HandleMsg(handler HandleConsumerMsgFunc) {
	c.handlerLock.Lock()
	defer c.handlerLock.Unlock()
	c.consumerMsgHandler = handler
}

// Start
//
//	@Description: 开启消费组, 该函数阻塞运行
//	@receiver c
func (c *Consumer) Start() {
	defer func() {
		fmt.Print("kafka consumer end\n")
	}()
	c.handlerLock.Lock()
	defer c.handlerLock.Unlock()
	if c.run {
		return
	}
	for i := 0; i < c.cfg.Workers; i++ {
		go c.worker()
	}
	c.run = true
	c.Handle()
}

// Close
//
//	@Description: 关闭所有消费组
//	@receiver c
//	@return error
func (c *Consumer) Close() error {
	c.cancel()
	c.wg.Wait()
	var err = c.group.Close()
	if c.closeCall != nil {
		c.closeCall()
	}
	return err
}

func (c *Consumer) worker() {
	defer func() {
		fmt.Print("consumer worker end\n")
	}()
	c.wg.Add(1)
	defer c.wg.Done()
	topics := []string{c.cfg.Topic}
	handler := consumerGroupHandler{c}
	for {
		if err := c.group.Consume(c.ctx, topics, handler); err != nil {
			if c.e != nil {
				c.e(err)
			}
		} else {
			if c.beforeReBalance != nil {
				c.beforeReBalance()
			}
			log.Printf("Kafka consumer rebalanced topic: " + c.cfg.Topic)
			// 休眠一下再重新建立连接
			time.Sleep(time.Second * 2)
		}

		if c.ctx.Err() != nil {
			return
		}
	}
}

func (c *Consumer) Handle() {
	for {
		select {
		case err, ok := <-c.group.Errors():
			{
				if !ok {
					return
				}
				if c.e != nil {
					c.e(err)
				}
			}
		}
	}
}

// ParserTrace
//
//	@Description: 从kafka消息头中提取traceId和spanId
//	@param headers  body any true "-"
//	@return traceId
//	@return spanId
func ParserTrace(headers []*sarama.RecordHeader) (traceId string, spanId string) {
	if nil == headers {
		return
	}
	for _, v := range headers {
		if string(v.Key) == enum.RequestId {
			traceId = string(v.Value)
		}
		if string(v.Key) == enum.SpanId {
			spanId = string(v.Value)
		}
	}
	return traceId, spanId
}
