package config

import (
	"time"
)

// 样本配置, 实际项目需要自行在config.go文件中定义Conf结构体
// Deprecated
type Conf struct {
	App        App
	Etcd       Etcd
	Nsq        Nsq
	Jwt        Jwt
	WeWork     WeWorkConfig
	Aws        Aws
	S3Cdn      S3Storage
	Ansible    Ansible
	ThirdParty ThirdParty
	Google     Google
	QwRobot    QwRobotConfig
	SqlRunner  SqlRunnerConfig
	Mysql      MysqlConfig
	Redis      RedisConfig
	GoCache    GoCacheConfig
	Report     Report
}

type GoCacheConfig struct {
	// 单位秒
	DefaultTTL time.Duration
	// 单位秒
	CleanInterval time.Duration
}

// 应用相关信息
type App struct {
	Name  string `yaml:"name"`  // 应用名称
	Stage string `yaml:"stage"` // 环境, 选项：local,develop,release,production
	Dev   bool   `yaml:"dev"`   // 是否为开发环境
}

type DbGroup struct {
	Mysql      DBConf `yaml:"mysql"`
	SlaveMysql DBConf `yaml:"slavemysql"`
	RedisDb    DBConf `yaml:"redisdb"`
	MongoDB    DBConf `yaml:"mongodb"`
}

type DBConf struct {
	Addr         string
	User         string
	Password     string
	PoolSize     int `yaml:"poolsize"`
	MinIdleConns int `yaml:"minidleconns"`
	Db           interface{}
	IsCluster    bool `yaml:"iscluster"`
	IsSrv        bool `yaml:"isSrv"`
	Addrs        []string
	// 前缀, 用于缓存中间件, 区分业务域
	Prefix        string
	AuthSource    string `yaml:"authSource"`
	RouteRandomly bool   `yaml:"routerandomly"`
}

type Ansible struct {
	//host
	Addrs string `yaml:"addrs"`
	//用户名
	User string `yaml:"user"`
	//密码
	Pass string `yaml:"pass"`
	//秘钥文件
	SshKeyPath string `yaml:"ssh_key_path"`
	//端口
	Port string `yaml:"port"`
}

type ThirdParty struct {
	ApiHost          string `yaml:"apiHost"`
	WebsiteHost      string `yaml:"websiteHost"`
	WebsiteWeb       string `yaml:"websiteWeb"`
	Payermax         string `yaml:"payermax"`
	FlexionPub       string `yaml:"flexionPub"`
	XsollaProductId  int    `yaml:"xsollaProductId"`
	XsollaMerchantId int    `yaml:"xsollaMerchantId"`
	XsollaApiSecret  string `yaml:"xsollaApiSecret"`
}

type Report struct {
	Secret string `yaml:"secret"`
}

type BroadcastConf struct {
	// upd 监听端口
	Port int
}

type Smtp struct {
	// Username 用户名称
	Username string `json:"username"`
	// Password 用户密码
	Password string `json:"password"`
	// Host host
	Host string `json:"host"`
}
