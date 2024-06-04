# 日志

日志组件调整了初始化的方式（由toolbox/tool/config/logger.go 里的结构体），以及砍掉了冗余接口，把微信机器人也挪出去了。

## 基本使用

下面是代码是：

1.  通过文件加载配置。
2.  配置序列化到 config.LogConfig 结构体。
3.  通过Log配置初始化 logger。
4.  使用logger

```go
package example

import (
	"context"
	"toolbox/tool/config"
	"toolbox/tool/config/loader"
	"toolbox/tool/logger"
	"toolbox/utils/trace"
)

func main() {
	// 通过文件配置
	conf, err := loader.InitConf(&loader.File{},
		loader.ConfOptWithPath("./config.yaml"),
		loader.ConfOptWithType("yaml"))
	if err != nil {
		return
	}
	
	// 配置解析道 config.LogConfig
	type AppConf struct {
		Logger config.LogConfig
	}
	var app = AppConf{}
	if err = conf.Unmarshal(&app); err != nil {
		return
	}

	// 通过配置初始化 logger
	logger.InitLoggerByConf(app.Logger)

	// 指定 zap logger
	logger.SwitchLogger(logger.AdapterZap)
	
	// 为日志设置 trace 信息
	var traceCtx = trace.NewTrace()
	logger.Ctx(traceCtx).Info("print log with new trace ctx")

	// 手动设置 trace 信息
	var newCtx = trace.NewContextWithRequestIdAndSpanId(context.TODO(), "req_111", "span_222")
	logger.Ctx(newCtx).Info("print log with set trace,span ctx")
}

```

## 各等级日志打印

各种等级的日志输出：

```go
import (
	"toolbox/tool/logger"
)
func LogAllLevel() {
	logger.Trace("this is Trace")
	logger.Debug("this is Debug")
	logger.Info("this is Info")
	logger.Warn("this is Warn")
	logger.Error("this is Error")
	logger.Crit("this is Critical")
	logger.Alert("this is Alert")
	logger.Emer("this is Emergency")
}

```

## 日志额外参数

设置日志额外的参数：context，trace，error，zap.Field.

```go
func LogExtraFields() {
	var ctx = trace.NewTrace()
	logger.Ctx(ctx).Info("log ctx")

	traceId, spanId := trace.NewRequestId(), trace.NewSpanId()
	logger.TraceId(traceId).SpanId(spanId).Trace("log traceId, spanId")

	var err = errors.New("This is err ")
	logger.SetErr(err).Error("I got a err", err)

	logger.SetExtra(logger.E().String("name", "xh").String("level", "2")).Info("log extra fields")
}

```

## 日志&企微机器人

如何设置在输出日志同时，能发给企微机器人？

这里初始化企微机器人，然后设置到全局 logger 的 callback 中。当在输出日志前调用 `Notify()` 就会在输出日志时在发一条消息给机器人

```go
import (
	"context"
	"errors"
	"fmt"
	"toolbox/tool/config"
	"toolbox/tool/config/loader"
	"toolbox/tool/logger"
	"toolbox/tool/wework/qwrobot"
	"toolbox/utils/trace"
)

func SetLoggerCallBack() {
	// 设置 qw roboto driver
	qwr := qwrobot.New(config.QwRobotConfig{})
	// 设置 callback 为 qw robot
	logger.SetCallBack(logger.NewQwRobot(qwr))
	// 一般的日志
	logger.Error("roboto test")
	// 调用 notify 后，会触发 callback，也就是 qw robot
	logger.Notify().Error("roboto test")
}

```

## 日志拓展规范

toolbox其他组件，也有日志输出的需求，这个时候是需要外部制定logger的，但大部分情况是组件内声明一个，这意味着无法从外部控制该日志组件的行为。所以，其他组件如果有 日志记录需求。

所以logger抽象出了接口，其他组件有日志输出时，调用日志抽象接口即可。