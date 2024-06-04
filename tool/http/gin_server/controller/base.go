package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/senyu-up/toolbox/enum"
	"github.com/spf13/cast"
	"net/http"
	"reflect"
)

type BaseController struct {
}

var validate = validator.New()

type CommonResp struct {
	Code      int32       `json:"code"`                 // 状态码, 0表示成功, 其他表示失败
	Msg       string      `json:"message"`              // 信息
	Data      interface{} `json:"data,omitempty"`       // 接口数
	RequestId string      `json:"request_id,omitempty"` //请求id
}

func (b *BaseController) ParamsValidator(ctx *gin.Context, params interface{}) error {
	refParam := reflect.ValueOf(params)
	if refParam.Kind() != reflect.Ptr {
		return fmt.Errorf("params must be a pointer")
	}
	err := b.ParseJson(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseController) ParseJson(c *gin.Context, params interface{}) error {
	err := c.ShouldBind(params)
	if err != nil {
		return err
	}
	return validate.Struct(params)
}

func (b *BaseController) Call_(ctx *gin.Context, params interface{}, handler interface{}) {
	refParam := reflect.ValueOf(params)
	if params != nil {
		if refParam.Kind() != reflect.Ptr {
			ctx.JSON(http.StatusBadRequest, CommonResp{Code: http.StatusBadRequest, Msg: "params must be a pointer"})
			return
		}
		err := b.ParseJson(ctx, params)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, CommonResp{Code: http.StatusBadRequest, Msg: err.Error()})
			return
		}
	}
	refHandler := reflect.ValueOf(handler)
	if refHandler.Kind() != reflect.Func {
		ctx.JSON(http.StatusBadRequest, CommonResp{Code: http.StatusBadRequest, Msg: "handler must be a function"})
		return
	}
	var rets []reflect.Value
	inParam := []reflect.Value{reflect.ValueOf(ctx)}
	if params != nil {
		inParam = append(inParam, refParam)
	}
	rets = refHandler.Call(inParam)
	// 对返回值进行处理
	resp := CommonResp{Code: enum.SuccessCode, Msg: "success", RequestId: b.GetRequestId(ctx)}
	if len(rets) > 1 {
		if !rets[1].IsNil() {
			err := rets[1].Interface().(error)
			resp.Code = enum.FailCode
			resp.Msg = err.Error()
		}
		resp.Data = rets[0].Interface()
	} else {
		if !rets[0].IsNil() {
			err := rets[0].Interface().(error)
			resp.Code = enum.FailCode
			resp.Msg = err.Error()
		}
	}
	ctx.JSON(http.StatusOK, resp)
}

func (b *BaseController) GetRequestId(ctx *gin.Context) string {
	r, _ := ctx.Get(enum.RequestId)
	if r == nil {
		return ""
	}
	return cast.ToString(r)
}
