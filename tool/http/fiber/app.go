package fiber

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/senyu-up/toolbox/tool/config"
	"time"
)

type App struct {
	a   *fiber.App
	ctf config.FiberConfig
}

// NewApp 返回一个新的fiberApp。 conf 配置，传入值结构至少有 Config 里的属性信息
func NewApp(cf *config.FiberConfig) (*App, error) {
	t := &App{}
	if cf == nil {
		return t, fmt.Errorf("请传入有效的 fiber 配置参数")
	}

	if cf.BodyLimit <= 0 {
		cf.BodyLimit = 4 * 1024 * 1024
	}

	t.ctf = *cf
	//这里的配置很多，后续可自定义
	t.a = fiber.New(fiber.Config{
		CaseSensitive:           cf.CaseSensitive,
		AppName:                 cf.Name,
		ReadTimeout:             time.Millisecond * time.Duration(cf.Timeout),
		WriteTimeout:            time.Millisecond * time.Duration(cf.Timeout),
		EnableTrustedProxyCheck: cf.EnableTrustedProxyCheck,
		BodyLimit:               cf.BodyLimit,
		EnablePrintRoutes:       cf.EnablePrintRoutes,
		Immutable:               true,
	})
	return t, nil
}

func (t *App) Fiber() *fiber.App {
	return t.a
}

func (t *App) Run() error {
	return t.a.Listen(t.ctf.Addr)
}
