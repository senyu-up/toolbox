package job

import (
	"context"
	"errors"
	"fmt"
	"github.com/senyu-up/toolbox/tool/logger"
	"reflect"
)

var AsyncPool = NewAsyncWork(context.TODO(), OptWithMaxJob(1<<10))

type Job struct {
	f  reflect.Value
	in []reflect.Value
}

type AsyncJob struct {
	ctx     context.Context
	channel chan Job
	finish  chan struct{}
	close   bool
	maxJobs int64

	log logger.Log
}

/*
	使用当前进程执行且保证了，且保证了所有任务都能被正常执行
*/

func NewAsyncWork(ctx context.Context, opts ...JobOption) *AsyncJob {
	job := &AsyncJob{ctx: ctx}

	// option apply
	for _, opt := range opts {
		opt(job)
	}

	// 缺省值补全
	job.channel = make(chan Job, job.maxJobs)
	job.finish = make(chan struct{}, job.maxJobs)
	if job.log == nil {
		job.log = logger.GetLogger()
	}

	go job.process()
	return job
}

func (job *AsyncJob) SubJob(f interface{}, args ...interface{}) error {
	if len(job.channel) >= int(job.maxJobs) {
		return errors.New("async job is full")
	}
	if job.close {
		return errors.New("bye !")
	}
	tf := reflect.ValueOf(f)
	if tf.Kind() != reflect.Func {
		return errors.New("func Type error")
	}
	isTypeCheck := true
	if argsNum := tf.Type().NumIn(); argsNum != len(args) {
		switch {
		//判断最后一个参数是否为slice,如果为slice则放弃校验传入参数
		case argsNum > 0 && tf.Type().In(argsNum-1).Kind() != reflect.Slice:
			return errors.New("args is err")
			//处理方法不存在传入参数
		case argsNum == 0:
			args = nil
		default:
			isTypeCheck = false
		}
	}
	in := make([]reflect.Value, len(args))
	for idx, arg := range args {
		val := reflect.ValueOf(arg)
		if val.Kind() != tf.Type().In(idx).Kind() && isTypeCheck {
			return fmt.Errorf("in %d: kind invalid", idx)
		}
		in[idx] = val
	}
	//等待完成的任务
	job.finish <- struct{}{}
	job.channel <- Job{
		f:  tf,
		in: in,
	}
	return nil
}

func (job *AsyncJob) process() {
	for {
		select {
		//接受到任务终止信号
		case work := <-job.channel:
			go func(work Job) {
				defer func() {
					if err := recover(); err != nil {
						job.log.Error("func:%s,args:%v,err:%v", work.f.Type().String(), work.in, err)
					}
					<-job.finish
				}()
				call := work.f.Call(work.in)
				if logger.Level() >= logger.LevelInformational {
					in, out := make([]interface{}, 0, 4), make([]interface{}, 0, 1)
					for _, val := range call {
						out = append(out, val.Interface())
					}
					for _, val := range work.in {
						in = append(in, val.Interface())
					}
					job.log.Info("%s:in->%v,out->%v", work.f.Type().String(), in, out)
				}
			}(work)
		}
	}
}

func (job *AsyncJob) WaitingStop() {
	job.close = true
	for len(job.finish) != 0 {
		//一致等待所有任务完成
	}
}
