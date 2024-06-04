package queue

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/senyu-up/toolbox/tool/str"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

/*
	因为redis是单线程模型,保证了所有的数据仅会被处理一次，避免相同的数据被多次处理

*/

const (
	QueueKey = "_queue:%s:%d"
)

type RedisQ struct {
	client        *redis.Client
	topics        *sync.Map
	checkInterval *time.Ticker //周期性检查
	channel       IQueueTopic
	w             *sync.WaitGroup
	process       Process
	pushChan      chan string
	isConsumer    bool
	coroutine     chan struct{}
	ctx           context.Context
}

// Process 处理方法
type Process func(qname string, data []byte) error

/*
	ctx:上下文
	client:redis实例化后的一个客户端
	interval:当没有收到消费通知的时候，用于周期性检查消费队列数据是否存在
	channel:定义的一个父的消费主题,通过hash值生成对应的子消费主题
	work:消费方的处理方法的实现

*/

func NewRedisMQClient(ctx context.Context, client *redis.Client, interval time.Duration, channel IQueueTopic, work Process) *RedisQ {
	check := time.NewTicker(interval)
	mqclient := RedisQ{
		ctx:           ctx,
		client:        client,
		checkInterval: check,
		topics:        new(sync.Map),
		channel:       channel,
		process:       work,
		coroutine:     make(chan struct{}, channel.HashValue), //控制最大协程开启的数量
	}
	for i := 0; i < channel.HashValue; i++ {
		queueKey := fmt.Sprintf(QueueKey, channel.QueueScheme, i)
		mqclient.topics.Store(queueKey, nil)
	}
	//是否为消费者模式
	if work != nil {
		mqclient.isConsumer = true
	}

	go mqclient.notify(channel.QueueScheme)

	//队列监控

	return &mqclient
}

func (rq *RedisQ) Push(hash string, data []byte) error {
	//哈希计算
	hashid := str.HashID(hash, rq.channel.HashValue)
	queueKey := fmt.Sprintf(QueueKey, rq.channel.QueueScheme, hashid)
	if _, ok := rq.topics.Load(queueKey); !ok {
		return errors.New("not found queue key")
	}
	s := base64.StdEncoding.EncodeToString(data)
	err := rq.client.RPush(queueKey, s).Err()
	if err != nil {
		return err
	}
	//跨进程处理
	rq.client.Publish(rq.channel.QueueScheme, queueKey)
	return nil
}

// 数据解码处理
func (rq *RedisQ) msgProcess(qname string) {

	for {
		//避免当消息过多时，进程不能迟迟结束
		select {
		case <-rq.ctx.Done():
			return
		default:
			result, err := rq.client.LPop(qname).Result()
			//队列为空
			if err == redis.Nil {
				return
			}
			if err != nil && err != redis.Nil {
				log.Error(err)
				continue
			}
			b, err := base64.StdEncoding.DecodeString(result)
			if err != nil {
				log.Error(err)
				continue
			}
			//控制协程执行数量
			rq.coroutine <- struct{}{}
			go func() {
				defer func() {
					if err := recover(); err != nil {
						log.Error(err)
					}
					<-rq.coroutine
				}()
				err = rq.process(qname, b)
				if err != nil {
					log.Errorf("topic:%s,data:%s,err:%v", qname, b, err)
				}
			}()

		}

	}
}

// 通知要消费的队列数据
func (rq *RedisQ) notify(channel string) {
	pubsub := rq.client.PSubscribe(channel)
	for {
		select {
		//接受消息订阅
		case msg := <-pubsub.Channel():
			if _, ok := rq.topics.Load(msg.Payload); ok && rq.isConsumer {
				rq.msgProcess(msg.Payload)
			}
		//定时检查队列是否存在数据
		case <-rq.checkInterval.C:
			if rq.isConsumer {
				rq.topics.Range(func(key, value interface{}) bool {
					rq.msgProcess(key.(string))
					return true
				})
			}
		//结束
		case <-rq.ctx.Done():
			return
		}

	}
}

func (rq *RedisQ) QueueClient() (q []*QueueMap) {
	rq.topics.Range(func(key, value interface{}) bool {
		q = append(q, &QueueMap{
			client:    rq.client,
			QueueName: key.(string),
		})
		return true
	})
	return
}

func (rq *RedisQ) Stop() {
	//等待任务执行完成
	for {
		if len(rq.coroutine) == 0 {
			fmt.Println("async task end")
			return
		}
	}
}
