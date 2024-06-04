package nsq

import (
	"github.com/nsqio/go-nsq"
	"time"
)

type QueueHandler struct {
	track   chan *Elem
	topic   string
	channel string
}

type Elem struct {
	Topic   string
	Channel string
	Payload []byte
}

func (p *QueueHandler) Take() <-chan *Elem {
	return p.track
}

func (p *QueueHandler) HandleMessage(m *nsq.Message) error {
	//todo pool
	p.track <- &Elem{Topic: p.topic, Channel: p.channel, Payload: m.Body}
	return nil
}

func InitConsumer(topic string, channel string, nsqd string, handler func(c <-chan *Elem)) error {
	config := nsq.NewConfig()
	config.HeartbeatInterval = time.Second * 3
	config.LookupdPollInterval = time.Second * 3
	config.MaxInFlight = 10
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return err
	}
	consumer.SetLoggerLevel(nsq.LogLevelError)
	qh := &QueueHandler{
		track:   make(chan *Elem),
		topic:   topic,
		channel: channel,
	}
	consumer.AddHandler(qh)
	err = consumer.ConnectToNSQD(nsqd)
	if err != nil {
		return err
	}
	handler(qh.Take())
	return nil
}

func InitConsumerLookup(topic string, channel string, lookUpd string, handler func(c <-chan *Elem)) error {
	config := nsq.NewConfig()
	config.HeartbeatInterval = time.Second * 3
	config.LookupdPollInterval = time.Second * 3
	config.MaxInFlight = 10
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return err
	}
	consumer.SetLoggerLevel(nsq.LogLevelError)
	qh := &QueueHandler{
		track:   make(chan *Elem),
		topic:   topic,
		channel: channel,
	}
	consumer.AddHandler(qh)
	err = consumer.ConnectToNSQLookupd(lookUpd)
	if err != nil {
		return err
	}
	handler(qh.Take())
	return nil
}
