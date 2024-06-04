package kafka

import (
	"github.com/IBM/sarama"
	"github.com/senyu-up/toolbox/tool/marshaler"
)

type KafkaOpt struct {
	m marshaler.Marshaler
	// 自定义分区方式, 为空默认使用 sarama.NewHashPartitioner
	pc sarama.PartitionerConstructor
}

type KafkaOption func(*KafkaOpt)

func KafkaOptWithMarshaler(m marshaler.Marshaler) KafkaOption {
	return func(option *KafkaOpt) {
		option.m = m
	}
}

// 自定义分区方式, 为空默认使用 sarama.NewHashPartitioner
func LogOptWithPartitionerConstructor(pc sarama.PartitionerConstructor) KafkaOption {
	return func(option *KafkaOpt) {
		option.pc = pc
	}
}
