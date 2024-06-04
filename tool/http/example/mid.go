package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/senyu-up/toolbox/tool/http/fiber/middleware"
)

func ExampleMid() {
	var app = fiber.New(fiber.Config{})
	app.Use(middleware.Cors(),
		middleware.Logger(),           // 日志
		middleware.SetRequestId(true), // 启用 traceId
		middleware.Health(),           // 健康检查路由开启
		middleware.Swagger(),          // swagger 文档展示
		middleware.PanicRecover(),     // cover panic
		//middleware.Pprof(),            // http pprof 采集开启
	)
}
