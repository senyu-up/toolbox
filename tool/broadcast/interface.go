package broadcast

import (
	"context"
	"github.com/senyu-up/toolbox/tool/config"
)

type BroadcastInter interface {
	Init(cnf config.Broadcast) error
	Subscribe(ctx context.Context) error
	Publish(ctx context.Context, msg []byte) error
	RegisterHandler(handler func(ctx context.Context, msg []byte))
}
