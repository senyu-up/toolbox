package aws_kafka

import (
	"github.com/segmentio/kafka-go"
)

func (k *Kafka) GetAdmin() *kafka.Conn {
	return k.Conn
}

func (k *Kafka) ListTopics() (map[string]kafka.Partition, error) {
	if k.Conn == nil {
		return nil, ErrKafkaClientIsNil
	}

	partitions, err := k.Conn.ReadPartitions()
	if err != nil {
		return nil, err
	}

	m := map[string]kafka.Partition{} // map topicName -> partitionInfo

	for _, p := range partitions {
		m[p.Topic] = p
	}
	return m, nil
}

func (k *Kafka) DescribeTopic(name string) (detail kafka.Partition, err error) {
	if k.Conn == nil {
		return detail, ErrKafkaClientIsNil
	}

	if ps, err := k.ListTopics(); err != nil {
		return detail, err
	} else if p, ok := ps[name]; ok {
		return p, err
	} else {
		return detail, ErrTopicNotFound
	}
	return detail, err
}

func (k *Kafka) CreateTopic(name string, numPartition int32, numReplication int16) error {
	if k.Conn == nil {
		return ErrKafkaClientIsNil
	}

	var topicConf = kafka.TopicConfig{
		Topic:             name,
		NumPartitions:     int(numPartition),
		ReplicationFactor: int(numReplication),
	}
	return k.Conn.CreateTopics(topicConf)
}

func (k *Kafka) DelTopic(name string) error {
	if k.Conn == nil {
		return ErrKafkaClientIsNil
	}

	return k.Conn.DeleteTopics(name)
}
