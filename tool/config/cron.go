package config

// CronConfig cronv2 配置
type Cron struct {
	TraceOn        bool `json:"traceon"`        // 开启链路追踪的打点，关闭不影响traceId生成
	AccurateSecond bool `json:"accuratesecond"` // 表达式 精确到秒
}
