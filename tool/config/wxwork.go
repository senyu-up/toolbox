package config

import "github.com/go-redis/redis"

type WeWorkConfig struct {
	CorpId  string `yaml:"corpId"`
	Secret  string `yaml:"secret"`
	AgentId string `yaml:"agentId"`

	//deprecated
	RefreshInterval uint32 `yaml:"refreshInterval"` // 刷新 token 的间隔时间， 单位秒
	TryTimes        uint32 `yaml:"TryTimes"`        // 请求失败了，重试的次数
	Debug           bool   `yaml:"debug"`           // 在调用企微接口时，是否带上 debug 参数
}

type QwRobotConfig struct {
	// 机器人地址
	Webhook string `yaml:"webhook"`
	// 常规频率限制, 支持 n/S n/M n/H
	InfoFreqLimit string `yaml:"infoFreqLimit"`
	// 警告消息频率限制, 支持 n/S n/M n/H
	WarnFreqLimit string `yaml:"warnFreqLimit"`
	// 错误消息频率限制, 支持 n/S n/M n/H
	ErrorFreqLimit string `yaml:"errorFreqLimit"`
	// 消息类型, 自定义, 区分消息类型
	MessageType string `yaml:"messageType"`
	// Redis key前缀, 区分业务域,非必填,
	Prefix string `yaml:"prefix"`

	RedisCli redis.UniversalClient `yaml:"-"` // 可用的 redis 连接
}
