package queue

import (
	"encoding/base64"
	"errors"
	"github.com/go-redis/redis"
	"github.com/senyu-up/toolbox/tool/logger"
)

var (
	TOPIC_EXISTS = errors.New("topic already exists")
)

// IQueueTopic 处理的topic元素
type IQueueTopic struct {
	QueueScheme string
	HashValue   int
}

type QueueMap struct {
	client    *redis.Client
	QueueName string
}

// IQueue 消费处理接口
type IQueue interface {
	Push(hash string, data []byte) error

	Pop() ([]byte, error)

	Length() int64
}

func (q *QueueMap) Push(hash string, data []byte) error {

	return errors.New("Push Unavailable")

}

func (q *QueueMap) Pop() ([]byte, error) {
	result, err := q.client.LPop(q.QueueName).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if err == redis.Nil {
		return []byte{}, nil
	}
	b, err := base64.StdEncoding.DecodeString(result)
	if err != nil {
		logger.SetErr(err).Error("%v", err)
		return nil, err
	}
	return b, nil
}

func (q *QueueMap) Length() int64 {

	length, _ := q.client.LLen(q.QueueName).Result()
	return length
}
