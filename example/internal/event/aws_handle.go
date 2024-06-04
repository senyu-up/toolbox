package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
)

func AwsMsgHandle(ctx context.Context, msg kafka.Message) error {
	fmt.Printf("got aws kafka msg %s\n", msg.Value)
	var data = map[string]string{}
	json.Unmarshal(msg.Value, &data)
	fmt.Printf("consume %s's msg, headers %v, time %s, data %v\n", msg.Topic, msg.Headers, msg.Time.String(), data)
	return nil
}
