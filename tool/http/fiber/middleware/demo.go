package middleware

import "github.com/gofiber/fiber/v2"

// 带config的自定义中间件示例
type DemoMiddlewareConfig struct {
	Next func(c *fiber.Ctx) bool
	Msg  string
}

var DemoMiddlewareConfigDefault = DemoMiddlewareConfig{
	Next: func(c *fiber.Ctx) bool {
		if !c.Secure() {
			return false
		}
		return true
	},
	Msg: "you shall not pass!",
}

// DemoMiddleware 带config的中间件示例
func DemoMiddleware(config ...DemoMiddlewareConfig) fiber.Handler {
	ctg := DemoMiddlewareConfigDefault
	if len(config) > 0 {
		ctg = config[0]
	}

	return func(ctx *fiber.Ctx) error {
		if ctg.Next(ctx) {
			return ctx.Next()
		}
		ctx.WriteString(ctg.Msg)
		return nil
	}

}

// DemoMiddlewareWithoutConfig 不带config的中间件示例
func DemoMiddlewareWithoutConfig() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		//get a == 1 时跳过此中间件
		if ctx.Query("a") == "1" {
			ctx.Next()
		}
		ctx.WriteString("about by DemoMiddlewareWithoutConfig")
		return nil
	}
}
