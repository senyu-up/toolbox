package appstorage

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/logger"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

// StoreOption options
type StoreOption func(box *DBStorage)

// StoreOptionWithChannelID nsq topic channel
func StoreOptionWithChannelID(channelID string) StoreOption {
	return func(option *DBStorage) {
		if "" != channelID {
			option.channelID = channelID
		}
	}
}

// StoreOptionWithAppStage Stage
func StoreOptionWithAppStage(appStage string) StoreOption {
	return func(option *DBStorage) {
		if "" != appStage {
			option.stage = appStage
		}
	}
}

// StoreOptionWithGorm gorm db handler
func StoreOptionWithGorm(db *gorm.DB) StoreOption {
	return func(option *DBStorage) {
		if nil != db {
			option.db = db
		}
	}
}

// StoreOptionWithRedisCli set redis cli
func StoreOptionWithRedisCli(cli *redis.UniversalClient) StoreOption {
	return func(option *DBStorage) {
		if nil != cli {
			option.redisCli = cli
		}
	}
}

// StoreOptionWithLog log
func StoreOptionWithLog(log logger.Log) StoreOption {
	return func(option *DBStorage) {
		if nil != log {
			option.log = log
		}
	}
}

// StoreOptionWithLogLevel log level
func StoreOptionWithLogLevel(level int) StoreOption {
	return func(option *DBStorage) {
		if level > 0 {
			option.logLevel = glogger.LogLevel(level)
		}
	}
}

// StoreOptionWithAppName app name
func StoreOptionWithAppName(appName string) StoreOption {
	return func(option *DBStorage) {
		if "" != appName {
			option.appName = appName
		}
	}
}

// StoreOptionWithMySqlConf set default mysql conf
func StoreOptionWithMySqlConf(conf *config.MysqlConfig) StoreOption {
	return func(option *DBStorage) {
		if nil != conf { // 非空则覆盖
			option.mysqlConf = conf
		}
	}
}

// StoreOptionWithContext set context
func StoreOptionWithContext(c context.Context) StoreOption {
	return func(option *DBStorage) {
		if c != nil {
			option.ctx = c
		}
	}
}

// StoreOptionWithInit set context
func StoreOptionWithInit(b bool) StoreOption {
	return func(option *DBStorage) {
		option.initImmediately = b
	}
}
