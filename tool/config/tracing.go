package config

const (
	ClientLogLevelInfo  = "info"
	ClientLogLevelError = "error"
)

type TraceConfig struct {
	ServerLogOn    bool       `yaml:"serverLogOn"`              // 是否开启链路请求日志
	ClientLogOn    bool       `yaml:"clientLogOn"`              // 是否开启客户端请求日志
	ClientLogLevel string     `yaml:"clientLogLevel,omitempty"` // 客户端日志输出等级 info => 所有请求都会进行打印, error => 仅对错误响应进行打印
	Jaeger         JaegerConf `yaml:"jaeger"`                   // Jaeger 配置
}
