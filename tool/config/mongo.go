package config

type MongoConfig struct {
	Dsn        string   `yaml:"dsn"`                  // 数据库连接DSN, 如果不为空优先用 dsn 而忽略其他配置
	Addr       string   `yaml:"addr"`                 // 数据库连接地址
	Addrs      []string `yaml:"addrs"`                // 数据库集群地址列表, iscluster 模式下用这个地址
	User       string   `yaml:"user"`                 // 数据库用户名
	Password   string   `yaml:"password"`             // 数据库密码
	Db         string   `yaml:"db"`                   // 默认Database
	AuthSource string   `yaml:"authSource,omitempty"` // 认证库名，默认为 admin

	// 慢查询阈值, 单位秒, 默认5
	SlowThreshold int64
	// 4=>info(非master,production默认) 3=>Warn(master,production默认) 2=>Error 1=>Silent
	LogLevel int64
	// 是否开启 联路追踪
	TraceOn bool

	IsCluster bool `yaml:"isCluster"` // 是否为集群
	IsSrv     bool `yaml:"isSrv"`     // 是否使用 SRV 模式连接
}
