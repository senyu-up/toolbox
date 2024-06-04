package config

type KafkaConfig struct {
	Brokers          []string `yaml:"brokers"`            // kafka集群地址列表
	Timeout          int      `yaml:"timeout"`            // 发送消息超时时间, 单位秒
	Level            int      `yaml:"level"`              // 消息等级 1. 允许消息出现丢失, leader确认收到消息即可, 性能较高 2. 不允许消息丢失, 所有broker均确认收到消息
	User             string   `yaml:"user,omitempty"`     // [可选] 用户名
	Password         string   `yaml:"password,omitempty"` // [可选] 密码
	SyncFullMetadata bool     `yaml:"syncFullMetadata"`   // 同步所有主题的metadata信息, 默认为 false
	DialTimeout      int      `yaml:"dialTimeout"`        // Dial网络时间配置, 单位秒, connection 超时时间, 默认30s
	ReadTimeout      int      `yaml:"readTimeout"`        // Read 时间, 单位秒, 默认30s
	WriteTimeout     int      `yaml:"writeTimeout"`       // Write 时间, 单位秒, 默认30s
	Version          string   `yaml:"version"`            // Kafka 版本，不传则使用：2.6.0.0

	MetadataTimeout               int `yaml:"metadataTimeout"`               // 获取metadata超时时间, 单位秒, 默认30s
	ProducerTimeout               int `yaml:"producerTimeout"`               // 投递超时时间, 单位秒
	ProducerMaxRetryTimes         int `yaml:"producerMaxRetryTimes"`         // 最大重新投递次数
	MetadataMaxRetryTimes         int `yaml:"metadataMaxRetryTimes"`         // 获取metadata信息最大重试次数, 默认3次
	MetadataRefreshIntervalSecond int `yaml:"metadataRefreshIntervalSecond"` // 元数据信息刷新间隔, 单位秒, 默认10分钟

	// 多个consumer配置
	Consumers []*KafkaConsumerConfig `yaml:"consumers"`

	// consumer 每个 topic 起多少个协程去处理消息
	Workers int `yaml:"workers"`
	// consumer 是否从最老的开始消费
	Oldest bool `yaml:"oldest"`

	SASL struct {
		Enable   bool   `yaml:"enable"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"sasl"` // consumer sasl 配置

	TraceOn bool `yaml:"traceOn"` // 是否开启 trace, 打开后会产生一层 span
}

type KafkaConsumerConfig struct {
	// broker的集群地址
	Brokers []string `yaml:"brokers"`
	Version string   `yaml:"version"` // Kafka 版本，不传则使用：2.6.0.0
	// topic 的名称
	Topic string `yaml:"topic"`
	// 消费组名称
	Group string `yaml:"group"`
	SASL  struct {
		Enable   bool   `yaml:"enable"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"sasl"`
	// 多少个协程
	Workers int `yaml:"workers"`
	// 是否从最老的开始消费
	Oldest bool `yaml:"oldest"`
}

type AwsKafkaConfig struct {
	SASL struct {
		// Whether or not to use SASL authentication when connecting to the broker
		// (defaults to false).
		Enable bool
		// SASLMechanism is the name of the enabled SASL mechanism.
		// Possible values: OAUTHBEARER, PLAIN (defaults to PLAIN).
		// option: PLAIN, AWS_MSK_IAM, SCRAM-SHA-512, SCRAM-SHA-256
		Mechanism string

		Region    string `yaml:"region"`    // 区域
		AccessId  string `yaml:"accessId"`  // 访问ID
		SecretKey string `yaml:"secretKey"` // 访问密钥

		UserName string `yaml:"userName"` // 用户名
		Password string `yaml:"password"` // 密码
	}

	Brokers []string `yaml:"brokers"` // kafka集群地址列表
	Timeout int      `yaml:"timeout"` // 消息处理超时时间,包含读写, 单位秒

	// 发送方，负载均衡策略, 默认为 Hash，可选：hash,least_bytes,round_robin,crc32balancer,murmur2_balancer,reference_hash
	Balancer string `yaml:"balancer"`

	TraceOn bool `yaml:"traceOn"` // 是否开启 trace, 打开后会产生一层 span
	Workers int  `yaml:"workers"` // 消费端多少个协程并发消费

	MaxAttempts    int   `yaml:"maxAttempts"`    // 最大重试次数
	CommitInterval int   `yaml:"commitInterval"` // flushes commits to Kafka every second,
	MaxBytes       int64 `yaml:"maxBytes"`       // 消息体最大字节数

	BackoffMin int `yaml:"backoffMin"` // 最小重试间隔时间, 单位 ms， 默认100ms
	BackoffMax int `yaml:"backoffMax"` // 最大重试间隔时间, 单位 ms, 默认1s

	// 消息等级0. 发了就不管， 不等待分区确认 1. 允许消息出现丢失, leader确认收到消息即可, 性能较高 2. 不允许消息丢失, 所有broker均确认收到消息
	RequiredAcks int8 `yaml:"requiredAcks"`

	IsolationLevel int8 `yaml:"isolationLevel"` // 事务隔离级别, 可选：0-ReadUncommitted, 1-ReadCommitted

	// 是否允许自动创建topic, 亚马逊kafka服务默认不允许自动创建topic
	AllowAutoTopicCreation bool `yaml:"allowAutoTopicCreation"`

	GroupId string `yaml:"-"` // 消费组名称
	Topic   string `yaml:"-"` // topic 的名称
	Async   bool   `yaml:"-"` // 是否异步发送
}
