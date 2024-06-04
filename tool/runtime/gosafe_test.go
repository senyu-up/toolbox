package runtime

import (
	"go.uber.org/zap/zapcore"
	"testing"
	"time"
	"toolbox/tool/config"
	"toolbox/tool/logger"
	"toolbox/tool/trace"
)

// 测试 go safe 和 go 两个函数
func TestGOSafe(t *testing.T) {
	task := func() {
		logger.Info("task start")
		panic("test")
	}
	var ctx = trace.NewTrace()
	GOSafe(ctx, "rest", task)
	time.Sleep(1 * time.Second)
	logger.Warn("go-safe down")

	logger.InitDefaultLoggerByConf(&config.LogConfig{
		AppName:    "xh_test",
		DefaultLog: "zap",
		Zap: &config.ZapConfig{
			LogLevel: int(zapcore.DebugLevel),
			Level:    "Info",
		},
	})
	Go(ctx, "rest", task) // 这个调用不会吃掉 panic
	time.Sleep(1 * time.Second)
	logger.Warn("go down")
}
