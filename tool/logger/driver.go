package logger

import "time"

// Driver
// @Description:  Driver interface
type Driver interface {
	LogWrite(when time.Time, msg string, level LogLevel, extra []Field) error
	Destroy()
	// CurrentLevel
	// @Description: 获取当前的日志输出级别
	CurrentLevel() LogLevel
	Name() string
}
