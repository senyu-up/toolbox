package config

type JaegerConf struct {
	JaegerOn bool   `yaml:"jaegerOn"` // 是否开启 Jaeger
	AppName  string `yaml:"-"`        // 应用名

	CollectorEndpoint string `yaml:"collectorEndpoint,omitempty"` // 收集器地址
	AgentPort         string `yaml:"agentPort,omitempty"`         // 本地 agent 端口，默认为 8888
	User              string `yaml:"user,omitempty"`              // 用户名
	Password          string `yaml:"password,omitempty"`          // 密码

	SamplerFreq              float64                `yaml:"samplerFreq,omitempty"`              // 采样频率，取值范围 (0.0 and 1.0]，默认为 1（每次都进行采集），该参数与 RateLimitPerSecond 二选一，RateLimitPerSecond 优先级更高
	RateLimitPerSecond       float64                `yaml:"rateLimitPerSecond,omitempty"`       // 每秒采集限制，默认不限制。如果指定了该值，则忽略采样频率
	Tags                     map[string]interface{} `yaml:"tags,omitempty"`                     // 服务标签
	QueueSize                int                    `yaml:"queueSize,omitempty"`                // 队列大小，默认为 100
	QueueFlushIntervalSecond int                    `yaml:"queueFlushIntervalSecond,omitempty"` // 缓冲刷新间隔，默认为 5 秒
}

type Jaeger struct {
	Sampler  JaegerSampler
	Reporter JaegerReporter
}

type JaegerSampler struct {
	Type  string `yaml:"type"`
	Param float64
}
type JaegerReporter struct {
	AgentAddr string `yaml:"agentAddr"`
}
