package enum

const (
	SuccessCode = 20000
	FailCode    = 40000

	InternalErrCode     = 50000
	ParamsErrCode       = 50001
	ServerIsBusyErrCode = 50002

	OpenFileErrCode    = 50003
	IsAccessModelsCode = 180000
	FileMustBeJsonCode = 50004
)

const (
	RPCFailedErrMsg    = "failed request rpc service "
	ParamsErrMsg       = "params error or missing "
	ServerIsBusyErrMsg = "server is busy, please try again later"
	UniqueDataMsg      = "数据已存在，请调整后提交"
	OpenFileMsg        = "文件打开失败"
	FileMustBeJsonMsg  = "文件必须为Json格式"
	DbFailedErrorMsg   = "数据库操作失败"
)
