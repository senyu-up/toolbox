package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/config/loader"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/trace"
	"github.com/senyu-up/toolbox/tool/wework/qwrobot"
)

func sample() {
	logger.Info("xh logger ")
	var err = errors.New("Boom! ")
	logger.SetErr(err).Error("You got err ")
}

func InitLoggerByConfExample() {
	// 通过文件配置
	conf, err := loader.InitConf(&loader.File{},
		loader.ConfOptWithPath("./config.yaml"),
		loader.ConfOptWithType("yaml"))
	if err != nil {
		fmt.Printf("load config err %v\n", err)
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
	if _, err = logger.InitDefaultLoggerByConf(&app.Logger); err != nil {
		fmt.Printf("init log by config err %v\n", err)
		return
	}

	logger.GetLogger()

	// 指定 zap logger
	logger.SwitchLogger(logger.AdapterZap) // qw机器人也作为 adatper 中的一个，

	// 在 lloger 加入 callback ，传入 气味机器人 的发送代码
	logger.Debug("Ok~")
}

func LogExtraInfo() {
	logger.SwitchLogger(logger.AdapterFile)
	logger.SwitchLogger(logger.AdapterConsole)
	// 错误信息
	var err = errors.New("Boom! ")
	logger.SetErr(err).Error("You got err ")

	// 为日志设置 trace 信息
	var traceCtx = trace.NewTrace()
	logger.Ctx(traceCtx).Info("print log with new trace ctx")

	// 手动设置 trace 信息
	var newCtx = trace.NewContextWithRequestIdAndSpanId(context.TODO(), "req_111", "span_222")
	logger.Ctx(newCtx).Info("print log with set trace,span ctx")

	logger.SpanId("span_222").TraceId("req_111").Info("Set TraceId, spanId by handle")
}

func SetLoggerCallBack() {
	// 设置 qw roboto driver
	qwr := qwrobot.New(&config.QwRobotConfig{})
	// 设置 callback 为 qw robot
	logger.SetCallBack(logger.NewQwRobot(qwr))
	// 一般的日志
	logger.Error("roboto test")
	// 调用 notify 后，会触发 callback，也就是 qw robot
	logger.Notify().Error("roboto test")
}

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

func LogExtraFields() {
	// 通过 ctx 传入 trace 信息
	var ctx = trace.NewTrace()
	logger.Ctx(ctx).Info("log ctx")

	// 也可手动设置 traceId， spanId，
	traceId, spanId := trace.NewTraceID(), trace.NewSpanID()
	logger.TraceId(traceId).SpanId(spanId).Trace("log traceId, spanId")

	var err = errors.New("This is err ")
	logger.SetErr(err).Error("I got a err", err)

	// 打印日志时，加入额外的字段
	logger.SetExtra(logger.E().String("name", "xh").String("level", "2")).Info("log extra fields")
}

func main() {
	sample()

	InitLoggerByConfExample()

	LogExtraInfo()

	SetLoggerCallBack()
}
