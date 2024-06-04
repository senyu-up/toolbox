package controller

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc/metadata"

	toolEnum "github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/encrypt"
	"github.com/senyu-up/toolbox/tool/logger"
	error2 "github.com/senyu-up/toolbox/tool/su_error"
	"github.com/senyu-up/toolbox/tool/trace"
	"github.com/senyu-up/toolbox/tool/validator"
)

// BaseController 合理封装一下常用方法
type BaseController struct {
	httpStatus int32
	inst       redis.UniversalClient
}
type CommonResp struct {
}
type AuthInfo struct {
	AppKey string
	Name   string
}

type JsonResponse struct {
	// 状态码, 若指定会直接使用, 反之会从error中获取 或者 使用默认值
	Code int32 `json:"code"`
	// 描述, 若指定会直接使用, 反之会从error中获取 或者 使用默认值
	Msg string `json:"message"`
	// 接口数据
	Data interface{} `json:"data"`
	//请求id
	RequestId string `json:"request_id"`
	// Http Status Code
	HttpStatus int `json:"-"`
}

type PageData struct {
	// 总记录数
	Total int `json:"total,omitempty"`
	// 基于标记为获取分页时使用, eg select * from log where id>{next} limit 100
	Next string `json:"next,omitempty"`
	// 当前返回的记录数 场景: 流式拉取列表使用, Items的记录在处理过程中有可能会被过滤掉, 导致前端通过len(items)判断是否还有更多不准确
	ItemCount int `json:"item_count,omitempty"`
	// 记录列表
	Items interface{} `json:"items"`
}

type JsonResponseWithPage struct {
	// 状态码, 若指定会直接使用, 反之会从error中获取 或者 使用默认值
	Code int32 `json:"code"`
	// 描述, 若指定会直接使用, 反之会从error中获取 或者 使用默认值
	Msg string `json:"message"`
	//具体数据
	Data PageData `json:"data"`
	//request id
	RequestId  string `json:"request_id"`
	HttpStatus int    `json:"-"`
}

// ResponseJson 返回json信息
func (b *BaseController) ResponseJson(c *fiber.Ctx, res JsonResponse) error {
	res.RequestId = b.GetRequestId(c)
	return c.JSON(res)
}

// SuccessOrServiceErrJson 返回成功或者失败json信息
func (b *BaseController) SuccessOrServiceErrJson(c *fiber.Ctx, res interface{}, err error) error {
	if err != nil {
		return b.ResponseJson(c, JsonResponse{
			Code: toolEnum.InternalErrCode,
			Msg:  toolEnum.InternalErrDesc,
			Data: nil,
		})
	}
	if res == nil {
		return b.ResponseJson(c, JsonResponse{
			Code: toolEnum.SuccessCode,
			Msg:  toolEnum.SuccessDesc,
			Data: nil,
		})
	}
	val := reflect.ValueOf(res)

	if val.IsValid() {
		elm := val.Elem()
		cod := int32(elm.FieldByName("Code").Int())
		msg := elm.FieldByName("Message").String()
		data := elm.FieldByName("Data").Interface()

		if cod != toolEnum.SuccessCode {
			return b.ResponseJson(c, JsonResponse{
				Code: cod,
				Msg:  msg,
				Data: data,
			})
		}

		return b.ResponseJson(c, JsonResponse{
			Code: toolEnum.SuccessCode,
			Msg:  toolEnum.SuccessDesc,
			Data: data,
		})
	} else {
		return b.ResponseJson(c, JsonResponse{
			Code: toolEnum.SuccessCode,
			Msg:  toolEnum.SuccessDesc,
			Data: nil,
		})
	}
}

func (b *BaseController) GetRequestId(ctx *fiber.Ctx) string {
	return cast.ToString(ctx.Context().Value(toolEnum.RequestId))
}

// ResponseSuccessDataJson 返回json信息
func (b *BaseController) ResponseSuccessDataJson(c *fiber.Ctx, data interface{}, others ...interface{}) error {
	if len(others) > 0 {
		for _, other := range others {
			if commonResp, ok := other.(*unionResp); ok {
				if commonResp != nil && commonResp.Code != toolEnum.RPCSuccessCode {
					res := JsonResponse{
						Code:      commonResp.Code,
						Msg:       commonResp.Msg,
						Data:      nil,
						RequestId: b.GetRequestId(c),
					}
					return c.JSON(res)
				}
				continue
			}
			if commonResp, ok := other.(*unionResp); ok {
				if commonResp != nil && commonResp.Code != toolEnum.RPCSuccessCode {
					res := JsonResponse{
						Code:      commonResp.Code,
						Msg:       commonResp.Msg,
						Data:      nil,
						RequestId: b.GetRequestId(c),
					}
					return c.JSON(res)
				}
				continue
			}
			if commonResp, ok := other.(*unionResp); ok {
				if commonResp != nil && commonResp.Code != toolEnum.RPCSuccessCode {
					res := JsonResponse{
						Code:      commonResp.Code,
						Msg:       commonResp.Msg,
						Data:      nil,
						RequestId: b.GetRequestId(c),
					}
					return c.JSON(res)
				}
				continue
			}
			if err, ok := other.(error); ok {
				if err != nil {
					res := JsonResponse{
						Code:      toolEnum.FailCode,
						Msg:       err.Error(),
						Data:      nil,
						RequestId: b.GetRequestId(c),
					}
					return c.JSON(res)
				}
				continue
			}
		}

	}
	if commonResp, ok := data.(*unionResp); ok {
		if commonResp != nil {
			res := JsonResponse{
				Code:      commonResp.Code,
				Msg:       commonResp.Msg,
				Data:      nil,
				RequestId: b.GetRequestId(c),
			}
			return c.JSON(res)
		}
	}
	if commonResp, ok := data.(*unionResp); ok {
		if commonResp != nil {
			res := JsonResponse{
				Code:      commonResp.Code,
				Msg:       commonResp.Msg,
				Data:      nil,
				RequestId: b.GetRequestId(c),
			}
			return c.JSON(res)
		}
	}
	res := JsonResponse{
		Code:      toolEnum.SuccessCode,
		Msg:       toolEnum.SuccessDesc,
		Data:      data,
		RequestId: b.GetRequestId(c),
	}
	return c.JSON(res)
}

func (b *BaseController) Err(msg string) error {
	logger.Error(msg)
	return errors.New(msg)
}

// CommonReturn 通常型请求模型 paramIn 为rpc入参且必须是指针,out为将要返回给前端的数据，对data进行赋值即可
func (b *BaseController) CommonReturn(c *fiber.Ctx, paramIn interface{}, grpcDo func(out *JsonResponse) error) (err error) {
	jsonResp := JsonResponse{
		Code: toolEnum.ParamsErrCode,
		Msg:  toolEnum.ParamsErrDesc,
		Data: nil,
	}
	defer func() {
		err = b.ResponseJson(c, jsonResp)
	}()
	if err = b.ParseJson(c, paramIn); err != nil {
		jsonResp.Msg = err.Error()
		return
	}
	jsonResp.Msg = ""
	err = grpcDo(&jsonResp)
	if err != nil {
		jsonResp.Code = toolEnum.FailCode
		if jsonResp.Msg == "" {
			jsonResp.Msg = toolEnum.RPCFailedErrMsg
		}
		return
	}

	jsonResp.Code = toolEnum.SuccessCode
	jsonResp.Msg = toolEnum.SuccessDesc
	jsonResp.RequestId = b.GetRequestId(c)
	return
}

// ParseJson param必须是指针。 从ctx里读取信息并赋值到param,且进行注解校验
func (b *BaseController) ParseJson(c *fiber.Ctx, param interface{}) error {
	err := c.BodyParser(param)
	if err != nil {
		return error2.NewSUError(toolEnum.ParamsErrCode, toolEnum.ParamsErrDesc+"err:"+err.Error())
	}
	return validator.StructValidator(param)
}

// Userinfo  获取登陆用户信息
func (b *BaseController) Userinfo(c *fiber.Ctx) *AuthInfo {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	info := c.Context().UserValue(toolEnum.AuthInfo)
	if info == nil {
		return nil
	}
	authInfo := info.(*AuthInfo)
	return authInfo
}

// RPCCtx 将上游数据传递到下游RPC服务
func (b *BaseController) RPCCtx(c *fiber.Ctx) (ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			ctx = c.Context()
			return
		}
	}()

	headers := c.GetReqHeaders()
	param := map[string]string{
		toolEnum.XhSdkVersion: headers["Xh-Sdk-Version"][0],
		toolEnum.XhSource:     headers["Xh-Source"][0],
		toolEnum.XhOs:         headers["Xh-Os"][0],
		toolEnum.XhAppKey:     headers["Xh-App-Key"][0],
		toolEnum.RequestId:    "",
		toolEnum.SpanId:       "",
		toolEnum.AuthInfo:     "",
	}

	c.Context().VisitUserValues(func(bytes []byte, i interface{}) {
		k := string(bytes)
		var v string
		if k == toolEnum.AuthInfo {
			if authInfo, ok := i.(*AuthInfo); ok {
				v, _ = jsoniter.MarshalToString(authInfo)
			}
			v = encrypt.Base64EncodeString(v)
		} else {
			v = cast.ToString(i)
		}
		param[k] = v
	})

	if reqId, _ := param[toolEnum.RequestId]; reqId == "" {
		param[toolEnum.RequestId] = trace.NewTraceID()
	}
	if spanId, _ := param[toolEnum.SpanId]; spanId == "" {
		param[toolEnum.SpanId] = trace.NewSpanID()
	}

	md := metadata.New(param)

	return metadata.NewOutgoingContext(context.Background(), md)
}

// 获取结构体中的某个字段
func GetStructField(in interface{}, columnName ...string) (res reflect.Value, err error) {
	if in == nil {
		err = errors.New("nil ptr")
		return
	}
	t := reflect.TypeOf(in)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		err = errors.New("not a Struct")
		return
	}
	nameMap := map[string]struct{}{}
	for _, v := range columnName {
		nameMap[strings.ToUpper(v)] = struct{}{}
	}
	fieldNum := t.NumField()
	val := reflect.ValueOf(in)
	for i := 0; i < fieldNum; i++ {
		n := strings.ToUpper(t.Field(i).Name)
		if _, ok := nameMap[n]; ok {
			v := val.Elem().Field(i)
			res = reflect.Indirect(v)
			break
		}
	}
	return
}

// RspCheck 对rpc的返回结果进行检查，若有错误则直接向前端返回错误信息并返回false,否则返回true
// res 结构体应为{Code:200 int,Message:string}
func (b *BaseController) RspCheck(c *fiber.Ctx, res interface{}, e error) bool {
	if e != nil {
		msg := toolEnum.RPCFailedErrMsg
		if e.Error() != "" {
			msg = e.Error()
		}
		_ = b.ResponseJson(c, JsonResponse{
			Code: toolEnum.FailCode,
			Msg:  msg,
			Data: nil,
		})
		return false
	}
	codeV, err := GetStructField(res, "Code")
	if err != nil {
		return false
	}
	msgV, err := GetStructField(res, "Message", "msg")
	if err != nil {
		return false
	}
	code := codeV.Int()
	errMsg := msgV.String()
	if code != toolEnum.RPCSuccessCode {
		if errMsg == "" {
			errMsg = toolEnum.RPCFailedDesc
		}
		_ = b.ResponseJson(c, JsonResponse{
			Code: toolEnum.FailCode,
			Msg:  errMsg,
			Data: nil,
		})
		return false
	}
	return true
}

func (b *BaseController) IP(ctx *fiber.Ctx) string {
	var ip string
	ip = ctx.Get("X-Real-Ip")
	ip = ctx.Get("True-Client-Ip")
	if ip == "" {
		ip = ctx.Get("X-Real-Ip")
	}
	if ip == "" {
		ip = ctx.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = ctx.IP()
	}
	if ip != "" {
		ip = strings.Split(ip, ",")[0]
		ip = strings.Split(ip, ":")[0]
	}
	return ip
}

// RequestLock , RequestUnlock 两个方法是对请求添加唯一锁
// 理论上,只有涉及到跟游戏进行交互的请求才需要加请求限定,确保对相同'资源'的操作唯一性
// 调用 RequestLock,必须手动调用 RequestUnlock 进行解锁
// RequestLock 的有效时长为3秒
func (b *BaseController) RequestLock(ctx *fiber.Ctx) bool {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("RequestLock:", r)
		}
	}()
	res := ctx.Context().UserValue(toolEnum.AuthInfo)
	if res == nil {
		logger.Error("RequestLock,auth info is nil.")
		return true
	}
	info := res.(*AuthInfo)
	key := ""
	locked, err := b.inst.SetNX(key, info.Name, toolEnum.RequestLockExpiredTime).Result() // TODO
	if err != nil {
		logger.Error("RequestLock:", err)
		return true
	}
	if !locked {
		_ = b.ResponseJson(ctx, JsonResponse{
			Code: toolEnum.ServerIsBusyErrCode,
			Msg:  toolEnum.ServerIsBusyErrMsg,
		})
	}
	return locked
}

func (b *BaseController) RequestUnlock(ctx *fiber.Ctx) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("RequestUnlock:", r)
		}
	}()
	key := ""
	b.inst.Del(key)
}

func msgAndCodeParser(code int32, msg string, err error) (int32, string) {
	if code > 0 && msg != "" {
		return code, msg
	}
	var errCode int32
	var errMsg string
	if err != nil {
		switch err.(type) {
		case *error2.SUError:
			tmp := strings.SplitN(err.Error(), "$", 2)
			if len(tmp) == 2 {
				c, _ := strconv.Atoi(tmp[0])
				errCode = int32(c)
				errMsg = tmp[1]
			}
		default:
			errCode = toolEnum.FailCode
			errMsg = err.Error()
		}
	}

	if code == 0 {
		if errCode > 0 {
			code = errCode
		} else {
			code = toolEnum.SuccessCode
		}
	}

	if msg == "" {
		if errMsg != "" {
			msg = errMsg
		} else {
			msg = toolEnum.SuccessDesc
		}
	}

	return code, msg
}

func (b *BaseController) Response(ctx *fiber.Ctx, respData JsonResponse, err error) error {
	if respData.HttpStatus > 0 {
		ctx.Status(respData.HttpStatus)
	}

	respData.Code, respData.Msg = msgAndCodeParser(respData.Code, respData.Msg, err)
	respData.RequestId = b.GetRequestId(ctx)

	return ctx.JSON(respData)
}

func (b *BaseController) ResponseWithPage(ctx *fiber.Ctx, respData JsonResponseWithPage, err error) error {
	if respData.HttpStatus > 0 {
		ctx.Status(respData.HttpStatus)
	}

	respData.Code, respData.Msg = msgAndCodeParser(respData.Code, respData.Msg, err)
	respData.RequestId = ctx.UserContext().Value(toolEnum.RequestId).(string)

	return ctx.JSON(respData)
}

//var jsonBeginWith = []byte("{")
//var emptyRequestIdMark = []byte(`"request_id":""`)

type unionResp struct {
	Code       int32       `json:"code"`
	Msg        string      `json:"msg"`
	RequestId  string      `json:"request_id"`
	CommonResp *CommonResp `json:"common_resp"`
	Items      interface{} `json:"items"`
	Data       interface{} `json:"data"`
}

// RPCCall
// @description rpc调用过程封装, 仅适合单个rpc调用, 注意, 当返回数据第一层级包含了 Items 时, 会自动将其调整到 Data.Items 下
// 示例:
// rpcParam := &dataoperation2.GetBaseInfoReq{}
// return  c.RPCCall(ctx, rpcParam, grpc_clients.DataOperationRoleClient().GetBaseInfo)
func (b *BaseController) RPCCall(ctx *fiber.Ctx, param interface{}, handler interface{}, overwriteParam ...bool) error {
	l := logger.Ctx(nil)

	refParam := reflect.ValueOf(param)
	if refParam.Kind() != reflect.Ptr {
		l.Error("param is not a pointer")
		return error2.NewSUError(toolEnum.FailCode, "WRONG param")
	}
	var overwriteFlag = true
	if len(overwriteParam) > 0 {
		overwriteFlag = overwriteParam[0]
	}
	if overwriteFlag && len(ctx.Body()) > 0 {
		err := b.ParseJson(ctx, param)
		if err != nil {
			return b.Response(ctx, JsonResponse{Code: toolEnum.ParamsCheckErrCode}, err)
		}
	}

	refHandler := reflect.ValueOf(handler)
	if refHandler.Kind() != reflect.Func {
		return error2.NewSUError(toolEnum.FailCode, "WRONG handler")
	}
	rpcCxt := b.RPCCtx(ctx)
	var rets []reflect.Value
	inParam := []reflect.Value{reflect.ValueOf(rpcCxt), refParam}
	rets = refHandler.Call(inParam)

	if !rets[1].IsNil() {
		err := rets[1].Interface().(error)
		return b.Response(ctx, JsonResponse{}, error2.NewSUError(toolEnum.InternalErrCode, err.Error()))
	}

	// 对rpc响应内容进行处理
	byteData, _ := jsoniter.Marshal(rets[0].Interface())
	reqId, _ := trace.ParseCurrentContext(rpcCxt)
	resp := responseDecorator(byteData, reqId)
	ctx.Response().Header.Add("Content-Type", fiber.MIMEApplicationJSON)
	_, err := ctx.WriteString(resp)

	return err
}

// RPCCall
// @description rpc调用过程封装, 仅适合单个rpc调用, 注意, 当返回数据第一层级包含了 Items 时, 会自动将其调整到 Data.Items 下
// 示例:
// rpcParam := &dataoperation2.GetBaseInfoReq{}
// return  c.RPCCall(ctx, rpcParam, grpc_clients.DataOperationRoleClient().GetBaseInfo)
func (b *BaseController) RPCCallHander(ctx *fiber.Ctx, param interface{}, handler interface{}, overwriteParam ...bool) error {
	l := logger.Ctx(nil)
	refParam := reflect.ValueOf(param)
	if refParam.Kind() != reflect.Ptr {
		l.Error("param is not a pointer")
		return error2.NewSUError(toolEnum.FailCode, "WRONG param")
	}
	var overwriteFlag = true
	var _buffer bytes.Buffer
	if len(overwriteParam) > 0 {
		overwriteFlag = overwriteParam[0]
	}
	if overwriteFlag && len(ctx.Body()) > 0 {
		err := b.ParseJson(ctx, param)
		if err != nil {
			return b.Response(ctx, JsonResponse{Code: toolEnum.ParamsCheckErrCode}, err)
		}
	}
	refHandler := reflect.ValueOf(handler)
	if refHandler.Kind() != reflect.Func {
		return error2.NewSUError(toolEnum.FailCode, "WRONG handler")
	}
	rpcCxt := b.RPCCtx(ctx)
	var rets []reflect.Value
	inParam := []reflect.Value{reflect.ValueOf(rpcCxt), refParam}
	rets = refHandler.Call(inParam)

	if !rets[1].IsNil() {
		err := rets[1].Interface().(error)
		return b.Response(ctx, JsonResponse{}, error2.NewSUError(toolEnum.InternalErrCode, err.Error()))
	}
	marshaler := jsonpb.Marshaler{
		// 是否将枚举值设定为整数，而不是字符串类型.
		EnumsAsInts: true,
		// 是否将字段值为空的渲染到JSON结构中
		EmitDefaults: true,
		//是否使用原生的proto协议中的字段
		OrigName: true,
	}
	// 对rpc响应内容进行处理
	valuep, ok := rets[0].Interface().(proto.Message)
	if !ok {
		return nil
	}
	marshaler.Marshal(&_buffer, valuep)
	byteData := _buffer.Bytes()
	reqId, _ := trace.ParseCurrentContext(rpcCxt)
	resp := responseDecorator(byteData, reqId)
	ctx.Response().Header.Add("Content-Type", fiber.MIMEApplicationJSON)
	_, err := ctx.WriteString(resp)
	return err
}

func responseDecorator(byteData []byte, reqId string) string {
	respData := strings.Builder{}
	dataVal := strings.Builder{}

	respData.WriteString(`{"request_id":"` + reqId + `"`)
	gjson.ParseBytes(byteData).ForEach(func(k, v gjson.Result) bool {
		curKey := k.String()
		switch curKey {
		case "msg", "message":
			respData.WriteString(fmt.Sprintf(`,"message":%v`, strconv.Quote(v.String())))
		case "code":
			respData.WriteString(`,"code":` + v.String())
		case "error":
			respData.WriteString(`,"error":"` + v.String() + `"`)
		case "common_resp":
			curMsg := v.Get("Msg").String()
			if curMsg == "" {
				curMsg = v.Get("message").String()
			}
			if curMsg == "" {
				curMsg = v.Get("msg").String()
			}
			curCode := v.Get("Code").String()
			if curCode == "" {
				curCode = v.Get("code").String()
			}
			respData.WriteString(fmt.Sprintf(`,"message":%v`, strconv.Quote(curMsg)))
			respData.WriteString(`,"code":` + curCode)
		case "request_id":
			// 不处理
		case "data":
			appendToDataVal(&dataVal, "", v.String(), v.Type)
		default:
			appendToDataVal(&dataVal, k.String(), v.String(), v.Type)
		}

		return true
	})
	respData.WriteString(`,"data":`)
	dataValStr := dataVal.String()
	if dataValStr == "" || dataValStr == "null" {
		respData.WriteString("null")
	} else {
		if strings.HasPrefix(dataValStr, "{") || strings.HasPrefix(dataValStr, "[") {
			respData.WriteString(dataValStr)
		} else {
			respData.WriteString("{" + dataValStr + "}")
		}
	}
	respData.WriteString("}")

	return respData.String()
}

func appendToDataVal(data *strings.Builder, key, val string, valType gjson.Type) {
	if data.Len() > 0 {
		data.WriteByte(',')
	}

	if key != "" {
		data.WriteString(`"` + key + `":`)
	}
	if valType == gjson.String {
		// 适配[]byte类型
		jsonByte, err := base64.StdEncoding.DecodeString(val)
		if err == nil && len(jsonByte) > 0 && (jsonByte[0] == '{' || jsonByte[0] == '[') {
			// 认为是json串
			data.WriteString(string(jsonByte))
		} else {
			data.WriteString(`"` + val + `"`)
		}
	} else if valType == gjson.Null {
		data.WriteString(`null`)
	} else {
		data.WriteString(val)
	}
}
