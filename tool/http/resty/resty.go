package resty

import (
	"crypto/tls"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// ResponseJson json 响应约定
type ResponseJson struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 请求超时时间
const httpClientTimeOut = 60

// 请求重试次数
const httpClientRetryCount = 1

// HttpGetResJson get request and json response
func HttpGetResJson(url string, queryParams map[string]string, result interface{}) (res *resty.Response, err error) {
	client := resty.New()
	client.SetTimeout(time.Second * httpClientTimeOut)
	client.SetRetryCount(httpClientRetryCount)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	res, err = client.R().
		SetQueryParams(queryParams).
		SetHeader("Accept", "application/json").
		SetResult(result).
		Get(url)
	return res, err
}

func HttpPostFormResJson(url string, formData map[string]string, result interface{}) (res *resty.Response, err error) {
	return HttpSendFormResJson(url, "POST", formData, result)
}
func HttpPutFormResJson(url string, formData map[string]string, result interface{}) (res *resty.Response, err error) {
	return HttpSendFormResJson(url, "PUT", formData, result)
}
func HttpPatchFormResJson(url string, formData map[string]string, result interface{}) (res *resty.Response, err error) {
	return HttpSendFormResJson(url, "PATCH", formData, result)
}
func HttpDeleteFormResJson(url string, formData map[string]string, result interface{}) (res *resty.Response, err error) {
	return HttpSendFormResJson(url, "DELETE", formData, result)
}
func HttpOptionsFormResJson(url string, formData map[string]string, result interface{}) (res *resty.Response, err error) {
	return HttpSendFormResJson(url, "OPTIONS", formData, result)
}
func HttpHeadFormResJson(url string, formData map[string]string, result interface{}) (res *resty.Response, err error) {
	return HttpSendFormResJson(url, "HEAD", formData, result)
}

func HttpPostJsonResJson(url string, formData map[string]string, result interface{}) (res *resty.Response, err error) {
	return HttpSendJsonResJson(url, "POST", formData, result)
}
func HttpPutJsonResJson(url string, formData map[string]string, result interface{}) (res *resty.Response, err error) {
	return HttpSendJsonResJson(url, "PUT", formData, result)
}
func HttpPatchJsonResJson(url string, formData map[string]string, result interface{}) (res *resty.Response, err error) {
	return HttpSendJsonResJson(url, "PATCH", formData, result)
}
func HttpDeleteJsonResJson(url string, formData map[string]string, result interface{}) (res *resty.Response, err error) {
	return HttpSendJsonResJson(url, "DELETE", formData, result)
}
func HttpOptionsJsonResJson(url string, formData map[string]string, result interface{}) (res *resty.Response, err error) {
	return HttpSendJsonResJson(url, "OPTIONS", formData, result)
}
func HttpHeadJsonResJson(url string, formData map[string]string, result interface{}) (res *resty.Response, err error) {
	return HttpSendJsonResJson(url, "HEAD", formData, result)
}

// HttpSendFormResJson send formData and response json
func HttpSendFormResJson(url, method string, formData map[string]string, result interface{}) (res *resty.Response, err error) {
	client := resty.New()
	client.SetTimeout(time.Second * httpClientTimeOut)
	client.SetRetryCount(httpClientRetryCount)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	req := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		SetFormData(formData).
		SetResult(result)
	switch strings.ToLower(method) {
	case "post":
		res, err = req.Post(url)
	case "put":
		res, err = req.Put(url)
	case "patch":
		res, err = req.Patch(url)
	case "delete":
		res, err = req.Delete(url)
	case "options":
		res, err = req.Options(url)
	default:
		res, err = req.Head(url)
	}
	return res, err
}

// HttpSendJsonResJson send json and response json
func HttpSendJsonResJson(url, method string, body interface{}, result interface{}) (res *resty.Response, err error) {
	client := resty.New()

	client.SetTimeout(time.Second * httpClientTimeOut)
	client.SetRetryCount(httpClientRetryCount)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	req := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(body).
		SetResult(result)

	switch strings.ToLower(method) {
	case "post":
		res, err = req.Post(url)
	case "put":
		res, err = req.Put(url)
	case "patch":
		res, err = req.Patch(url)
	case "delete":
		res, err = req.Delete(url)
	case "options":
		res, err = req.Options(url)
	default:
		res, err = req.Head(url)
	}
	return res, err
}
