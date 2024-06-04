package config

type GrpcServerConfig struct {
	Host string `yaml:"host"` // Host Ip地址
	Port uint32 `yaml:"port"` // Port

	SlowThreshold int64 `yaml:"slowThreshold"` // 请求处理超过多少时间则判定为慢请求，慢请求会产生 warn 日志，单位：time.Millisecond
	TimeOut       int64 `yaml:"timeOut"`       // 每次grpc请求最长处理时间， 单位：time.Millisecond
	RequestLogOn  bool  `yaml:"requestLogOn"`  // 是否记录每次 grpc 请求， 记录内容包含： 参数，与返回值
	TraceOn       bool  `yaml:"traceOn"`       // 是否开启链路追踪
}

type GrpcClientConfig struct {
	DebugLocal       string `yaml:"debugLocal"`       // debug local
	DevopsServerHost string `yaml:"devopsServerHost"` // devops server host
	RPCTLS           bool   `yaml:"rpcTls"`           // 是否使用 tls

	RetryMax      uint32 `yaml:"retryMax"`      // grpc 请求失败后重试次数
	RetryInterval int64  `yaml:"retryInterval"` // grpc 请求失败后重试间隔, 单位：time.Millisecond
	HoldLiveTime  int64  `yaml:"holdLiveTime"`  // grpc 连接保持存活最长时间， 单位：time.Second

	SlowThreshold int64 `yaml:"slowThreshold"` // 请求处理超过多少时间则判定为慢请求，慢请求会产生 warn 日志
	TimeOut       int64 `yaml:"timeOut"`       // 每次grpc请求最多等长时间
	ClientLogOn   bool  `yaml:"clientLogOn"`   // 发起 grpc 请求是否记录，记录内容有： 参数，与返回值
	TraceOn       bool  `yaml:"traceOn"`       // 是否开启链路追踪

	ServiceName string `yaml:"serviceName"` // 服务名, 用于服务发现, 命名格式化
}
