package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/trace"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

var (
	infoStr      = "%s"
	warnStr      = "%s"
	errStr       = "%s"
	traceStr     = "[%.3fms] [rows:%v] %s"
	traceWarnStr = "%s [%.3fms] [rows:%v] %s"
	traceErrStr  = "%s [%.3fms] [rows:%v] %s"
)

type GormLogger struct {
	logDriver                 logger.Log // logger 接口
	level                     glogger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
	traceOn                   bool
}

func NewGormLogger(opts ...SqlLogOption) *GormLogger {
	var ml = &GormLogger{
		level:         glogger.Info,
		slowThreshold: time.Second * 5,
		logDriver:     logger.GetLogger(),
	}
	// 先用默认值初始化，下面应用用户穿参数
	for _, opt := range opts {
		opt(ml)
	}
	return ml
}

func (m *GormLogger) LogMode(level glogger.LogLevel) glogger.Interface {
	m.level = level

	return m
}

func (m *GormLogger) checkLevel(level glogger.LogLevel) bool {
	return level <= m.level
}

func (m *GormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	if m.checkLevel(glogger.Info) {
		msg := m.Printf(ctx, infoStr+s, i...)
		m.logDriver.Ctx(ctx).Info(msg)
	}
}

func (m *GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	if m.checkLevel(glogger.Warn) {
		msg := m.Printf(ctx, warnStr+s, i...)
		m.logDriver.Ctx(ctx).Warn(msg)
	}
}

func (m *GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	//TODO implement me
	if m.checkLevel(glogger.Error) {
		msg := m.Printf(ctx, errStr+s, i...)
		m.logDriver.Ctx(ctx).Error(msg)
	}
}

func (m *GormLogger) Printf(ctx context.Context, s string, i ...interface{}) string {
	traceId, _ := trace.ParseCurrentContext(ctx)
	if traceId == "" {
		traceId = "-"
	}
	return fmt.Sprintf("[trace:"+traceId+"] "+s, i...)
}

func (m *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if m.level <= glogger.Silent {
		return
	}
	var isSlow bool
	var sql string = ""
	var rows int64 = 0
	// 开启链路追踪, 且 traceId 不为空，才记录
	if m.traceOn && len(trace.ParseRequestId(ctx)) > 0 {
		var traceId, spanId, newSpanId string
		ctx, traceId, spanId, newSpanId = trace.ParseOrGenContext(ctx)
		span := trace.NewJaegerSpan("mysql", traceId, newSpanId, spanId, nil, nil)
		defer func() {
			span.SetTag("sql", sql)
			if isSlow {
				span.SetTag("slowsql", 1)
			}
			span.Finish()
		}()
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && m.level >= glogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !m.ignoreRecordNotFoundError):
		sql, rows = fc()
		var msg string
		if rows == -1 {
			msg = m.Printf(ctx, traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			msg = m.Printf(ctx, traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
		m.logDriver.Ctx(ctx).Error(msg)
	case elapsed > m.slowThreshold && m.slowThreshold != 0 && m.level >= glogger.Warn:
		sql, rows = fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", m.slowThreshold)
		isSlow = true
		var msg string
		if rows == -1 {
			msg = m.Printf(ctx, traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			msg = m.Printf(ctx, traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}

		m.logDriver.Ctx(ctx).Warn(msg)
	case m.level == glogger.Info:
		sql, rows = fc()
		var msg string
		if rows == -1 {
			// utils.FileWithLineNum()
			msg = m.Printf(ctx, traceStr, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			msg = m.Printf(ctx, traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
		m.logDriver.Ctx(ctx).Info(msg)
	}
}
