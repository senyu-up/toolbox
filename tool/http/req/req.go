package req

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/imroc/req/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/opentracing/opentracing-go"
	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/trace"
	"golang.org/x/net/http2"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	ctx     context.Context
	c       *req.Client
	r       *req.Request
	timeout time.Duration
	traceOn bool // 设置是否开启 trace
}

// New
// @description 获取client示例
func New(ctx context.Context) *Client {
	c := req.NewClient().SetJsonMarshal(jsoniter.Marshal).SetJsonUnmarshal(jsoniter.Unmarshal)
	r := c.NewRequest()
	r.SetContext(ctx)

	return &Client{
		ctx: ctx,
		r:   r,
		c:   c,
	}
}

func (c *Client) Proxy(proxy string) *Client {
	c.c.SetProxyURL(proxy)

	return c
}

// Timeout
// @description 设置超时时间
func (c *Client) Timeout(d time.Duration) *Client {
	c.c.SetTimeout(d)
	return c
}

// TryTimes
// @description 请求失败后的重试次数
func (c *Client) TryTimes(n int) *Client {
	c.r.SetRetryCount(n)

	return c
}

// Retry Config
type RetryConfig struct {
	// 重试次数, 可选
	Count int
	// 重试间隔, 可选
	Interval time.Duration
	// 重试判断的条件, 比如基于响应内容code判断是否服务繁忙, 可选
	Condition func(resp *req.Response, err error) bool
}

// 自定义重试方案
func (c *Client) Retry(cnf *RetryConfig) *Client {
	if cnf.Count > 0 {
		c.r.SetRetryCount(cnf.Count)
	}
	if cnf.Interval > 0 {
		c.r.SetRetryFixedInterval(cnf.Interval)
	}

	if cnf.Condition != nil {
		c.r.SetRetryCondition(cnf.Condition)
	}

	return c
}

type TLSConfig struct {
	CertFile string
	KeyFile  string
}

// TLS
// @description 设置证书
// 参考示例:
func (c *Client) TLS(cnf *TLSConfig) *Client {
	cert, err := tls.LoadX509KeyPair(cnf.CertFile, cnf.KeyFile)
	if err != nil {
		fmt.Println(err, "密钥对不合法")
		return c
	}
	ssl := &tls.Config{
		Certificates: []tls.Certificate{cert},
		//InsecureSkipVerify: true,
	}
	// 有了证书, 选择使用http2进行通信
	c.c.GetClient().Transport = &http2.Transport{
		TLSClientConfig: ssl,
	}

	return c
}

// Insecure
// @description 涉及到证书校验时, 忽略校验
func (c *Client) Insecure() *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c.c.GetClient().Transport = tr

	return c
}

func (c *Client) SetJsonMarshal(marshal func(v interface{}) ([]byte, error)) *Client {
	c.c.SetJsonMarshal(marshal)

	return c
}

func (c *Client) SetJsonUnmarshal(unmarshal func(data []byte, v interface{}) error) *Client {
	c.c.SetJsonUnmarshal(unmarshal)

	return c
}

// Debug
// @description 开启调试模式, 收到响应后, 会将响应的原始数据全量打印出来
func (c *Client) Debug() *Client {
	c.c.EnableDumpAll()

	return c
}

// Context
// @description 中途修改context
func (c *Client) Context(ctx context.Context) *Client {
	reqId, spanId := trace.ParseCurrentContext(ctx)
	if reqId == "" {
		reqId = trace.NewTraceID()
	}
	if spanId == "" {
		spanId = trace.NewSpanID()
	}
	c.r.SetContext(ctx)
	c.r.SetHeader(enum.RequestId, reqId)
	c.r.SetHeader(enum.SpanId, spanId)

	return c
}

// RequestId
// @description 使用脚本触发, 自定义requestId
func (c *Client) RequestId(id string) *Client {
	c.r.SetHeader(enum.RequestId, id)

	return c
}

// Header
// @description 设置header
func (c *Client) Header(header map[string]string) *Client {
	if header == nil {
		return c
	}
	for k, v := range header {
		c.r.SetHeader(k, v)
	}

	return c
}

// SetTraceOn
// @description 设置是否开启trace
func (c *Client) SetTraceOn(traceOn bool) *Client {
	c.traceOn = traceOn
	return c
}

// SetBasicToken
// @description 设置basic认证的token, 会在token前追加 Basic
func (c *Client) SetBasicToken(token, psd string) *Client {
	c.r.SetBasicAuth(token, psd)

	return c
}

// SetBearerToken
// @description 设置bearer认证的token, 会在token前追加 Bearer
func (c *Client) SetBearerToken(token string) *Client {
	c.r.SetBearerAuthToken(token)

	return c
}

func (c *Client) SetCertFromFile(clientPemFile, keyPemFile string) *Client {
	c.c.SetCertFromFile(clientPemFile, keyPemFile)

	return c
}

func (c *Client) SetCertificate(cert tls.Certificate) *Client {
	c.c.SetCerts(cert)

	return c
}

// SetRootCertsFromFile
// @description 设置根证书
func (c *Client) SetRootCertsFromFile(pemFiles ...string) *Client {
	c.c.SetRootCertsFromFile(pemFiles...)

	return c
}

// SetRootCertFromString
// @description 设置根证书
func (c *Client) SetRootCertFromString(certs ...string) *Client {
	c.c.SetRootCertsFromFile(certs...)

	return c
}

// SetToken
// @description 自定义token
func (c *Client) SetToken(token string) *Client {
	c.r.Headers.Set("Authorization", token)

	return c
}

func (c *Client) ForceHttp1() *Client {
	c.c.EnableForceHTTP1()
	return c
}

func (c *Client) ForceHttp2() *Client {
	c.c.EnableForceHTTP2()

	return c
}

// afterRequest
// @description 回收
func (c *Client) afterRequest(resp *req.Response) {
}

// ContentType
// @description 设置编码方式
func (c *Client) ContentType(t string) *Client {
	c.r.SetContentType(t)

	return c
}

// Cookie
// @description 设置cookie
func (c *Client) Cookie(data []*http.Cookie) *Client {
	c.r.SetCookies(data...)
	return c
}

// FormDataByValue
// @description 设置表单数据
func (c *Client) FormDataByValue(val url.Values) *Client {
	c.r.SetContentType("application/x-www-from-urlencoded")
	c.r.SetFormDataFromValues(val)

	return c
}

// FormData
// @description 以表单形式请求
func (c *Client) FormData(data map[string]string) *Client {
	c.r.SetContentType("application/x-www-from-urlencoded")
	c.r.SetFormData(data)

	return c
}

// File
// @description 表单形式上传文件
func (c *Client) File(name, path string) *Client {
	c.r.SetContentType("multipart/form-data")
	c.r.SetFile(name, path)

	return c
}

// BodyJson
// @description json编码格式, data 支持 string, []byte, io.Reader, map, struct, slice 类型
func (c *Client) BodyJson(data interface{}) *Client {
	c.r.SetContentType("application/json")
	c.r.SetBody(data)

	return c
}

// Body
// @description 设置body内容, data 支持 string, []byte, io.Reader, map, struct, slice 类型
func (c *Client) Body(data interface{}) *Client {
	c.r.SetBody(data)
	return c
}

func (c *Client) Post(url string) (resp *req.Response, err error) {
	span := c.Span(url)
	defer func() {
		if span != nil {
			span.Finish()
		}
	}()
	resp, err = c.r.Post(url)
	if err != nil {
		return nil, err
	}
	c.afterRequest(resp)
	return
}

func (c *Client) Span(uri string) opentracing.Span {
	if c.traceOn {
		parse, err := url.Parse(uri)
		if err == nil {
			traceId, pSpanId := trace.ParseCurrentContext(c.ctx)
			if traceId == "" {
				traceId = trace.NewTraceID()
			}
			spanId := trace.NewSpanID()
			return trace.NewJaegerSpan("remote", traceId, spanId, pSpanId, map[string]interface{}{"host": parse.Host}, nil)
		}
	}
	return nil
}

func (c *Client) Get(url string) (resp *req.Response, err error) {
	span := c.Span(url)
	defer func() {
		if span != nil {
			span.Finish()
		}
	}()
	resp, err = c.r.Get(url)
	if err != nil {
		return nil, err
	}
	c.afterRequest(resp)

	return
}

func (c *Client) Delete(url string) (resp *req.Response, err error) {
	span := c.Span(url)
	defer func() {
		if span != nil {
			span.Finish()
		}
	}()
	resp, err = c.r.Delete(url)
	if err != nil {
		return nil, err
	}
	c.afterRequest(resp)

	return
}

func (c *Client) Put(url string) (resp *req.Response, err error) {
	span := c.Span(url)
	defer func() {
		if span != nil {
			span.Finish()
		}
	}()
	resp, err = c.r.Put(url)
	if err != nil {
		return nil, err
	}
	c.afterRequest(resp)

	return
}

func (c *Client) Patch(url string) (resp *req.Response, err error) {
	span := c.Span(url)
	defer func() {
		if span != nil {
			span.Finish()
		}
	}()
	resp, err = c.r.Patch(url)
	if err != nil {
		return nil, err
	}
	c.afterRequest(resp)

	return
}
