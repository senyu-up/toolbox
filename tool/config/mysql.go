package config

import (
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

// MysqlConfig
// @Description: mysql 连接配置, 替代DbGroup配置, 封装读写分离逻辑
type MysqlConfig struct {
	Dsn         string
	MaxOpenConn int
	MaxIdleConn int
	// 最大空闲时间, 单位秒
	MaxIdleTime time.Duration
	// 最大生命周期, 单位秒
	MaxLifeTime time.Duration
	Master      MysqlSingleConfig
	// 从节点配置
	Slave []MysqlSingleConfig
	// 慢查询阈值, 单位秒, 默认5
	SlowThreshold int
	// 4=>info(非master,production默认) 3=>Warn(master,production默认) 2=>Error 1=>Silent. 详情看gorm.io/gorm/logger/logger.go
	LogLevel int
	// 结构体转表名是否为复数，填 true 会在结构体后+s作为表名，中台标准是 false
	PluralTable bool
	// gorm logger 调用栈跳过层数，不填则默认3
	CallDepth int

	Logger glogger.Interface
}

type MysqlSingleConfig struct {
	// 可选
	Dsn      string
	User     string
	Password string
	Addr     string
	Db       string
}

type MysqlPageConfig struct {
	// 可选, 默认为200
	MaxPageSize int
	// 可选, 最大分页数, 0 表示不限制
	MaxPage int
	// 可选, 默认 "id"
	PrimaryKey string
	// 可选, 默认 "total", 用于统计记录数的字段
	TotalKey string
	// 必选
	DB *gorm.DB
}
