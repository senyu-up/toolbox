package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/senyu-up/toolbox/tool/config"
)

var (
	localBrokers = []string{"127.0.0.1:9092", "127.0.0.1:9093", "127.0.0.1:9094"}
)

func TestNew(t *testing.T) {
	var (
		pushTopic = "xh_push"
	)
	kafkaClient, err := New(&config.KafkaConfig{
		Brokers: localBrokers,
		Timeout: 3,
		Consumers: []*config.KafkaConsumerConfig{
			{
				Brokers: localBrokers,
				Topic:   "test",
				Group:   "test_1",
				Workers: 10,
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	defer kafkaClient.Close()

	// 发送 异步消息成功回调
	kafkaClient.Producer().HandleSucceed(func(msg *sarama.ProducerMessage) {
		fmt.Printf("push msg handle success %v \n", msg)
	})
	// 发送 异步消息错误回调
	kafkaClient.Producer().HandleError(func(err error) {
		fmt.Printf("push msg handle err %v \n", err)
	})

	// 消费消息
	go func() {
		kafkaClient.RegisterConsumerHandler("xh_consume", "test_1", func(ctx context.Context, msg *sarama.ConsumerMessage) error {
			fmt.Printf("consume msg  %v \n", msg)
			return nil
		}, func(err error) {
			fmt.Print("consume msg handle err \n")
		})

		kafkaClient.StartConsume()
	}()

	for {
		// push 消息
		var ctx = context.TODO()
		var data = map[string]string{"name": "heihei", "addr": "jingronghui"}
		b, _ := json.Marshal(data)
		var msg = &sarama.ProducerMessage{
			Topic: pushTopic,
			Value: sarama.ByteEncoder(b),
		}

		kafkaClient.PushSyncRaw(ctx, msg)

		time.Sleep(time.Second)
	}

	time.Sleep(time.Second * 15)
	return
}

func TestKafkaClose(t *testing.T) {
	var (
		pushTopic = "xh_push"
	)
	kafkaClient, err := New(&config.KafkaConfig{
		Brokers: localBrokers,
		Timeout: 3,
		Consumers: []*config.KafkaConsumerConfig{
			{
				Brokers: localBrokers,
				Topic:   "test",
				Group:   "test_1",
				Workers: 10,
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	defer kafkaClient.Close()

	// 发送 异步消息成功回调
	kafkaClient.Producer().HandleSucceed(func(msg *sarama.ProducerMessage) {
		fmt.Printf("push msg handle success %v \n", msg)
	})
	// 发送 异步消息错误回调
	kafkaClient.Producer().HandleError(func(err error) {
		fmt.Printf("push msg handle err %v \n", err)
	})

	// 消费消息
	go func() {
		kafkaClient.RegisterConsumerHandler("xh_consume", "xh_g", func(ctx context.Context, msg *sarama.ConsumerMessage) error {
			fmt.Printf("consume msg  %v \n", msg)
			return nil
		}, func(err error) {
			fmt.Print("consume msg handle err \n")
		})

		go kafkaClient.StartConsume()
		time.Sleep(time.Second)
		kafkaClient.StartConsume()
	}()

	// push 消息
	var ctx = context.TODO()
	var data = map[string]string{"name": "heihei", "addr": "jingronghui"}
	b, _ := json.Marshal(data)
	var msg = &sarama.ProducerMessage{
		Topic: pushTopic,
		Value: sarama.ByteEncoder(b),
	}

	kafkaClient.PushSyncRaw(ctx, msg)

	time.Sleep(time.Second * 5)
	return
}

func TestKafkaClose2(t *testing.T) {
	client, err := NewConsumerWithHandler(&config.KafkaConsumerConfig{
		Brokers: localBrokers,
		Topic:   "xh_topic",
		Group:   "xh_topic_c1",
		Workers: 5,
	}, func(message *sarama.ConsumerMessage) {
		traceId, spanId := ParserTrace(message.Headers)
		fmt.Printf("consume the msg (%s), topic is  %s , traceId: %s, spanId: %s\n",
			message.Value, message.Topic, traceId, spanId)
	})
	if err != nil {
		fmt.Printf("init kafka consumer err %v", err)
		return
	} else {
		fmt.Printf("init kafka consumer success \n")
		time.Sleep(10 * time.Second)
		client.Close()
	}
	client.Close()
	time.Sleep(1 * time.Second)
}
