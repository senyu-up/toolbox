package index

import (
	"github.com/senyu-up/toolbox/example/internal/event"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/mq/aws_kafka"
)

func RegisterKafkaConsumer(consumer *aws_kafka.Kafka) (err error) {
	err = consumer.RegisterConsumerHandler("data_manager_dqc_develop", "data_manager_dqc_develop_c2",
		event.AwsMsgHandle, func(err error) {
			logger.SetErr(err).Error("consumer topic: data_manager_dqc_develop error")
		})
	if err != nil {
		return err
	}

	return
}
