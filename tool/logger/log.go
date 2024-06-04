package logger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

// 日志等级，从0-7，日优先级由高到低
type LogLevel int

const (
	LevelEmergency     LogLevel = iota // 系统级紧急，比如磁盘出错，内存异常，网络不可用等
	LevelAlert                         // 系统级警告，比如数据库访问异常，配置文件出错等
	LevelCritical                      // 系统级危险，比如权限出错，访问异常等
	LevelError                         // 用户级错误
	LevelWarning                       // 用户级警告
	LevelInformational                 // 用户级信息
	LevelDebug                         // 用户级调试
	LevelTrace                         // 用户级基本输出
)

const DefaultCallDepth = 3 // call stack 获取时，跳过多少层级

// LevelMap 日志等级和描述映射关系
var LevelMap = map[string]LogLevel{
	"EMER": LevelEmergency,
	"ALRT": LevelAlert,
	"CRIT": LevelCritical,
	"EROR": LevelError,
	"WARN": LevelWarning,
	"INFO": LevelInformational,
	"DEBG": LevelDebug,
	"TRAC": LevelTrace,
	"":     LevelInformational,
}

// 注册实现的适配器， 当前支持控制台，文件
var adapters = make(map[string]Driver)
var adapterRwLock = sync.RWMutex{}
var ErrAdapterNotFound = fmt.Errorf("The logger driver was not found by adapter name ")

// 日志记录等级字段
var levelPrefix = [LevelTrace + 1]string{
	"EMER",
	"ALRT",
	"CRIT",
	"EROR",
	"WARN",
	"INFO",
	"DEBG",
	"TRAC",
}

var ErrInvalidLogLevel = errors.New("invalid log level")

const (
	LogTimeDefaultFormat = "2006-01-02 15:04:05" // 日志输出默认格式
	AdapterConsole       = "console"             // 控制台输出配置项
	AdapterFile          = "file"                // 文件输出配置项a
	AdapterConn          = "conn"
	AdapterZap           = "zap"
)

var appSn = os.Getenv("APPSN")

const produce = "PRODUCE"

var env = os.Getenv("ENV")

// log 接口
type Log interface {
	Ctx(ctx context.Context) Log

	Clone() Log

	TraceId(string) Log

	SpanId(string) Log

	SetErr(error) Log

	SetCallBack(Driver) Log

	Notify() Log

	//
	// SetExtra
	//  @Description: 设置额外参数用
	//  @return Log
	//
	SetExtra(extra *Extras) Log

	Panic(format string, args ...interface{})

	// 这里的告警登记，依次降低
	// 系统级别紧急
	Emer(format string, args ...interface{})

	Alert(format string, args ...interface{})

	Crit(format string, args ...interface{})

	Error(format string, args ...interface{})

	Warn(format string, args ...interface{})

	Info(format string, args ...interface{})

	Debug(format string, args ...interface{})

	Trace(format string, args ...interface{})
}

type loginfo struct {
	Time    string
	Level   string
	Path    string
	Name    string
	Content string
}

type nameLogger struct {
	name   string
	config interface{}
}

// register
// @description 日志输出适配器注册，log需要实现Init，LogWrite，Destroy方法
func register(name string, log Driver) {
	if log == nil {
		panic("logs: register provide is nil")
	}
	if _, ok := adapters[name]; ok {
		//panic("logs: register called twice for provider " + name)
	}
	adapterRwLock.Lock()
	defer adapterRwLock.Unlock()
	adapters[name] = log
}

// GetAdapter
// @description 日志输出适配器注册，log需要实现Init，LogWrite，Destroy方法
func GetAdapter(name string) (Driver, error) {
	adapterRwLock.RLock()
	defer adapterRwLock.RUnlock()
	if val, ok := adapters[name]; ok {
		return val, nil
	} else {
		return nil, ErrAdapterNotFound
	}
}

func CheckFormats(format string) int {
	if 2 > len(format) {
		return 0
	}
	var count int
	for i := 0; i < len(format); i++ {
		if format[i] == '%' {
			if i+1 < len(format) && format[i+1] != '%' {
				count++
			}
			i++
		}
	}
	return count
}

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		var vCount = CheckFormats(msg)
		if vCount >= len(v) {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v)-vCount)
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	if 1 > len(v) {
		return msg
	}
	return fmt.Sprintf(msg, v...)
}

func Level() LogLevel {
	var lv = LevelTrace
	if env == produce {
		lv = LevelInformational
	}

	return lv
}

func stringTrim(s string, cut string) string {
	ss := strings.SplitN(s, cut, 2)
	if 1 == len(ss) {
		return ss[0]
	}
	return ss[1]
}
