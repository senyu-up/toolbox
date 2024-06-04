package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

func HandleUserLogin(ctx context.Context, msg *sarama.ConsumerMessage) error {
	fmt.Printf("got msg %v\n", msg)
	var data = map[string]string{}
	json.Unmarshal(msg.Value, &data)
	fmt.Printf("consume %s's msg, headers %v, time %s, data %v\n", msg.Topic, msg.Headers, msg.Timestamp.String(), data)
	return nil
}
