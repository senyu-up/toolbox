package su_logger

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/trace"
	"github.com/senyu-up/toolbox/tool/wework/qwrobot"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewExtra() *Extra {
	return E()
}

func E() *Extra {
	return &Extra{
		fields: make([]zap.Field, 0),
	}
}

type Extra struct {
	fields []zap.Field
}

func (e *Extra) String(key string, value string) *Extra {
	e.fields = append(e.fields, zap.String(key, value))
	return e
}

func (e *Extra) Int(key string, v int) *Extra {
	e.fields = append(e.fields, zap.Int(key, v))

	return e
}

func (e *Extra) Int8(key string, v int8) *Extra {
	e.fields = append(e.fields, zap.Int8(key, v))

	return e
}

func (e *Extra) Int16(key string, v int16) *Extra {
	e.fields = append(e.fields, zap.Int16(key, v))

	return e
}

func (e *Extra) Int32(key string, v int32) *Extra {
	e.fields = append(e.fields, zap.Int32(key, v))

	return e
}

func (e *Extra) Int64(key string, v int64) *Extra {
	e.fields = append(e.fields, zap.Int64(key, v))

	return e
}

func (e *Extra) Uint(key string, v uint) *Extra {
	e.fields = append(e.fields, zap.Uint(key, v))

	return e
}

func (e *Extra) Uint8(key string, v uint8) *Extra {
	e.fields = append(e.fields, zap.Uint8(key, v))

	return e
}

func (e *Extra) Uint16(key string, v uint16) *Extra {
	e.fields = append(e.fields, zap.Uint16(key, v))

	return e
}

func (e *Extra) Uint32(key string, v uint32) *Extra {
	e.fields = append(e.fields, zap.Uint32(key, v))

	return e
}

func (e *Extra) Uint64(key string, v uint64) *Extra {
	e.fields = append(e.fields, zap.Uint64(key, v))

	return e
}

func (e *Extra) Float64(key string, v float64) *Extra {
	e.fields = append(e.fields, zap.Float64(key, v))

	return e
}

func (e *Extra) Bool(key string, v bool) *Extra {
	e.fields = append(e.fields, zap.Bool(key, v))

	return e
}

func (e *Extra) Any(key string, value interface{}) *Extra {
	e.Interface(key, value)

	return e
}

func (e *Extra) Error(err error) *Extra {
	e.fields = append(e.fields, zap.Error(err))

	return e
}

func (e *Extra) NamedError(key string, err error) *Extra {
	e.fields = append(e.fields, zap.NamedError(key, err))

	return e
}

func (e *Extra) Interface(key string, v interface{}) *Extra {
	if v == nil {
		e.fields = append(e.fields, zap.String(key, ""))
	} else {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice, reflect.Array, reflect.Struct, reflect.Map, reflect.Ptr:
			// 对结构体, map, slice 进行特判
			s, err1 := jsoniter.MarshalToString(v)
			if err1 != nil {
				e.fields = append(e.fields, zap.Any(key, v))
			} else {
				e.fields = append(e.fields, zap.String(key, s))
			}
		default:
			e.fields = append(e.fields, zap.Any(key, v))
		}
	}

	return e
}

func (e *Extra) Ctx(ctx context.Context) *Extra {
	if curC, ok := ctx.(C); ok {
		e.Trace(curC.traceId)
		e.Span(curC.spanId)
	} else {
		traceId, spanId := trace.ParseFromContext(ctx)
		e.Trace(traceId)
		e.Span(spanId)
	}

	return e
}

func (e *Extra) Trace(traceId string) *Extra {
	if traceId != "" {
		e.fields = append(e.fields, zap.String("trace", traceId))
	}

	return e
}

func (e *Extra) NowMs() *Extra {
	e.fields = append(e.fields, zap.Int64("ms", time.Now().UnixNano()/1e6))

	return e
}

func (e *Extra) Span(spanId string) *Extra {
	if spanId != "" {
		e.fields = append(e.fields, zap.String("span", spanId))
	}

	return e
}

var logLevel zapcore.Level = -2

func isEnable(level zapcore.Level) bool {
	if logLevel == -2 {
		// 获取当前的stage, 并设置对应的日志级别
		stage := os.Getenv(enum.StageKey)
		if stage == enum.EvnStageProduction || stage == "master" {
			logLevel = zapcore.InfoLevel
		} else {
			logLevel = zapcore.DebugLevel
		}
	}

	if level < logLevel {
		return false
	}

	return true
}

func assembleField(ctx context.Context, err error, extra []*Extra) []zap.Field {
	var e *Extra
	if len(extra) == 0 || extra[0] == nil {
		e = E()
	} else {
		e = extra[0]
	}

	if err != nil {
		e.Error(err)
	}

	e.Ctx(ctx)

	return e.fields
}

func WithCaller() *Extra {
	return nil
}

func Info(ctx context.Context, msg string, extra ...*Extra) {
	if isEnable(zapcore.InfoLevel) {
		if len(extra) > 1 {
			loggerWithCaller.Info(msg, assembleField(ctx, nil, extra)...)
		} else {
			loggerWithoutCaller.Info(msg, assembleField(ctx, nil, extra)...)
		}
	}
}

func InfoWithNotify(ctx context.Context, msg string, extra ...*Extra) {
	fields := assembleField(ctx, nil, extra)
	loggerWithCaller.Info(msg, fields...)
	Notify(zapcore.InfoLevel, msg, fields)
}

func Debug(ctx context.Context, msg string, extra ...*Extra) {
	if isEnable(zapcore.DebugLevel) {
		if len(extra) > 1 {
			loggerWithCaller.Debug(msg, assembleField(ctx, nil, extra)...)
		} else {
			loggerWithoutCaller.Debug(msg, assembleField(ctx, nil, extra)...)
		}
	}
}

func Warn(ctx context.Context, msg string, extra ...*Extra) {
	if isEnable(zapcore.WarnLevel) {
		loggerWithCaller.Warn(msg, assembleField(ctx, nil, extra)...)
	}
}

func WarnWithNotify(ctx context.Context, msg string, extra ...*Extra) {
	fields := assembleField(ctx, nil, extra)
	loggerWithCaller.Warn(msg, fields...)
	Notify(zapcore.WarnLevel, msg, fields)
}

func Notify(level zapcore.Level, msg string, fields []zap.Field) {
	robot := qwrobot.Get()
	if robot != nil {
		content := fieldsToStringV2(fields)
		qwMsg := qwrobot.Message{
			Title:   msg,
			Content: content,
		}
		if level == zapcore.InfoLevel {
			robot.Info(qwMsg)
		} else if level == zapcore.WarnLevel {
			robot.Warn(qwMsg)
		} else if level == zapcore.ErrorLevel || level == zapcore.PanicLevel || level == zapcore.DPanicLevel || level == zapcore.FatalLevel {
			robot.Error(qwMsg)
		}
	}
}

func Error(ctx context.Context, err error, msg string, extra ...*Extra) {
	if isEnable(zapcore.ErrorLevel) {
		loggerWithCaller.Error(msg, assembleField(ctx, err, extra)...)
	}
}

func ErrorWithNotify(ctx context.Context, err error, msg string, extra ...*Extra) {
	fields := assembleField(ctx, err, extra)
	loggerWithCaller.Error(msg, fields...)
	Notify(zapcore.ErrorLevel, msg, fields)
}

func Panic(ctx context.Context, msg string, extra ...*Extra) {
	stack := FormatStackInline(debug.Stack(), 5)
	err := errors.New(stack)

	loggerWithCaller.DPanic(msg, assembleField(ctx, err, extra)...)
}

func PanicWithNotify(ctx context.Context, msg string, extra ...*Extra) {
	stack := FormatStackInline(debug.Stack(), 5)
	err := errors.New(stack)
	fields := assembleField(ctx, err, extra)
	loggerWithCaller.DPanic(msg, fields...)

	Notify(zapcore.DPanicLevel, msg, fields)
}

// FormatStackInline
// 调用栈信息, 进行压缩, 否则在日志平台显示出来内容过多
func FormatStackInline(stack []byte, jumpLevel int) string {
	if stack == nil {
		return ""
	}
	r := bytes.NewReader(stack)
	reader := bufio.NewReader(r)
	s := strings.Builder{}
	sep := []byte{'/'}
	var i int
	for {
		line, _, err := reader.ReadLine()
		if i >= jumpLevel {
			if err == nil && s.Len() > 0 {
				s.WriteString("->")
			}

			cnt := bytes.Count(line, sep)
			if cnt > 3 {
				byteSlice := bytes.Split(line, sep)
				sliceLen := len(byteSlice)
				fromPos := sliceLen - 3
				s.Write(bytes.Join(byteSlice[fromPos:], []byte{'/'}))
			} else {
				s.Write(line)
			}
		}
		i++

		if err == io.EOF {
			break
		}
	}

	return s.String()
}

func fieldsToStringV2(fields []zap.Field) string {
	s := strings.Builder{}
	// 事件点
	fields = append(fields, zap.Field{
		Key:    "ts",
		Type:   zapcore.StringType,
		String: time.Now().Format("2006-01-02 15:04:05.000"),
	})
	// 触发代码位置
	_, file, line, ok := runtime.Caller(3)
	if ok {
		caller := toShortCaller(file)
		fields = append(fields, zap.Field{
			Key:    "caller",
			Type:   zapcore.StringType,
			String: caller + ":" + cast.ToString(line),
		})
	}

	for i, _ := range fields {
		if s.Len() > 0 {
			s.WriteByte('\n')
		}
		s.WriteString("> ")
		s.WriteString(fields[i].Key)
		s.WriteString(": ")
		switch fields[i].Type {
		case zapcore.StringType:
			s.WriteString(fields[i].String)
		case zapcore.Int8Type, zapcore.Int16Type, zapcore.Int32Type, zapcore.Int64Type, zapcore.Uint8Type, zapcore.Uint16Type, zapcore.Uint32Type, zapcore.Uint64Type, zapcore.BoolType, zapcore.Float32Type, zapcore.Float64Type:
			s.WriteString(cast.ToString(fields[i].Integer))
		default:
			s.WriteString(fmt.Sprintf("%v", fields[i].Interface))
		}
	}

	return s.String()
}

func toShortCaller(line string) string {
	list := strings.Split(line, "/")
	if len(list) > 1 {
		max := len(list)
		return strings.Join(list[max-2:max], "/")
	} else {
		return strings.Join(list, "/")
	}
}

type C struct {
	context.Context
	traceId string
	spanId  string
}

func (c C) TraceId() string {
	return c.traceId
}
func (c C) SpanId() string {
	return c.spanId
}

func (c C) Trace() (traceId, spanId string) {
	return c.traceId, c.spanId
}

func Ctx(ctx context.Context) context.Context {
	traceId, spanId := trace.ParseFromContext(ctx)

	return C{
		Context: nil,
		traceId: traceId,
		spanId:  spanId,
	}
}

// 尝试从当前ctx中获取链路, 获取失败则自动填充链路信息, 用于 脚本, 定时任务, 没有上下文的场景
func AutoCtx(ctx context.Context) context.Context {
	traceId, spanId := trace.ParseFromContext(ctx)
	if traceId == "" {
		traceId = trace.NewTraceID()
	}
	if spanId == "" {
		spanId = trace.NewSpanID()
	}

	return C{
		Context: nil,
		traceId: traceId,
		spanId:  spanId,
	}
}
