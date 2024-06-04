package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/env"
	fiber_app "github.com/senyu-up/toolbox/tool/http/fiber"
	"github.com/senyu-up/toolbox/tool/http/fiber/middleware"
	"github.com/senyu-up/toolbox/tool/trace"
)

func InitJaeger() {
	// 初始化配置
	var config = &config.TraceConfig{
		ServerLogOn:    true,
		ClientLogOn:    true,
		ClientLogLevel: "info",

		Jaeger: config.JaegerConf{
			JaegerOn: true,           // 开启状态才会初始化 jaeger
			AppName:  "example_test", // 当前记录到jeager时，上报的 app name。

			// jeager 采集器地址信息
			CollectorEndpoint: "127.0.0.1:6831", // 采集器地址 UDP 采集！
			//CollectorEndpoint: "http://localhost:14268/api/traces",
			//AgentPort:         "14268",
			User:     "",
			Password: "",

			SamplerFreq:              1,
			QueueFlushIntervalSecond: 1,
		},
	}
	trace.Init(config) // 通过配置初始化 trace，jaeger 客户端
}

func Example1() {
	InitJaeger() // 初始化 jaeger

	// 生命 fiber http 服务的配置
	var conf = &config.FiberConfig{
		Name:              env.GetAppInfo().Name, // 应用名
		Addr:              "0.0.0.0:8181",        //
		CaseSensitive:     false,                 // http 路径是否大小写敏感
		Timeout:           60,                    // 请求处理超时
		EnablePrintRoutes: true,                  // 启动后打印所有路由不
	}
	// 初始化 fiber http 服务
	app, err := fiber_app.NewApp(conf)
	if err != nil {
		fmt.Printf("fiber.NewApp error: %v", err)
		return
	}

	var mid = middleware.Prometheus(app.Fiber(), "local_test", "test") // 请求 `/metrics` 路由返回 prometheus 指标

	// fiber 中间件启用
	app.Fiber().Use(
		mid,
		middleware.Cors(),
		middleware.Logger(),
		middleware.PanicRecover(), // 捕获 panic
		//middleware.Pprof(),            // 请求 `/debug/pprof` 路由查看当前服务的 pprof 信息
		middleware.Health(),           // 请求 `/system/health` 路由返回 "ok"
		middleware.SetRequestId(true), // 启用 traceId
		middleware.TrustIp(map[string]struct{}{ // 设置访问 ip 白名单
			"127.0.0.1": {},
			"0.0.0.0":   {},
		}),
	)
	// fiber 路由
	app.Fiber().Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"code": enum.SuccessCode, "msg": "ok"})
	})

	app.Fiber().Get("/trace", func(c *fiber.Ctx) error {
		reqId, spanId := trace.ParseFromContext(c.Context())
		return c.JSON(fiber.Map{"code": enum.SuccessCode, "msg": fmt.Sprintf("trace info : %s, %s", reqId, spanId)})
	})

	app.Fiber().Get("/do_panic", func(c *fiber.Ctx) error {
		panic("do panic")
		return c.JSON(fiber.Map{"code": enum.SuccessCode, "msg": "Boom!"})
	})

	app.Fiber().Get("/swagger/*", middleware.Swagger()) // 请求 `/swagger` 路由查看当前 api 的 swagger 文档

	// 运行 fiber http 服务
	err = app.Run()
	if err != nil {
		fmt.Printf("app.Run error: %v", err)
	}
}

func main() {
	Example1()
}
