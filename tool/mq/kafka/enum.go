package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/marshaler"
	"sync"
)

const (
	msgCommit = ""
)

var (
	defaultVersion = sarama.V2_6_0_0
)

var (
	ErrConsumerNotFound   = errors.New("Consumer not found ")
	ErrConsumerTopicExist = errors.New("Consumer topic exist ")
	ErrConsumerReRegister = errors.New("Already register consumer topic - group exist ")
)

type (
	HandleErrorFunc   func(error)
	HandleSucceedFunc func(*sarama.ProducerMessage)

	HandleConsumerMsgFunc    func(message *sarama.ConsumerMessage)
	HandleConsumerMsgCtxFunc func(context.Context, *sarama.ConsumerMessage) error
)

type Producer struct {
	config        *sarama.Config
	asyncProducer sarama.AsyncProducer // 异步生产者
	syncProducer  sarama.SyncProducer  // 同步生产者

	e         HandleErrorFunc
	s         HandleSucceedFunc
	m         marshaler.Marshaler
	closeChan chan struct{} // 协程关闭通知

	hostName string
	ip       string
	traceOn  bool
}

type Consumer struct {
	run bool

	ctx    context.Context
	cancel context.CancelFunc

	group  sarama.ConsumerGroup
	config *sarama.Config
	wg     *sync.WaitGroup

	cfg   *config.KafkaConsumerConfig
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

type Event struct {
	Time     string          `json:"Time,omitempty"`
	Hostname string          `json:"Hostname,omitempty"`
	From     string          `json:"From,omitempty"`
	Type     string          `json:"Type,omitempty"`
	Data     json.RawMessage `json:"Data,omitempty"`
}

func GetDefaultVersion() sarama.KafkaVersion {
	return defaultVersion
}
