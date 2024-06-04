package config

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"time"
)

type SqlRunnerConfig struct {
	LockKey   string `yaml:"lockKey,omitempty"`   // SQL 锁的 key，默认使用文件的 hash
	RunPolicy string `yaml:"runPolicy,omitempty"` // SQL 执行策略，可选值为 once 或 always，默认为 always
	SqlFile   string `yaml:"sqlFile,omitempty"`   // SQL 文件路径，可选，基于 stage 自动取

	Stage      string                `yaml:"_"` // 当前环境
	LockTTL    time.Duration         `yaml:"-"` // SQL 锁的 TTL，单位秒，默认为 1 分钟
	RedisCli   redis.UniversalClient `yaml:"-"` // 必填, Redis 客户端
	DBInstList []*gorm.DB            `yaml:"-"` // 必填, GORM 数据库实例列表
}
