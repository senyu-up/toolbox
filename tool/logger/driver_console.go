package logger

import (
	"github.com/jinzhu/copier"
	"github.com/senyu-up/toolbox/tool/config"
	"os"
	"runtime"
	"sync"
	"time"
)

type brush func(string) string

func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

// 鉴于终端的通常使用习惯，一般白色和黑色字体是不可行的,所以30,37不可用，
var colors = []brush{
	newBrush("1;41"), // Emergency          红色底
	newBrush("1;35"), // Alert              紫色
	newBrush("1;34"), // Critical           蓝色
	newBrush("1;31"), // Error              红色
	newBrush("1;33"), // Warn               黄色
	newBrush("1;36"), // Informational      天蓝色
	newBrush("1;32"), // Debug              绿色
	newBrush("1;32"), // Trace              绿色
}

type Console struct {
	sync.Mutex
	Level    string `json:"level"`
	Colorful bool   `json:"color"`
	LogLevel LogLevel
}

func (c *Console) InitByConf(conf config.ConsoleConfig) (err error) {
	copier.CopyWithOption(c, conf, copier.Option{IgnoreEmpty: true})
	if runtime.GOOS == "windows" {
		c.Colorful = false
	}

	if l, ok := LevelMap[c.Level]; ok {
		c.LogLevel = l
	} else {
		return ErrInvalidLogLevel
	}

	return err
}

func (c *Console) LogWrite(when time.Time, msg string, level LogLevel, extra []Field) error {
	if level > c.LogLevel {
		return nil
	}
	if c.Colorful {
		msg = colors[level](msg)
	}
	c.printlnConsole(when, msg)
	return nil
}

func (c *Console) WithCaller() {
}

func (c *Console) CurrentLevel() LogLevel {
	return c.LogLevel
}

func (c *Console) Destroy() {

}

func (c *Console) Name() string {
	return AdapterConsole
}

func (c *Console) printlnConsole(when time.Time, msg string) {
	c.Lock()
	defer c.Unlock()
	_, _ = os.Stdout.WriteString(msg + "\n")
}
