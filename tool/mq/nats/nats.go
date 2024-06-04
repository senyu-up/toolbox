package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"strings"
	"time"
)

type NatsHandle struct {
	client *nats.Conn
}

func InitNats(hosts []string) (*NatsHandle, error) {
	client, err := nats.Connect(strings.Join(hosts, ","))
	if err != nil {
		return nil, err
	}
	return &NatsHandle{client: client}, nil
}

func (p *NatsHandle) Publish(topic string, data []byte) error {
	return p.client.Publish(topic, data)
}

func (p *NatsHandle) Subscribe(topic string, handler func(data *nats.Msg)) error {
	_, err := p.client.Subscribe(topic, handler)
	if err != nil {
		return err
	}
	return nil
}

func (p *NatsHandle) Response(topic string, handler func(request *nats.Msg) []byte) error {
	var result []byte
	_, err := p.client.Subscribe(topic, func(msg *nats.Msg) {
		result = handler(msg)
		_ = msg.Respond(result)
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *NatsHandle) Push(topic string, data []byte) error {
	return p.Publish(topic, data)
}

func (p *NatsHandle) Pop(topic string, groupId string, handler func(data *nats.Msg)) error {
	_, err := p.client.QueueSubscribe(topic, groupId, handler)
	if err != nil {
		return err
	}
	return nil
}

func (p *NatsHandle) Request(topic string, body []byte) (<-chan *nats.Msg, error) {
	reqTopic := fmt.Sprintf("%s.request", topic)
	rspTopic := fmt.Sprintf("%s.response", topic)
	ch := make(chan *nats.Msg, 1)
	err := p.client.PublishRequest(reqTopic, rspTopic, body)
	if err != nil {
		return nil, err
	}
	_, err = p.client.Subscribe(rspTopic, func(msg *nats.Msg) {
		ch <- msg
	})
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (p *NatsHandle) RequestWithTimeout(topic string, body []byte, response func(data *nats.Msg), timeout time.Duration) error {
	msg, err := p.client.Request(topic, body, timeout)
	if err != nil {
		return err
	}
	response(msg)
	return nil
}
