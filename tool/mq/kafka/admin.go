package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"strings"
)

func (k *Kafka) GetAdmin() (sarama.ClusterAdmin, error) {
	var admCfg = *k.producer.config
	admCfg.Metadata.Full = true
	admin, err := sarama.NewClusterAdmin(k.conf.Brokers, &admCfg)
	if err != nil {
		return nil, err
	} else {
		return admin, nil
	}
}

func (k *Kafka) ListTopics() (map[string]sarama.TopicDetail, error) {
	adm, err := k.GetAdmin()
	if err != nil {
		return nil, err
	}
	defer adm.Close()
	return adm.ListTopics()
}

func (k *Kafka) DescribeTopic(name string) (detail *sarama.TopicMetadata, err error) {
	adm, err := k.GetAdmin()
	if err != nil {
		return nil, err
	}
	defer adm.Close()
	topics, err := adm.DescribeTopics([]string{name})
	if err != nil {
		return nil, err
	}
	if topics == nil {
		return nil, fmt.Errorf("topic not found")
	}
	return topics[0], nil
}

func (k *Kafka) DelTopic(name string) error {
	adm, err := k.GetAdmin()
	if err != nil {
		return err
	}
	defer adm.Close()

	err = adm.DeleteTopic(name)
	if err != nil {
		// 针对不存在的情况也认为删除成功
		if strings.Contains(err.Error(), "does not exist") {
			return nil
		}
	}
	return err
}

func (k *Kafka) CreateTopic(name string, numPartition int32, numReplication int16) error {
	adm, err := k.GetAdmin()
	if err != nil {
		return err
	}
	defer adm.Close()
	metadata, err := adm.DescribeTopics([]string{name})
	if err != nil {
		return err
	} else if metadata == nil {
		if metadata[0].Err == sarama.ErrUnknownTopicOrPartition {
			// 创建topic
			return adm.CreateTopic(name, &sarama.TopicDetail{
				NumPartitions:     numPartition,
				ReplicationFactor: numReplication,
			}, false)
		} else {
			if metadata[0].Err != sarama.ErrNoError {
				return fmt.Errorf("topic already exist")
			} else {
				return metadata[0].Err
			}
		}
	}
	return err
}
