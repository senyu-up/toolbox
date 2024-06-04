package cronv2

import (
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/logger"
)

type cronLogger struct {
	log logger.Log
}

func NewCronLogger(l logger.Log) cronLogger {
	if l == nil {
		l, _ = logger.InitLoggerByConf(&config.LogConfig{
			CallDepth: 4,
		})
	}
	return cronLogger{log: l}
}

func (c cronLogger) Info(msg string, keysAndValues ...interface{}) {
	c.log.Info(msg, keysAndValues...)
}

func (c cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	c.log.SetErr(err).Error(msg, keysAndValues...)
}
