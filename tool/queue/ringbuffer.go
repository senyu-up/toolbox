package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var RingBufferOverFlowError = errors.New("Ring Buffer Over Flow Error. ")
var RingBufferStoppedError = errors.New("Ring Buffer Stopped. ")

// 扩缩容因子
const ExpansionFactor = 1024
const ExpansionCounter = 5

// RingBuffer 环形 []byte buffer
type RingBuffer struct {
	cap     int
	head    int           //头指针 head <---> tail
	max     int           //最大容量
	tail    int           //尾指针
	list    []interface{} //内容队列
	length  int
	lock    sync.Mutex
	oldList []interface{}

	start bool

	//缩容
	count    int
	shrink   bool
	interval time.Duration

	cancel context.CancelFunc
	ctx    context.Context

	//ticker
	rTicker *time.Ticker

	//ring buffer name
	name string
}

type RingBufferConf struct {
	// 队列名
	Name string
	// 容量
	Cap int
	// 最大值, 容量 <= 最大值
	Max int
	// 定时触发检测缩容的间隔
	Interval time.Duration
}

var bufferMap = make(map[string]*RingBuffer)
var lock = sync.Mutex{}

func NewRingBuffer(conf *RingBufferConf) *RingBuffer {
	lock.Lock()
	defer lock.Unlock()
	// 基于name去重
	if buf, exists := bufferMap[conf.Name]; exists {
		return buf
	}

	if conf.Max < conf.Cap {
		conf.Max = conf.Cap
	}

	/*最小周期1分钟*/
	if conf.Interval > 0 && conf.Interval < time.Minute {
		conf.Interval = time.Minute
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	buf := &RingBuffer{
		cap:      conf.Cap,
		max:      conf.Max,
		list:     make([]interface{}, conf.Cap),
		oldList:  nil,
		lock:     sync.Mutex{},
		name:     conf.Name,
		interval: conf.Interval,
		start:    true,
		head:     0,
		tail:     0,
		cancel:   cancelFunc,
		ctx:      ctx,
	}
	bufferMap[conf.Name] = buf
	if buf.interval > 0 {
		//fixme
		buf.rTicker = time.NewTicker(buf.interval)
		go buf.reduce()
	}

	return buf
}

func (p *RingBuffer) toOldList() {
	// 判断队列是否为空
	if p.head == p.tail && p.length == 0 {
		return
	}
	i := 0
	p.oldList = make([]interface{}, 0, p.length)
	for i < p.length {
		p.oldList = append(p.oldList, p.list[p.head])
		i++
		p.head++
		if p.head >= p.cap {
			p.head %= p.cap
		}
		if p.head == p.tail {
			break
		}
	}
}

func (p *RingBuffer) Put(data interface{}) error {
	if !p.start {
		return RingBufferStoppedError
	}
	if p.length >= p.max {
		return RingBufferOverFlowError
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	// 判断是否需要进行拓容
	if p.length == p.cap && p.cap < p.max {
		p.toOldList()

		p.oldList = p.list
		if p.cap <= ExpansionFactor {
			p.cap *= 2
		} else if p.cap > ExpansionFactor && p.cap < p.max {
			tmp := float32(p.cap) * 1.2
			p.cap = int(tmp)
			if p.cap > p.max {
				p.cap = p.max
			}
		}
		p.list = make([]interface{}, p.cap, p.cap)
		copy(p.list, p.oldList)
		// 这里需要对指针进行重置
		p.head = 0
		p.tail = p.length
		p.oldList = nil
	}
	p.tail %= p.cap
	p.list[p.tail] = data
	p.tail++
	p.length++

	if p.tail >= p.cap {
		p.tail = 0
	}

	return nil
}

func (p *RingBuffer) GetN(c int) ([]interface{}, int) {
	p.lock.Lock()
	defer p.lock.Unlock()
	// 判断队列是否为空
	if p.head == p.tail && p.length == 0 {
		return nil, 0
	}
	// 也许拿不到c个数据
	if p.length < c {
		c = p.length
	}
	i := 0
	rs := make([]interface{}, 0, c)
	for i < c {
		rs = append(rs, p.list[p.head])
		i++
		p.head++
		if p.head >= p.cap {
			p.head %= p.cap
		}
		if p.head == p.tail {
			break
		}
	}

	p.length = (p.tail - p.head + p.cap) % p.cap

	return rs, len(rs)
}

func (p *RingBuffer) AssignN(target []interface{}, n int) int {
	rs, _ := p.GetN(n)

	return copy(target, rs)
}

func (p *RingBuffer) Length() int {
	return p.length
}

func (p *RingBuffer) Cap() int {
	return p.cap
}

func (p *RingBuffer) reduce() {
	for true {
		select {
		case <-p.rTicker.C:
			fmt.Println("Ring Buffer [", p.name, "] Ring buffer current cap:[", p.cap, "],", "elements number:[", p.length, "]")
			if p.cap < ExpansionFactor {
				p.count = 0
				p.shrink = false
				continue
			}
			target := int(float32(p.cap) * 0.8)
			if p.length < target {
				if p.shrink && p.count >= ExpansionCounter {
					p.lock.Lock()
					if p.length >= target {
						p.lock.Unlock()
						continue
					}
					p.toOldList()
					p.list = make([]interface{}, target, target)
					copy(p.list, p.oldList)
					p.oldList = nil
					p.tail = p.length
					p.head = 0
					p.cap = target
					p.lock.Unlock()
					p.shrink = false
					p.count = 0
					continue
				}
				p.shrink = true
				p.count += 1

				continue
			}
			p.count = 0
			p.shrink = false
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *RingBuffer) Free() {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.list = nil
	p.start = false
	p.length = 0
	p.tail = 0
	p.head = 0
	p.cancel()
	if p.rTicker != nil {
		p.rTicker.Stop()
	}

	fmt.Println("[Ring Buffer] [", p.name, "] Closed!")
}
