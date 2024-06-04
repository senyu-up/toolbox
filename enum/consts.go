package enum

import (
	"time"
)

const (
	EvnStageLocal      = "local"
	EvnStageDevelop    = "develop"
	EvnStageRelease    = "release"
	EvnStagePre        = "pre"
	EvnStageProduction = "production"
)

// context value key
const (
	AuthInfo = "AuthInfoKey"
	// 请求的链路id
	RequestId    = "RequestId"
	ParentSpanId = "ParentSpanId"
	// 请求的每一跳的id
	SpanId = "SpanId"
	//XhSdkVersion sdk版本号
	XhSdkVersion = "XhSdkVersion"
	//XhSource 来源
	XhSource = "XhSource"
	//XhOs 系统
	XhOs = "XhOs"
	//XhAppKey appkey
	XhAppKey      = "XhAppKey"
	CommonRespKey = "CommonResp"
)

// expired time
const (
	RequestLockExpiredTime = time.Minute * 3
)

const (
	DayUnixTime   = 3600 * 24
	WeekUnixTime  = DayUnixTime * 7
	MonthUnixTime = DayUnixTime * 30
	YearUnixTime  = DayUnixTime * 365
)

const (
	StageKey            string = "stage"
	DebugLocalKey       string = "debugLocal"
	DevopsServerHostKey string = "devops_server_host"
	ServiceNameKey      string = "service_name"
)

const (
	ExportLimitNum = 2000 //导出每次写入数量
	ImportMaxNum   = 3000 //最大导入数量
	DefaultPage    = 1    //默认页码 =
)

const (
	ActionLockTime = 20
)
