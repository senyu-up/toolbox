package gin_server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/senyu-up/toolbox/tool/config"
)

type App struct {
	a   *gin.Engine
	ctf config.GinConfig
}

// NewApp 返回一个新的fiberApp。 conf 配置，传入值结构至少有 Config 里的属性信息
func NewApp(cf *config.GinConfig) (*App, error) {
	t := &App{}
	if cf == nil {
		return t, fmt.Errorf("请传入有效的 gin 配置参数")
	}

	if cf.BodyLimit <= 0 {
		cf.BodyLimit = 4 * 1024 * 1024
	}

	t.ctf = *cf
	t.a = gin.New()
	return t, nil
}

func (t *App) Gin() *gin.Engine {
	return t.a
}

func (t *App) Run() error {
	return t.a.Run(t.ctf.Addr)
}
