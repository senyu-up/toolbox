package adapter

import "context"

type BroadcastBase struct {
	handlers []func(ctx context.Context, msg []byte)
}

/*RegisterHandler
* @Description: 注册消息处理函数
* @param handler
 */
func (b *BroadcastBase) RegisterHandler(handler func(ctx context.Context, msg []byte)) {
	b.handlers = append(b.handlers, handler)
}

/*Broadcast
* @Description: 广播消息
* @param ctx
* @param msg
* @return error
 */
func (b *BroadcastBase) broadcast(ctx context.Context, msg []byte) {
	for i, _ := range b.handlers {
		b.handlers[i](ctx, msg)
	}
}
