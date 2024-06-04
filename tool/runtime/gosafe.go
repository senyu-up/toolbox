package runtime

import (
	"context"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/senyu-up/toolbox/tool/logger"
	"os"
	"os/signal"
	"reflect"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
)

func Recover(ctx context.Context, module string, data interface{}) {
	if err := recover(); err != nil {
		logger.Ctx(ctx).SetExtra(logger.E().Any("err", err).Any("data", data)).Crit(module)
	}
}

// GOSafe
//
//	@Description:
//	@param ctx  body any true "-"
//	@param tag  body any true "-"
//	@param f  body any true "-"
func GOSafe(ctx context.Context, tag string, f func()) {
	go func(gCtx context.Context) {
		log := logger.Ctx(ctx)
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.SetExtra(logger.E().Any("err", err).String("tag", tag)).Notify().Crit(stack) // 记录日志，发企业机器人
			}
		}()
		f()
	}(ctx)
}

func Go(ctx context.Context, tag string, f func()) {
	go func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Ctx(ctx).SetExtra(logger.E().Any("err", err)).Crit(tag)
			}

			f()
		}()
	}(ctx)
}

func GoWithPanic(ctx context.Context, tag string, f func()) {
	go func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Ctx(ctx).SetExtra(logger.E().Any("err", err)).Notify().Crit(tag)
			}

			f()
		}()
	}(ctx)
}

type PoolConfig struct {
	// 模块名
	Module string
	// * 协程池大小
	Size int
	// * 处理句柄
	Handler func(interface{})
	// ? 过期时间, 默认1h
	ExpiryDuration time.Duration
	// ? 是否预分配
	PreAlloc bool
	// ? 最大阻塞的任务量, 当 Nonblocking 设置为true, 此参数无效
	MaxBlockingTasks int
	// ? panic 处理句柄
	PanicHandler func(interface{})
	// ? 是否在panic时发送机器人消息, 默认 false
	PanicNotify bool
	// ? 非阻塞模式
	Nonblocking bool
	// ? 自定义日志打印
	logger ants.Logger
}

func dftPanicHandler(notify bool, module string) (f func(interface{})) {
	if notify {
		return func(i interface{}) {
			logger.SetExtra(logger.E().Any("msg", i)).Notify().Crit(module)
		}
	} else {
		return func(i interface{}) {
			logger.SetExtra(logger.E().Any("msg", i)).Crit(module)
		}
	}
}

type poolLogger struct {
	ctx    context.Context
	notify bool
	module string
}

func (l poolLogger) Printf(format string, args ...interface{}) {
	logger.Ctx(l.ctx).SetExtra(logger.E().String("module", l.module)).Info(format, args...)
}

func Pool(ctx context.Context, cfg PoolConfig) *ants.PoolWithFunc {
	opt := ants.Options{
		ExpiryDuration:   cfg.ExpiryDuration,
		PreAlloc:         cfg.PreAlloc,
		MaxBlockingTasks: cfg.MaxBlockingTasks,
		Nonblocking:      false,
		PanicHandler:     nil,
		Logger:           nil,
	}

	if cfg.ExpiryDuration <= 0 {
		cfg.ExpiryDuration = time.Hour
	}

	opt.ExpiryDuration = cfg.ExpiryDuration

	if cfg.PanicHandler == nil {
		cfg.PanicHandler = dftPanicHandler(cfg.PanicNotify, cfg.Module)
	}
	opt.PanicHandler = cfg.PanicHandler

	if cfg.logger == nil {
		cfg.logger = poolLogger{ctx: ctx, notify: cfg.PanicNotify, module: cfg.Module}
	}
	opt.Logger = cfg.logger

	options := ants.WithOptions(opt)
	pool, err := ants.NewPoolWithFunc(cfg.Size, cfg.Handler, options)
	if err != nil {
		panic(err)
	}

	return pool
}

func format(v ...interface{}) string {
	return fmt.Sprintf("%v \033[0m\n", strings.TrimRight(fmt.Sprintln(v...), "\n"))
}

func RecoverPanic(where string) {
	if x := recover(); x != nil {
		stack := string(debug.Stack())
		errorMsg := format("caught panic in ", where, x, stack)
		logger.Error(errorMsg)
	}
}

func GetType(msg interface{}) string {
	if t := reflect.TypeOf(msg); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func WaitForStopSignal(cb func()) {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan // wait for SIGINT or SIGTERM
	cb()
}
