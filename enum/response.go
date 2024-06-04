package enum

import (
	"errors"
	"github.com/senyu-up/toolbox/tool/su_error"
)

const (
	NotFundCache = "未找到缓存"
)

// rpc server 响应通用code
const (
	RPCSuccessCode   = 20000
	RPCFailedCode    = 10000
	NoLoginCode      = 40001
	NoPermissionCode = 40002
	UserInfoErrCode  = 40003
	//参数检测接口参数有误时使用,其他接口请务使用
	ParamsCheckErrCode = 41000
)

const (
	RPCFailedDesc  = "Failed"
	RPCSuccessDesc = "Success"
)

// gateway
const (
	// 成功
	SuccessCode         = 2000
	FailCode            = 4000
	ParamsMissingErr    = 4002
	InternalErrCode     = 5000
	ParamsErrCode       = 5001
	ServerIsBusyErrCode = 5002
)
const (
	SuccessDesc        = "Success"
	RPCFailedErrMsg    = "请求rpc服务出错"
	ParamsErrMsg       = "参数错误或丢失"
	ParamsErrDesc      = "参数错误或丢失"
	NoPermissionDesc   = "无此权限"
	NoLoginDesc        = "您还未登录,请使用企业微信扫码登录"
	UserInfoErrDesc    = "用户信息错误"
	ServerIsBusyErrMsg = "服务器繁忙,请稍后再试"
	InternalErrDesc    = "服务器内部错误"
	ParamsMissing      = "参数缺失"
)

var GmConfError = &su_error.SUError{Code: 250010, Msg: "配置获取失败"}
var TheTransactionIDBeingExecutedError = errors.New("The transaction ID being executed! ")
