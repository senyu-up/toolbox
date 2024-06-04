package logger

import (
	"context"
)

// LocalLogger
// @Description: 这个文件，是为了把 baseLogger 再包一层，避免 runtime.caller 返回不一样问题
type LocalLogger struct {
	bl *baseLogger
}

func NewLocalLogger(name string, d Driver, opts ...LogOption) Log {
	return &LocalLogger{bl: newBaseLogger(name, d, opts...)}
}

func (l *LocalLogger) Clone() Log {
	return &LocalLogger{bl: l.bl.Clone().(*baseLogger)}
}

func (l *LocalLogger) AdapterName() string {
	return l.bl.AdapterName()
}

func (l *LocalLogger) SetCallBack(cb Driver) Log {
	return l.bl.SetCallBack(cb)
}

func (l *LocalLogger) Notify() Log {
	return l.bl.Clone().Notify()
}

func (l *LocalLogger) Ctx(c context.Context) Log {
	return l.bl.Clone().Ctx(c)
}

func (l *LocalLogger) TraceId(t string) Log {
	return l.bl.Clone().TraceId(t)
}

func (l *LocalLogger) SetExtra(e *Extras) Log {
	return l.bl.Clone().SetExtra(e)
}

func (l *LocalLogger) SpanId(s string) Log {
	return l.bl.Clone().SpanId(s)
}

func (l *LocalLogger) SetErr(e error) Log {
	return l.bl.Clone().SetErr(e)
}

func (l *LocalLogger) Panic(format string, args ...interface{}) {
	//l.bl.doLog(LevelEmergency, format, args...)
	doLog(l.bl.ctx, l.bl.adapter, l.bl.callBack, l.bl.usePath, l.bl.spanId, l.bl.traceId, l.bl.timeFormat, format,
		l.bl.ShowCallerLevel, LevelEmergency, l.bl.callDepth, l.bl.notify, l.bl.err, l.bl.extra.fields, args...)
}

func (l *LocalLogger) Emer(format string, args ...interface{}) {
	//l.bl.doLog(LevelEmergency, format, args...)
	doLog(l.bl.ctx, l.bl.adapter, l.bl.callBack, l.bl.usePath, l.bl.spanId, l.bl.traceId, l.bl.timeFormat, format,
		l.bl.ShowCallerLevel, LevelEmergency, l.bl.callDepth, l.bl.notify, l.bl.err, l.bl.extra.fields, args...)
}

func (l *LocalLogger) Alert(format string, args ...interface{}) {
	//l.bl.doLog(LevelAlert, format, args...)
	doLog(l.bl.ctx, l.bl.adapter, l.bl.callBack, l.bl.usePath, l.bl.spanId, l.bl.traceId, l.bl.timeFormat, format,
		l.bl.ShowCallerLevel, LevelAlert, l.bl.callDepth, l.bl.notify, l.bl.err, l.bl.extra.fields, args...)
}

func (l *LocalLogger) Crit(format string, args ...interface{}) {
	//l.bl.doLog(LevelCritical, format, args...)
	doLog(l.bl.ctx, l.bl.adapter, l.bl.callBack, l.bl.usePath, l.bl.spanId, l.bl.traceId, l.bl.timeFormat, format,
		l.bl.ShowCallerLevel, LevelCritical, l.bl.callDepth, l.bl.notify, l.bl.err, l.bl.extra.fields, args...)
}

func (l *LocalLogger) Error(format string, args ...interface{}) {
	//l.bl.doLog(LevelError, format, args...)
	doLog(l.bl.ctx, l.bl.adapter, l.bl.callBack, l.bl.usePath, l.bl.spanId, l.bl.traceId, l.bl.timeFormat, format,
		l.bl.ShowCallerLevel, LevelError, l.bl.callDepth, l.bl.notify, l.bl.err, l.bl.extra.fields, args...)
}

func (l *LocalLogger) Warn(format string, args ...interface{}) {
	//l.bl.doLog(LevelWarning, format, args...)
	doLog(l.bl.ctx, l.bl.adapter, l.bl.callBack, l.bl.usePath, l.bl.spanId, l.bl.traceId, l.bl.timeFormat, format,
		l.bl.ShowCallerLevel, LevelWarning, l.bl.callDepth, l.bl.notify, l.bl.err, l.bl.extra.fields, args...)
}

func (l *LocalLogger) Info(format string, args ...interface{}) {
	//l.bl.doLog(LevelInformational, format, args...)
	doLog(l.bl.ctx, l.bl.adapter, l.bl.callBack, l.bl.usePath, l.bl.spanId, l.bl.traceId, l.bl.timeFormat, format,
		l.bl.ShowCallerLevel, LevelInformational, l.bl.callDepth, l.bl.notify, l.bl.err, l.bl.extra.fields, args...)
}

func (l *LocalLogger) Debug(format string, args ...interface{}) {
	//l.bl.doLog(LevelDebug, format, args...)
	doLog(l.bl.ctx, l.bl.adapter, l.bl.callBack, l.bl.usePath, l.bl.spanId, l.bl.traceId, l.bl.timeFormat, format,
		l.bl.ShowCallerLevel, LevelDebug, l.bl.callDepth, l.bl.notify, l.bl.err, l.bl.extra.fields, args...)
}

func (l *LocalLogger) Trace(format string, args ...interface{}) {
	//l.bl.doLog(LevelTrace, format, args...)
	doLog(l.bl.ctx, l.bl.adapter, l.bl.callBack, l.bl.usePath, l.bl.spanId, l.bl.traceId, l.bl.timeFormat, format,
		l.bl.ShowCallerLevel, LevelTrace, l.bl.callDepth, l.bl.notify, l.bl.err, l.bl.extra.fields, args...)
}
