package queue

import (
	"fmt"
	"sync"
)

// SqQueue 结构体定义
type SqQueue struct {
	lock    sync.Mutex
	data    []interface{}
	maxSize int
	front   int
	rear    int
}

// New 新建空队列
func New(size int) *SqQueue {
	return &SqQueue{
		front:   0,
		rear:    0,
		lock:    sync.Mutex{},
		maxSize: size,
		data:    make([]interface{}, size),
	}
}

// Length 队列长度
func (q *SqQueue) Length() int {
	return (q.rear + q.maxSize - q.front) % q.maxSize
}

func (q *SqQueue) IsFull() bool {
	return (q.rear+1)%q.maxSize == q.front
}

func (q *SqQueue) IsEmpty() bool {
	return q.rear == q.front
}

// Enqueue 入队
func (q *SqQueue) Enqueue(e interface{}) error {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.IsFull() {
		return fmt.Errorf("quque is full")
	}
	q.data[q.rear] = e
	q.rear = (q.rear + 1) % q.maxSize

	return nil
}

// Dequeue 出队
func (q *SqQueue) Dequeue() (e interface{}, err error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.IsEmpty() {
		return e, fmt.Errorf("quque is empty")
	}
	e = q.data[q.front]
	q.front = (q.front + 1) % q.maxSize

	return e, nil
}
