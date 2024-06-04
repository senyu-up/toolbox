package nats

import (
	"github.com/nats-io/nats.go"
	"time"
)

type INats interface {
	//PubSub
	Publish(topic string, data []byte) error
	Subscribe(topic string, handler func(data *nats.Msg)) error
	//Queue
	Push(topic string, data []byte) error
	Pop(topic string, groupId string, handler func(data *nats.Msg)) error
	//Request-Response
	Request(topic string, body []byte) (<-chan *nats.Msg, error)
	RequestWithTimeout(topic string, body []byte, response func(data *nats.Msg), timeout time.Duration) error
	Response(topic string, request func(data *nats.Msg) []byte) error
}
