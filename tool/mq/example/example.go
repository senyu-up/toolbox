package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/marshaler"
	"github.com/senyu-up/toolbox/tool/mq/aws_kafka"
	"github.com/senyu-up/toolbox/tool/mq/kafka"
	"github.com/senyu-up/toolbox/tool/runtime"
	"github.com/senyu-up/toolbox/tool/storage"
	"github.com/senyu-up/toolbox/tool/trace"
	"github.com/spf13/cast"
)

var (
	//localBrokers = []string{"127.0.0.1:9092", "127.0.0.1:9093", "127.0.0.1:9094"}
	localBrokers   = []string{"127.0.0.1:9092"} // local test
	pushTopic      = "xh_push"
	syncPushTopic  = "xh_sync_push"
	asyncPushTopic = "xh_async_push"

	lock = sync.Mutex{}
)

var (
	accessId = "AKIXXXXXX"
	secret   = "/F27Uxxxxxxxxxx"
	brokers  = []string{"kafka.us-west-2.amazonaws.com:9198",
		"kafka.us-west-2.amazonaws.com:9198"}
	topic1 = "test_topic"
	topic2 = "test_topic22"
	//topic2 = "data_manager_dqc_develop"
)

// NewKafka
//
//	@Description: 初始化kafka， 同时声明了消费者，和生产者
//	@return *kafka.Kafka
func NewKafka() *kafka.Kafka {
	var ctx = trace.NewTrace() // 初始化trace
	var cfg = &config.KafkaConfig{
		Brokers: localBrokers,
		Timeout: 3,
		Consumers: []*config.KafkaConsumerConfig{
			{
				Brokers: localBrokers,
				Topic:   "test",
				Group:   "test_1",
				Workers: 10,
			},
			{
				Brokers: localBrokers,
				Topic:   pushTopic,
				Group:   "xh_topic_g3",
				Workers: 10,
			},
		},
		Workers: 10,
	}
	// 初始化kafka, 并制定了消息序列化方式 jsonMarshaler, （默认就是json）当然你也可以实现 marshaler.Marshaler 接口作为自定义序列化方式
	kafkaClient, err := kafka.New(cfg, kafka.KafkaOptWithMarshaler(marshaler.JsonMarshaler{}))
	if err != nil {
		fmt.Printf("init kafka err %v", err)
		return nil
	}
	// 项目中别这样写，实际应该加入到 app shutdown callback 中
	defer kafkaClient.Close()

	// 发送消息成功回调
	kafkaClient.Producer().HandleSucceed(func(msg *sarama.ProducerMessage) {
		fmt.Printf("push async msg handle success \n")
	})
	// 捕获错误
	kafkaClient.Producer().HandleError(func(err error) {
		fmt.Printf("push msg handle err %v \n", err)
	})

	// 消费消息
	go func() {
		// 消费指定topic
		kafkaClient.RegisterConsumerHandler("test", "t1", func(ctx context.Context, msg *sarama.ConsumerMessage) error {
			var data = map[string]string{}
			json.Unmarshal(msg.Value, &data)
			fmt.Printf("consume msg topic is %s, headers %v, data %v\n", msg.Topic, msg.Headers, data)
			return nil
		}, func(err error) {
			fmt.Printf("consume msg handle err %v \n", err)
		})

		// 消费
		kafkaClient.RegisterConsumerHandler(pushTopic, "p1", func(ctx context.Context, msg *sarama.ConsumerMessage) error {
			var data = map[string]string{}
			json.Unmarshal(msg.Value, &data)
			fmt.Printf("consume %s's msg, headers %v, time %s, data %v ctx %v\n",
				msg.Topic, msg.Headers, msg.Timestamp.String(), data, ctx)
			return nil
		}, func(err error) {
			fmt.Printf("consume msg handle err %v \n", err)
		})

		kafkaClient.StartConsume() // 开始消费
	}()

	var i = 0
	for {
		// push 消息
		i += 1
		var data = map[string]string{"name": "heihei", "addr": "jingronghui", "times": cast.ToString(i)}
		b, _ := json.Marshal(data)
		var msg = &sarama.ProducerMessage{
			Topic:     pushTopic,
			Value:     sarama.ByteEncoder(b),
			Timestamp: time.Now(),
			Headers: []sarama.RecordHeader{
				{
					Key:   []byte("header-1"),
					Value: []byte("value-1"),
				},
			},
		}

		if i%2 == 0 {
			kafkaClient.PushSyncRaw(ctx, msg) // 发送原始消息
		} else {
			kafkaClient.PushAsync(ctx, "test", data) // 直接发送数据，使用内置的序列化
			//kafkaClient.PushAsyncRaw(ctx, msg)
		}

		time.Sleep(time.Second)
	}
	return nil
}

// ProducerSyncPush
//
//	@Description: 初始化生产者，同步push
//	@param k  body any true "-"
func ProducerSyncPush() {
	var ctx = trace.NewTrace() // 带上trace
	producer, err := kafka.NewProducer(&config.KafkaConfig{
		Brokers: localBrokers,
		Timeout: 3})
	if err != nil {
		fmt.Printf("init kafka producer err %v", err)
		return
	}

	// 同步发送错误信息，会在调用时返回，不需要异步监听

	// push 定制化消息
	var data = map[string]string{"name": "heihei", "addr": "jingronghui"}
	var i = 0
	for {
		lock.Lock()
		i += 1
		data["times"] = cast.ToString(i)
		data["from"] = "sync"
		lock.Unlock()
		var start = time.Now()
		// push 同步消息
		if partition, offset, err := producer.PushSync(ctx, syncPushTopic, data); err != nil {
			fmt.Printf("sync push event err %v \n", err)
			return
		} else {
			fmt.Printf("sync push event sucess partition %d, offset %d\n", partition, offset)
		}
		fmt.Printf("sync push event i: %d duration %d us\n", i, time.Now().Sub(start))
		time.Sleep(time.Second)
	}
}

// ProducerAsyncPush
//
//	@Description: 初始化生产者，异步push
//	@param k  body any true "-"
func ProducerAsyncPush() {
	var ctx = trace.NewTrace() // 带上trace
	producer, err := kafka.NewProducer(&config.KafkaConfig{
		Brokers: localBrokers,
		Timeout: 3})
	if err != nil {
		fmt.Printf("init kafka producer err %v", err)
		return
	}

	// 发送消息成功回调
	producer.HandleSucceed(func(msg *sarama.ProducerMessage) {
		fmt.Printf("async msg push handle success \n")
	})
	// 捕获错误
	producer.HandleError(func(err error) {
		fmt.Printf("push msg handle err %v \n", err)
	})

	// push 定制化消息
	var data = map[string]string{"name": "heihei", "addr": "jingronghui"}
	var i = 0
	for {
		lock.Lock()
		i += 1
		data["times"] = cast.ToString(i)
		data["from"] = "async"
		lock.Unlock()
		// push 定制化消息
		b, _ := json.Marshal(data)
		var testKey = sarama.StringEncoder("Bar")
		var msg = &sarama.ProducerMessage{
			Topic:     asyncPushTopic,
			Value:     sarama.ByteEncoder(b),
			Timestamp: time.Now(),
			Key:       testKey,
		}
		var start = time.Now()
		producer.PushAsyncRaw(ctx, msg)
		fmt.Printf("async push event i: %d, duration %d us\n", i, time.Now().Sub(start).Microseconds())
		time.Sleep(time.Second)
	}
}

func Consumer() {
	var cfg = &config.KafkaConsumerConfig{
		Brokers: localBrokers,
		Topic:   syncPushTopic,
		Group:   "xh_topic_c1",
		Workers: 5,
	}
	conClient, err := kafka.NewConsumer(cfg)
	if err != nil {
		fmt.Printf("init kafka consumer err %v", err)
		return
	}
	// 消费
	conClient.HandleMsg(func(message *sarama.ConsumerMessage) {
		traceId, spanId := kafka.ParserTrace(message.Headers)
		fmt.Printf("consume the msg (%s), topic is  %s , traceId: %s, spanId: %s\n",
			message.Value, message.Topic, traceId, spanId)
	})

	conClient.Start() // 开始消费
}

func Consumer2() {
	var cfg = &config.KafkaConsumerConfig{
		Brokers: localBrokers,
		Topic:   asyncPushTopic,
		Group:   "xh_topic_c1",
		Workers: 5,
	}
	conClient, err := kafka.NewConsumer(cfg)
	if err != nil {
		fmt.Printf("init kafka consumer err %v", err)
		return
	}
	// 消息处理 callback
	var msgCallBack = func(message *sarama.ConsumerMessage) {
		traceId, spanId := kafka.ParserTrace(message.Headers)
		fmt.Printf("consume the msg (%s), topic is  %s , traceId: %s, spanId: %s\n",
			message.Value, message.Topic, traceId, spanId)
	}
	// 消费
	conClient.HandleMsg(msgCallBack)

	conClient.Start() // 开始消费 阻塞
}

func Consumer3() {
	client, err := kafka.NewConsumerWithHandler(&config.KafkaConsumerConfig{
		Brokers: localBrokers,
		Topic:   asyncPushTopic,
		Group:   "xh_topic_c1",
		Workers: 5,
	}, func(message *sarama.ConsumerMessage) {
		traceId, spanId := kafka.ParserTrace(message.Headers)
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
	time.Sleep(10 * time.Second)
}

func Consumer4() {
	var cfg = &config.KafkaConsumerConfig{
		Brokers: localBrokers,
		Topic:   asyncPushTopic,
		Group:   "xh_topic_c1",
		Workers: 5,
	}
	conClient, err := kafka.NewConsumer(cfg)
	if err != nil {
		fmt.Printf("init kafka consumer err %v", err)
		return
	}
	conClient.HandleError(func(err error) {
		fmt.Printf("cb consume msg err %v \n", err)
	})
	conClient.HandleClose(func() {
		fmt.Printf("cb consume msg close \n")
	})
	conClient.HandleReBalance(func() {
		fmt.Printf("cb consume msg rebalance \n")
	})

	// 消息处理 callback
	var msgCallBack = func(message *sarama.ConsumerMessage) {
		traceId, spanId := kafka.ParserTrace(message.Headers)
		fmt.Printf("consume the msg (%s), topic is  %s , traceId: %s, spanId: %s\n",
			message.Value, message.Topic, traceId, spanId)
	}
	// 消费
	conClient.HandleMsg(msgCallBack)

	go conClient.Start() // 开始消费 阻塞

	runtime.WaitForStopSignal(
		func() {
			conClient.Close()
		})
}

func Case1() {
	NewKafka()
}

func Case2() {
	go Consumer()

	ProducerSyncPush()
}

func Case3() {
	go Consumer()
	go Consumer2()

	go ProducerAsyncPush()
	ProducerSyncPush()
}

func Case4() {
	go Consumer3()
	ProducerAsyncPush()

}

func Case5() {
	go ProducerAsyncPush()
	Consumer4()
}

func AwsPush1() {
	var topicName = "test_topic22"
	var conf = &config.AwsKafkaConfig{
		Brokers: brokers,

		Topic: topicName,
	}
	conf.SASL.Enable = true
	conf.SASL.Mechanism = "AWS_MSK_IAM"
	conf.SASL.Region = storage.AWSUwRegion
	conf.SASL.AccessId = accessId
	conf.SASL.SecretKey = secret

	var p = aws_kafka.NewProducer(conf)

	var ctx = context.Background()
	err := p.Producer.WriteMessages(ctx, kafkago.Message{
		Topic: topicName,
		Value: []byte("hello world"),
	})
	log.Printf("write msg err %v", err)
}

func AwsConsume() {
	var topicName = "test_topic22"
	var conf = &config.AwsKafkaConfig{
		GroupId: "xh_topic_c1",

		Brokers: brokers,

		Topic: topicName,
	}
	conf.SASL.Enable = true
	conf.SASL.Mechanism = "AWS_MSK_IAM"
	conf.SASL.Region = storage.AWSUwRegion
	conf.SASL.AccessId = accessId
	conf.SASL.SecretKey = secret

	var c = aws_kafka.NewConsumer(conf, aws_kafka.KafkaOptWithTopic(topicName))
	for {
		m, err := c.Reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("read msg err %v", err)
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	if err := c.Reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

func AwsGetTopic() {
	var ctx = context.Background()
	var topicName = "test_topic22"
	var conf = &config.AwsKafkaConfig{
		Brokers: brokers,

		Topic: topicName,
	}
	conf.SASL.Enable = true
	conf.SASL.Mechanism = "AWS_MSK_IAM"
	conf.SASL.Region = storage.AWSUwRegion
	conf.SASL.AccessId = accessId
	conf.SASL.SecretKey = secret

	k, err := aws_kafka.New(ctx, conf)
	if err != nil {
		log.Printf("new kafka err %v", err)
		return
	}
	defer k.Conn.Close()

	k.Conn.ReadPartitions()

	partitions, err := k.Conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	for k := range m {
		fmt.Println(k)
	}
}

func Case6() {
	//AwsGetTopic()
	AwsPush1()
	AwsConsume()
	time.Sleep(10 * time.Second)
}

func getAwsKafkaConf() config.AwsKafkaConfig {
	var conf = config.AwsKafkaConfig{
		Timeout: 10,

		Brokers: brokers,
		Workers: 3,
		Async:   false,
		TraceOn: true,

		//Topic: topic2,
	}
	conf.SASL.Enable = true
	conf.SASL.Mechanism = "AWS_MSK_IAM"
	conf.SASL.Region = storage.AWSUwRegion
	conf.SASL.AccessId = accessId
	conf.SASL.SecretKey = secret

	return conf
}

func getLocalAwsKafkaConf() config.AwsKafkaConfig {
	var localBroker = []string{"localhost:9092"}
	return config.AwsKafkaConfig{
		Timeout: 10,

		Brokers: localBroker,
		Workers: 3,
		Async:   false,
		TraceOn: true,

		//Topic: topic2,
	}
}

func ConsumeAwsKafkaMsg(ctx context.Context) {
	var conf = getAwsKafkaConf()
	log.Printf("get aws kafka conf %+v", conf)
	k, err := aws_kafka.New(ctx, &conf)
	if err != nil {
		log.Printf("new kafka err %v", err)
		return
	}

	// 订阅 topic2， 并设置 consumer group, 设置消息订阅callback， 和消息错误callback
	k.RegisterConsumerHandler(topic2, topic2+"_c", func(ctx context.Context, msg kafkago.Message) error {
		// 消费 consumer 代码逻辑
		log.Printf("consume msg %s", string(msg.Value))
		return nil
	}, func(err error) {
		// 错误处理
		log.Printf("consume msg err %v", err)
	})

	k.StartConsume() // 启动协程开始消费
	<-ctx.Done()     // 等待退出
	k.Close()        // 关闭消费
}

func PushAwsKafkaMsg(ctx context.Context) {
	var conf = getAwsKafkaConf()
	k, err := aws_kafka.New(ctx, &conf)
	if err != nil {
		log.Printf("new kafka err %v", err)
		return
	}

	var ticker = time.Tick(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			k.Close()
			return
		case <-ticker:
			k.PushSyncRawMsgs(ctx, []kafkago.Message{{
				Topic: topic2,
				Value: []byte("hello world"),
			}})

			time.Sleep(time.Second * 3)

			var data = map[string]string{"name": "heihei", "addr": "jingronghui"}
			k.PushSyncMsgs(ctx, topic2, data)
		}
	}
}

// aws kafka 发送消息， 和 消费消息
func Case7() {
	ctx, cancel := context.WithCancel(context.Background())
	go ConsumeAwsKafkaMsg(ctx)
	go PushAwsKafkaMsg(ctx)

	runtime.WaitForStopSignal(func() {
		cancel() // 目测没用
	})
	time.Sleep(time.Second * 10)
}

func PushAwsKafkaMsgCb(ctx context.Context) {
	var conf = getAwsKafkaConf()
	k, err := aws_kafka.New(ctx, &conf)
	if err != nil {
		log.Printf("new kafka err %v", err)
		return
	}

	// 设置同步推送消息成功 callback
	k.SyncProducer().HandleSucceed(func(msg kafkago.Message) {
		log.Printf("sync push msg succeed %s", string(msg.Value))
	})

	// 设置同步推送消息失败 callback
	k.SyncProducer().HandleError(func(err error) {
		log.Printf("async push msg err %v", err)
	})

	// 异步推送消息
	k.PushAsyncRawMsgs(ctx, []kafkago.Message{{
		Topic: topic2,
		Value: []byte("hello world"),
	}})

	// 同步推送消息
	k.PushSyncRawMsgs(ctx, []kafkago.Message{{
		Topic: topic2,
		Value: []byte("hello world"),
	}})
}

func GetListByAwsKafka() {
	var conf = getLocalAwsKafkaConf()

	k, err := aws_kafka.New(context.Background(), &conf)
	if err != nil {
		log.Printf("new kafka err %v", err)
		return
	}
	// 列出所有 topic
	list, err := k.ListTopics()
	if err != nil {
		log.Printf("list kafka topic err %v", err)
	} else {
		log.Printf("list %+v", list)
	}

	// 创建 topic
	err = k.CreateTopic(topic2, 1, 1)
	if err != nil {
		log.Printf("create kafka topic err %v", err)
	}
}

func Case8() {
	GetListByAwsKafka()
}

func Case71() {
	PushAwsKafkaMsgCb(context.Background())
}

func Case72() {
	ConsumeAwsKafkaMsg(context.Background())
}

func main() {
	// new kafka client
	//Case1()

	// 同步发消息， 消费
	//Case2()

	//同步、异步发消息，消费，对比
	//Case3()

	// kafka 异步发送， 消费
	//Case4()

	// consumer call back
	//Case5()

	// aws kafka
	//Case6()

	// aws kafka sdk
	//Case7()
	//Case71() // aws kafka push
	Case72() // aws kafka consumer

	// local aws test
	//Case8()
}
