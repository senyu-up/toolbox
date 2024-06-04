package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/encrypt"
	"github.com/senyu-up/toolbox/tool/http/gin_server/controller"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/su_slice"
	"github.com/senyu-up/toolbox/tool/trace"
	"net/http"
	"runtime/debug"
	"strings"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许跨域的域名，可以使用通配符 "*" 允许所有域
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// 设置允许的HTTP方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}, ","))

		// 设置允许的自定义请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 允许发送Cookie
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理OPTIONS预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		// 继续处理请求
		c.Next()
	}
}

// PanicRecoverMiddleware 是一个用于捕获panic并进行恢复的中间件
func PanicRecoverMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// 发生panic时的处理
			logger.Error("Recovered from panic:", err)
			logger.Error(string(debug.Stack()))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()

	// 继续处理请求
	c.Next()
}

func SetRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 判断header中是否包含链路id, 没有则自动生成
		var reqId, pSpanId string
		if spanIdFromC, exists := c.Get(enum.SpanId); exists {
			pSpanId = spanIdFromC.(string)
		}
		if reqIdFromC, exists := c.Get(enum.RequestId); exists {
			reqId = reqIdFromC.(string)
		}
		if reqId == "" {
			reqId = trace.NewTraceID()
		}

		spanId := trace.NewSpanID()
		var url = c.Request.URL
		var httpPath = url.RequestURI()
		opName := "http " + httpPath
		var span = trace.NewJaegerSpan(opName, reqId, spanId, pSpanId, nil, nil)
		defer span.Finish()

		c.Set(enum.SpanId, spanId)
		c.Set(enum.RequestId, reqId)
		c.Next()
	}
}

// AuthMiddleware 鉴权中间件
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 如果是不需要鉴权的接口，跳过 token 验证
		if skipAuthToken(c.Request.URL.Path) {
			c.Next()
			return
		}
		// 获取请求头中的 token
		token := c.GetHeader("token")
		// 如果请求头中没有 token，返回未授权的错误
		if token == "" {
			c.JSON(http.StatusUnauthorized, controller.CommonResp{
				Code: http.StatusUnauthorized,
				Msg:  "token缺失，请先登录获取token",
			})
			c.Abort()
			return
		}
		// 在这里可以对 token 进行解析和验证的逻辑
		claims, err := encrypt.ParseToken(token, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, controller.CommonResp{
				Code: http.StatusUnauthorized,
				Msg:  "token解析失败,err:" + err.Error(),
			})
			c.Abort()
			return
		}
		// 将解析出的用户信息存储到请求的上下文中
		c.Set("user_info", claims.Data)
		// 鉴权通过，继续处理请求
		c.Next()
	}
}

var SkipPath []string

// skipAuthToken 跳过token验证
func skipAuthToken(path string) bool {
	return su_slice.InArray(path, SkipPath)
}
