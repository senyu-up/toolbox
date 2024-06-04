package index

import (
	"github.com/senyu-up/toolbox/example/internal/event"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/mq/kafka"
)

func RegisterConsumer(consumer *kafka.Kafka) (err error) {
	err = consumer.RegisterConsumerHandler("xh_push", "g1", event.HandleUserLogin, func(err error) {
		logger.SetErr(err).Error("consumer topic: xh_push error")
	})
	if err != nil {
		return err
	}

	return
}
