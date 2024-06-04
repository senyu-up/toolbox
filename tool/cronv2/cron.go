package cronv2

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/trace"
)

type JobFunc func(ctx context.Context) error // func(ctx context.Context) error

type Client struct {
	client  *cron.Cron
	log     cron.Logger
	traceOn bool
	second  bool
}

func New(opts ...CronOption) *Client {
	var c = &Client{traceOn: false}
	for _, opt := range opts {
		opt(c)
	}
	var l = c.log // 用外部的 logger
	if l == nil { // 如果外部不可用，用自己的
		l = NewCronLogger(nil)
	}
	var cronOpts = []cron.Option{cron.WithChain(cron.Recover(l)), cron.WithLogger(l)}
	if c.second {
		cronOpts = append(cronOpts, cron.WithSeconds())
	}
	c.client = cron.New(cronOpts...)

	return c
}

// wrapCronFunc
//
//	@Description: 包装 cron 任务函数，注入 ctx 和trace 信息
//	@param name  body any true "-"
//	@param cmd  body any true "-"
//	@return func()
func (c *Client) wrapCronFunc(name string, cmd JobFunc) func() {
	return func() {
		var ctx = trace.NewTrace()
		traceId, pSpanId := trace.ParseCurrentContext(ctx)
		if c.traceOn {
			var opName = "cron " + name
			spanId := trace.NewSpanID()
			span := trace.NewJaegerSpan(opName, traceId, spanId, pSpanId, nil, nil)
			defer span.Finish()
		}

		if err := cmd(ctx); err != nil {
			logger.Ctx(ctx).SetErr(err).Error("Execute cron job [" + name + "] error")
		}
	}
}

// Register
//
//	@Description: 注册 cron 任务
//	@receiver c
//	@param spec  body any true "-"
//	@param name  body any true "-"
//	@param cmd  body any true "-"
//	@return cron.EntryID
//	@return error
func (c *Client) Register(spec string, name string, cmd JobFunc) (cron.EntryID, error) {
	return c.client.AddFunc(spec, c.wrapCronFunc(name, cmd))
}

func (c *Client) Restart(id cron.EntryID, spec string, name string, cmd JobFunc) (cron.EntryID, error) {
	c.client.Remove(id)
	return c.client.AddFunc(spec, c.wrapCronFunc(name, cmd))
}

func (c *Client) Start() {
	c.client.Start()
}

func (c *Client) Stop() {
	c.client.Stop()
}

// Remove 按照Id移除某个任务
// 注意，删除任务后，如果任务正在执行，不会停止正在运行的任务
func (c *Client) Remove(taskId cron.EntryID) {
	c.client.Remove(taskId)
}

type CronOption func(*Client)

func CronOptionWithTrace(on bool) CronOption {
	return func(option *Client) {
		option.traceOn = on
	}
}

func CronOptionWithSecond(s bool) CronOption {
	return func(option *Client) {
		option.second = s
	}
}

func CronOptionWithLogger(log cron.Logger) CronOption {
	return func(option *Client) {
		option.log = log
	}
}
