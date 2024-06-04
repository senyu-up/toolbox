package middleware

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/su_slice"
	"github.com/senyu-up/toolbox/tool/trace"
	"strings"
)

const (
	TRACE_RECORD_BODY_LIMIT = 10240 // trace记录body的最大长度
)

var (
	ignoreTracePath = []string{"/system/health", "/ping", "/swagger", "/metrics"}
)

// Cors 跨域中间件
func Cors() func(c *fiber.Ctx) error {
	return cors.New(cors.Config{
		//Next func(c *fiber.Ctx) bool ，返回值若为true，跳过此中间件
		Next: nil,
		//来源设定
		AllowOrigins: "*",
		//允许的方法
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
			fiber.MethodOptions,
		}, ","),
	})
}

// Logger 请求打印
func Logger() func(c *fiber.Ctx) error {
	return logger.New()
}

func SetRequestId(traceOn bool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// 判断header中是否包含链路id, 没有则自动生成
		var reqId string
		var pSpanId = c.Get(enum.SpanId)
		if reqId = c.Get(enum.RequestId); reqId == "" {
			reqId = trace.NewTraceID()
		}
		spanId := trace.NewSpanID()
		var url = c.Request().URI()
		var httpPath = string(url.Path())

		// 判断是否需要记录trace, 并且是业务路由
		if traceOn && !su_slice.InArray(httpPath, ignoreTracePath) {
			opName := "http " + httpPath
			// full url
			tags := map[string]interface{}{"uri": string(url.FullURI()), "query": string(url.QueryString())}
			// body, 仅当body不是stream时、不是文件上传才打印, 最多记录1k字符
			if c.Get("Content-Type") != "application/octet-stream" && !c.Request().IsBodyStream() {
				var fullBody = c.Request().Body()
				if len(fullBody) > TRACE_RECORD_BODY_LIMIT {
					tags["body"] = string(fullBody[:TRACE_RECORD_BODY_LIMIT])
				} else if fullBody != nil && len(fullBody) > 0 {
					tags["body"] = string(fullBody)
				}
			}

			var span = trace.NewJaegerSpan(opName, reqId, spanId, pSpanId, tags, nil)
			defer span.Finish()
		}
		c.Context().SetUserValue(enum.RequestId, reqId)
		c.Context().SetUserValue(enum.SpanId, spanId)

		return c.Next()
	}
}

// PanicRecover
func PanicRecover() func(c *fiber.Ctx) error {
	return recover.New(recover.Config{
		EnableStackTrace: true,
	})
}

// Monitor 监控中间件 监视cpu 内存 响应时间
func Monitor() func(c *fiber.Ctx) error {
	return monitor.New()
}

// Pprof pprof页面
// 不建议在http服务是那个做 pprof 暴露，这回让业务方有机会访问pprof，建议使用 healthcheck 服务来做pprof暴露
//Deprecated
//func Pprof() func(c *fiber.Ctx) error {
//	return pprof.New()
//}

// Limiter 限频器
func Limiter() func(c *fiber.Ctx) error {
	return limiter.New()
}

var defaultTrustIp = map[string]struct{}{
	"127.0.0.1": {},
	"0.0.0.0":   {},
	//@todo 加一个公司内网地址
}

// TrustIp 仅开放给某些ip
func TrustIp(trustMap ...map[string]struct{}) func(c *fiber.Ctx) error {
	var tm map[string]struct{}
	tm = defaultTrustIp
	if len(trustMap) > 0 {
		tm = trustMap[0]
	}
	return func(c *fiber.Ctx) error {
		if _, ok := tm[c.IP()]; ok {
			c.Next()
		} else {
			c.WriteString("abort")
		}
		return nil
	}
}

// Swagger 文档
func Swagger() func(c *fiber.Ctx) error {
	return swagger.New(swagger.Config{
		DeepLinking:       true,
		DocExpansion:      "none",
		OAuth:             nil,
		OAuth2RedirectUrl: "",
		URL:               "",
	})
}

// Prometheus 数据监控,这个和标准middleware有些不同
//
//	@Description:
//	@param app      body any true "-"
//	@param appName	body any true "上报指标时，的应用名，建议用 hostname"
//	@param url      body any true "可选参数，指定获取 metrics 的路由，默认为 /metrics"
//	@return func(c *fiber.Ctx) error
func Prometheus(app *fiber.App, appName string, url ...string) func(c *fiber.Ctx) error {
	s := fiberprometheus.New(appName)
	route := "/metrics"
	if len(url) > 0 {
		route = url[0]
	}
	s.RegisterAt(app, route)
	return s.Middleware
}

// Health 给k8s健康检查探针使用
func Health() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/system/health" {
			c.WriteString("ok")
			return nil
		}
		return c.Next()
	}
}
