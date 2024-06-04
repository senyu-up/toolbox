package global

import (
	"context"
	"github.com/senyu-up/toolbox/combz/facade"
	"github.com/senyu-up/toolbox/example/config"
)

var (
	conf       *config.Config                             // 全局配置，不能直接编辑
	ConfigPath string                                     // 配置文件路径
	ErrChan    = make(chan error, 10)                     // 全局错误 channel
	Ctx, Canel = context.WithCancel(context.Background()) // 全局 context
	f          *facade.ToolFacade                         // 门面
)

func GetConfig() *config.Config {
	if conf == nil {
		panic("config not init")
	}
	return conf
}

func SetConfig(c *config.Config) {
	conf = c
}

func GetFacade() *facade.ToolFacade {
	if f == nil {
		panic("facade not init")
	}
	return f
}

func SetFacade(tf *facade.ToolFacade) {
	f = tf
}
