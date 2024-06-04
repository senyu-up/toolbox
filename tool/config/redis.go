package config

type RedisConfig struct {
	Addrs    []string `yaml:"addrs"`    // Redis 集群地址列表，如果是单点，就填一个，会去第一个数组值作为地址
	User     string   `yaml:"user"`     // Redis 访问用户名
	Password string   `yaml:"password"` // Redis 访问密码
	//deprecated
	DB            int    `yaml:"db"`               // Redis 数据库编号（非 cluster 模式下生效, 测试、线上环境的redis都是cluster，所以这个配置实际上没用）
	Prefix        string `yaml:"prefix,omitempty"` // Redis 存储 key 的前缀
	IsCluster     bool   `yaml:"isCluster"`        // 是否为 Redis 集群模式，
	PoolSize      int    `yaml:"poolSize"`         // Redis 连接池大小
	RouteRandomly bool   `yaml:"routeRandomly"`    // 在集群模式下是否随机路由
	MinIdleConn   int    `yaml:"minIdleConn"`      // 最小空闲连接数
}
