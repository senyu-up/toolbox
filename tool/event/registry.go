//注册表

package event

import (
	"github.com/senyu-up/toolbox/tool/logger"
	sy "sync"
	"time"
)

type Option struct {
	//是否异步
	Async bool
	//重试次数
	Retry []*RetryOption

	Log logger.Log // 日志
}

// NewRegistry
// 实例化一个事件注册表
func NewRegistry(o *Option) *registry {
	var retry = make([]*RetryOption, len(o.Retry)+1)
	retry[0] = &RetryOption{T: 0}
	for k, try := range o.Retry {
		retry[k+1] = try
	}
	if o.Log == nil {
		o.Log = logger.GetLogger()
	}
	return &registry{
		book:  make(map[string][]Listener, 0),
		retry: retry,
		async: o.Async,
		log:   o.Log,
	}
}

type RetryOption struct {
	T time.Duration
}

type registry struct {
	lock  sy.RWMutex
	book  map[string][]Listener
	async bool
	retry []*RetryOption
	log   logger.Log // 日志
}

// RegisterListener
// 注册监听者
func (r *registry) RegisterListener(event Eventer, l Listener) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if v, ok := r.book[event.Name()]; ok {
		r.book[event.Name()] = append(v, l)
	} else {
		v = make([]Listener, 1)
		v[0] = l
		r.book[event.Name()] = v
	}

	return nil
}

// TriggerEvent
// 触发事件
func (r *registry) TriggerEvent(event Eventer, val interface{}) (err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	evts, ok := r.book[event.Name()]
	if !ok {
		return DoesNotExistErr
	}

	for _, evt := range evts {
		if r.async {
			r.asyncRun(evt, val)
		} else {
			if err = r.syncRun(evt, val); err != nil {
				return err
			}
		}
	}
	return nil
}

// 异步
func (r *registry) asyncRun(evt Listener, val interface{}) {
	err := pool.Submit(func() {
		for _, try := range r.retry {
			if try.T > 0 {
				time.Sleep(try.T)
			}
			if err := evt.Handle(val); err == nil {
				break
			}
		}
	})
	if err != nil {
		r.log.Error("triggerEvent pool err %s", err)
	}
}

// 同步
func (r *registry) syncRun(evt Listener, val interface{}) (err error) {
	for _, try := range r.retry {
		if try.T > 0 {
			time.Sleep(try.T)
		}
		if err = evt.Handle(val); err == nil {
			break
		}
	}
	return err
}
