//事件总线

package event

type Listener interface {
	//Handle 处理监听事件
	Handle(val interface{}) error
}

type Eventer interface {
	//Name 事件名称
	Name() string
}
