package broadcast

import (
	"context"
)

// IAdapterInterface 广播实现
type IAdapterInterface interface {
	Subscribe(ctx context.Context, topic Topic, c chan<- *Message) error
	Publish(msg *Message) error
}
