package index

import (
	"github.com/gofiber/fiber/v2"
	"github.com/senyu-up/toolbox/example/internal/service/center"
)

func RegisterRouter(app *fiber.App) {
	app.Static("/WW_verify_5G2EPZrrFcMRhlcf.txt", "./WW_verify_5G2EPZrrFcMRhlcf.txt")
	// event-bus socket链接
	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})

	c := app.Group("/center")
	{
		c.Get("/app_dsns", center.CenterCtl.ListAppDsns)
		c.Get("/app_dsns/:id", center.CenterCtl.GetAppDsns)
	}
}
