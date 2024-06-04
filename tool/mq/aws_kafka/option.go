package aws_kafka

import (
	"github.com/segmentio/kafka-go"
	"github.com/senyu-up/toolbox/tool/marshaler"
	"log"
)

type KafkaOpt struct {
	GroupId string `yaml:"-"` // 消费组名称
	Topic   string `yaml:"-"` // topic 的名称
	Async   bool   `yaml:"-"` // 是否异步发送

	// 自定义序列化方式, 为空默认使用 marshaler.JsonMarshaler
	m marshaler.Marshaler
	// 自定义分区方式, 为空默认使用 sarama.NewHashPartitioner
	pc kafka.Balancer
}

type KafkaOption func(*KafkaOpt)

// Marshaler
func KafkaOptWithMarshaler(m marshaler.Marshaler) KafkaOption {
	return func(option *KafkaOpt) {
		option.m = m
	}
}

func KafkaOptWithTopic(t string) KafkaOption {
	return func(option *KafkaOpt) {
		option.Topic = t
	}
}

// GroupId
func KafkaOptWithGroupId(g string) KafkaOption {
	return func(option *KafkaOpt) {
		option.GroupId = g
	}
}

// Async
func KafkaOptWithAsync(a bool) KafkaOption {
	return func(option *KafkaOpt) {
		option.Async = a
	}
}

func getKafkaBalancerByStr(lb string) kafka.Balancer {
	switch lb {
	case "hash":
		return &kafka.Hash{}
	case "round_robin":
		return &kafka.RoundRobin{}
	case "least_bytes":
		return &kafka.LeastBytes{}
	case "crc32balancer":
		return &kafka.CRC32Balancer{}
	case "murmur2_balancer":
		return &kafka.Murmur2Balancer{}
	case "reference_hash":
		return &kafka.ReferenceHash{}
	default:
		log.Fatal("kafka producer balancer error: ", lb)
		return nil
	}
}

// Balancer
// 发送方，负载均衡策略, 默认为 Hash，可选：hash,least_bytes,round_robin,crc32balancer,murmur2_balancer,reference_hash
func KafkaOptWithBalancer(lb string) KafkaOption {
	return func(option *KafkaOpt) {
		option.pc = getKafkaBalancerByStr(lb)
	}
}
