package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/senyu-up/toolbox/tool/trace"
)

type baseLogger struct {
	adapterName  string
	callBackName string
	adapter      Driver
	callBack     Driver

	ctx     context.Context
	extra   *Extras
	exLock  sync.Mutex
	err     error
	traceId string
	spanId  string
	notify  bool

	output          *nameLogger
	appName         string
	callDepth       int
	timeFormat      string
	usePath         string
	ShowCallerLevel LogLevel
}

func newBaseLogger(name string, d Driver, opts ...LogOption) *baseLogger {
	var bl = &baseLogger{
		adapterName:     name,
		adapter:         d,
		appName:         "[NONE]",
		ShowCallerLevel: LevelTrace,
		timeFormat:      LogTimeDefaultFormat,
		extra:           NewExtras(),
		exLock:          sync.Mutex{},
		output:          &nameLogger{name: name, config: ""},
	}
	// 应用 opts
	for _, opt := range opts {
		opt(bl)
	}
	return bl
}

func (b *baseLogger) Clone() Log {
	return &baseLogger{adapterName: b.adapterName,
		adapter:         b.adapter,
		output:          b.output,
		appName:         b.appName,
		callDepth:       b.callDepth,
		timeFormat:      b.timeFormat,
		usePath:         b.usePath,
		ShowCallerLevel: b.ShowCallerLevel,
		callBack:        b.callBack,
		extra:           NewExtras(),
	}
}

// AdapterName
//
//	@Description: 获取当前日志的类型
//	@receiver b
//	@return string
func (b *baseLogger) AdapterName() string {
	return b.adapterName
}

func (b *baseLogger) SetCallBack(cb Driver) Log {
	b.callBack = cb
	return b
}

// Notify
//
//	@Description: 是否同时发送 qw机器人通知, 调用则启用
//	@receiver b
//	@return Log
func (b *baseLogger) Notify() Log {
	b.notify = true
	return b
}

// Ctx
//
//	@Description: 把ctx 传入 logger, 并且从 ctx 解析出 traceId，spanId
//	@receiver b
//	@param ctx  body any true "-"
//	@return Log
func (b *baseLogger) Ctx(ctx context.Context) Log {
	b.ctx = ctx

	// ctx 解析 traceId
	b.traceId, b.spanId = trace.ParseCurrentContext(ctx)
	if b.traceId == "" {
		b.traceId = trace.GetRequestId(ctx)
	}

	return b
}

func (b *baseLogger) Extra(e *Extras) Log {
	b.extra = e
	return b
}

func (b *baseLogger) TraceId(s string) Log {
	b.traceId = s
	return b
}

func (b *baseLogger) SpanId(s string) Log {
	b.spanId = s
	return b
}

func (b *baseLogger) SetErr(err error) Log {
	b.err = err
	return b
}

// SetExtra
//
//	@Description: 设置日志额外字段信息
//	@receiver b
//	@param e  body any true "-"
//	@return Log
func (b *baseLogger) SetExtra(e *Extras) Log {
	b.extra = e
	return b
}

func (b *baseLogger) writeMsg(when time.Time, msg string, logLevel LogLevel, extra []Field) (err error) {
	if reqId := trace.GetRequestId(b.ctx); 0 < len(reqId) {
		msg = reqId + ": " + msg
	}
	var (
		src        = ""
		file       = ""
		lineno int = 0
		ok         = false
		strim      = "src/"
	)
	if logLevel <= b.ShowCallerLevel {
		_, file, lineno, ok = runtime.Caller(b.callDepth)
		if ok {
			file = toShortCaller(file)
		}
	}
	if b.usePath != "" {
		strim = b.usePath
	}
	if ok {
		src = strings.Replace(fmt.Sprintf("%s:%d", stringTrim(file, strim), lineno), "%2e", ".", -1)
	}

	//l.writeToLoggers(when, msgSt, logLevel)
	msgStr := when.Format(b.timeFormat) + " [" + levelPrefix[logLevel] + "] " + "[" + src + "] " + msg
	//ad, err := GetAdapter(l.adapterName)
	//if err != nil {
	//	return err
	//}
	err = b.adapter.LogWrite(when, msgStr, logLevel, extra)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to WriteMsg to adapter:%v,error:%v\n", b.output.name, err)
	}

	return nil
}

// writeMsg 函数式
func writeMsg(ctx context.Context, adapter Driver, usePath, timeFormat string, when time.Time, msg string,
	showCallerLevel, logLevel LogLevel, callDepth int, extra []Field) (err error) {
	if reqId := trace.GetRequestId(ctx); 0 < len(reqId) {
		msg = reqId + ": " + msg
	}
	var (
		src        = ""
		file       = ""
		lineno int = 0
		ok         = false
		strim      = "src/"
	)
	if logLevel <= showCallerLevel {
		_, file, lineno, ok = runtime.Caller(callDepth)
		if ok {
			file = toShortCaller(file)
		}
	}
	if usePath != "" {
		strim = usePath
	}
	if ok {
		src = strings.Replace(fmt.Sprintf("%s:%d", stringTrim(file, strim), lineno), "%2e", ".", -1)
	}

	msgStr := when.Format(timeFormat) + " [" + levelPrefix[logLevel] + "] " + "[" + src + "] " + msg
	err = adapter.LogWrite(when, msgStr, logLevel, extra)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to WriteMsg to adapter:%v,error:%v\n", adapter.Name(), err)
	}

	return nil
}

func (b *baseLogger) doLog(logLevel LogLevel, format string, args ...interface{}) {
	//ad, err := GetAdapter(b.adapterName)
	//if err != nil {
	//	fmt.Printf("get adapter error %s %v \n", b.adapterName, adapters)
	//	panic(err)
	//}
	if logLevel > b.adapter.CurrentLevel() { // get log level
		return
	}
	var extra = b.assembleField(b.extra)
	var msg = formatLog(format, args...)
	var t = time.Now()
	extra.Int64("ms", t.UnixNano()/1e6)
	if b.adapterName == AdapterZap {
		_ = b.adapter.LogWrite(t, msg, logLevel, extra.fields)
	} else {
		//b.writeMsg(t, msg, logLevel, extra.fields)
		writeMsg(b.ctx, b.adapter, b.usePath, b.timeFormat, t, msg, b.ShowCallerLevel, logLevel, b.callDepth, extra.fields)
	}

	// call back 机器人
	if b.notify && b.callBack != nil {
		b.callBack.LogWrite(t, msg, logLevel, extra.fields)
	}
	// 输出日志后，值清空
	b.ctx = nil
	b.spanId = ""
	b.traceId = ""
	b.err = nil
	b.notify = false
	b.exLock.Lock()
	if b.extra != nil {
		b.extra = nil
	}
	b.exLock.Unlock()
}

func doLog(ctx context.Context, adapter, callBack Driver, usePath, spanId, traceId, timeFormat, format string,
	showCallerLevel, logLevel LogLevel, callDepth int, notify bool, err error, f []Field, args ...interface{}) {
	if logLevel > adapter.CurrentLevel() {
		return
	}
	var extra = setExtraField(f, err, spanId, traceId)
	var msg = formatLog(format, args...)
	var t = time.Now()
	extra.Int64("ms", t.UnixNano()/1e6)
	if adapter.Name() == AdapterZap {
		_ = adapter.LogWrite(t, msg, logLevel, extra.fields)
	} else {
		writeMsg(ctx, adapter, usePath, timeFormat, t, msg, showCallerLevel, logLevel, callDepth, extra.fields)
	}

	// call back 机器人
	if notify && callBack != nil {
		callBack.LogWrite(t, msg, logLevel, extra.fields)
	}
}

// Panic
//
//	@Description: 记录 panic 信息。注意：DPanic在非development下则退化为error模式，最后的info照样会输出，这样子在production下比较安全一点。
//	@receiver b
//	@param format  body any true "-"
//	@param args  body any true "-"
func (b *baseLogger) Panic(format string, args ...interface{}) {
	//b.doLog(LevelEmergency, format, args...)
	doLog(b.ctx, b.adapter, b.callBack, b.usePath, b.spanId, b.traceId, b.timeFormat, format,
		b.ShowCallerLevel, LevelEmergency, b.callDepth, b.notify, b.err, b.extra.fields, args...)
}

func (b *baseLogger) Emer(format string, args ...interface{}) {
	//b.doLog(LevelEmergency, format, args...)
	doLog(b.ctx, b.adapter, b.callBack, b.usePath, b.spanId, b.traceId, b.timeFormat, format,
		b.ShowCallerLevel, LevelEmergency, b.callDepth, b.notify, b.err, b.extra.fields, args...)
}

func (b *baseLogger) Alert(format string, args ...interface{}) {
	//b.doLog(LevelAlert, format, args...)
	doLog(b.ctx, b.adapter, b.callBack, b.usePath, b.spanId, b.traceId, b.timeFormat, format,
		b.ShowCallerLevel, LevelAlert, b.callDepth, b.notify, b.err, b.extra.fields, args...)
}

func (b *baseLogger) Crit(format string, args ...interface{}) {
	//b.doLog(LevelCritical, format, args...)
	doLog(b.ctx, b.adapter, b.callBack, b.usePath, b.spanId, b.traceId, b.timeFormat, format,
		b.ShowCallerLevel, LevelCritical, b.callDepth, b.notify, b.err, b.extra.fields, args...)
}

func (b *baseLogger) Error(format string, args ...interface{}) {
	//b.doLog(LevelError, format, args...)
	doLog(b.ctx, b.adapter, b.callBack, b.usePath, b.spanId, b.traceId, b.timeFormat, format,
		b.ShowCallerLevel, LevelError, b.callDepth, b.notify, b.err, b.extra.fields, args...)
}

func (b *baseLogger) Warn(format string, args ...interface{}) {
	//b.doLog(LevelWarning, format, args...)
	doLog(b.ctx, b.adapter, b.callBack, b.usePath, b.spanId, b.traceId, b.timeFormat, format,
		b.ShowCallerLevel, LevelWarning, b.callDepth, b.notify, b.err, b.extra.fields, args...)
}

func (b *baseLogger) Info(format string, args ...interface{}) {
	//b.doLog(LevelInformational, format, args...)
	doLog(b.ctx, b.adapter, b.callBack, b.usePath, b.spanId, b.traceId, b.timeFormat, format,
		b.ShowCallerLevel, LevelInformational, b.callDepth, b.notify, b.err, b.extra.fields, args...)
}

func (b *baseLogger) Debug(format string, args ...interface{}) {
	//b.doLog(LevelDebug, format, args...)
	doLog(b.ctx, b.adapter, b.callBack, b.usePath, b.spanId, b.traceId, b.timeFormat, format,
		b.ShowCallerLevel, LevelDebug, b.callDepth, b.notify, b.err, b.extra.fields, args...)
}

func (b *baseLogger) Trace(format string, args ...interface{}) {
	//b.doLog(LevelTrace, format, args...)
	doLog(b.ctx, b.adapter, b.callBack, b.usePath, b.spanId, b.traceId, b.timeFormat, format,
		b.ShowCallerLevel, LevelTrace, b.callDepth, b.notify, b.err, b.extra.fields, args...)
}

func (b *baseLogger) assembleField(extra *Extras) (e *Extras) {
	b.exLock.Lock()
	e = E()
	if extra != nil && extra.fields != nil {
		e = &Extras{fields: make([]Field, len(b.extra.fields))}
		copy(e.fields, extra.fields)
	}
	b.exLock.Unlock()

	if b.err != nil {
		e.Error(b.err)
	}

	//e.Ctx(c.ctx)
	e.Span(b.spanId)
	e.Trace(b.traceId)
	return e
}

func setExtraField(f []Field, err error, spanId, traceId string) (e *Extras) {
	e = &Extras{fields: f}

	if err != nil {
		e.Error(err)
	}

	//e.Ctx(c.ctx)
	e.Span(spanId)
	e.Trace(traceId)
	return e
}
