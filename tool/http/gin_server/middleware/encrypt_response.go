package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/senyu-up/toolbox/tool/encrypt"
	"github.com/senyu-up/toolbox/tool/http/gin_server/controller"
	"net/http"
)

// EncryptResponse 自定义中间件，用于拦截响应并加密
func EncryptResponse(stage string) gin.HandlerFunc {

	key := []byte("54q55oCc6I2J5pyo6Z2S") // 16位密码串

	return func(c *gin.Context) {

		//本机测试环境不加密
		//if  == enum.Local || config.AppConfigVar.App.Stage == enum.Develop {
		//	c.Next()
		//	return
		//}
		// 替换Writer
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		// 调用下一个处理程序，使得响应体被写入到writer.body中
		c.Next()

		// 获取响应体内容
		bodyBytes := writer.body.Bytes()

		// 解析响应体到CommonResp结构体
		var rsp controller.CommonResp
		if err := json.Unmarshal(bodyBytes, &rsp); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// 对响应体中的Data字段进行加密
		if rsp.Data != nil {
			encrypted := encrypt.AesCBCEncrypt(rsp.Data.([]byte), key)
			rsp.Data = encrypted

			// 重新将修改后的rsp序列化为JSON并写回响应体
			modifiedBody, err := json.Marshal(rsp)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			// 将修改后的响应体写回原始ResponseWriter
			_, _ = writer.ResponseWriter.Write(modifiedBody)
		}
	}

}

// 自定义ResponseWriter，用于截获响应体
type responseWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

// Write 重写 Write 方法，以便将响应数据写入自定义缓冲区
func (r *responseWriter) Write(data []byte) (int, error) {
	return r.body.Write(data)
}
