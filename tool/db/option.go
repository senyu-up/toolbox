package db

import (
	"context"
	"github.com/senyu-up/toolbox/tool/logger"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	glogger "gorm.io/gorm/logger"
	"time"
)

type SqlLogOption func(*GormLogger)

func GormLogOptWithTraceOn(on bool) SqlLogOption {
	return func(option *GormLogger) {
		option.traceOn = on
	}
}

func GormLogOptWithLevel(level glogger.LogLevel) SqlLogOption {
	return func(option *GormLogger) {
		if 0 < level {
			option.level = level
		}
	}
}

func GormLogOptWithLogDriver(log logger.Log) SqlLogOption {
	return func(option *GormLogger) {
		if log != nil {
			option.logDriver = log
		}
	}
}

func GormLogOptWithSlowThreshold(slow time.Duration) SqlLogOption {
	return func(option *GormLogger) {
		if 0 < slow {
			option.slowThreshold = slow
		}
	}
}

func GormLogOptWithIgnoreRecordNotFoundError(ignore bool) SqlLogOption {
	return func(option *GormLogger) {
		option.ignoreRecordNotFoundError = ignore
	}
}

type MgoOption func(*MongoOpts)

func MgoOptWithLogDriver(log logger.Log) MgoOption {
	return func(option *MongoOpts) {
		if log != nil {
			option.logDriver = log
		}
	}
}

func MgoOptWithTraceOn(on bool) MgoOption {
	return func(option *MongoOpts) {
		option.traceOn = on
	}
}

func MgoOptWithContext(c context.Context) MgoOption {
	return func(option *MongoOpts) {
		option.ctx = c
	}
}

func MgoOptWithSlowThreshold(slow time.Duration) MgoOption {
	return func(option *MongoOpts) {
		if 0 < slow {
			option.slowThreshold = slow
		}
	}
}

func MgoOptWithPref(pref readpref.Mode) MgoOption {
	return func(option *MongoOpts) {
		if 0 < pref {
			option.preferMode = pref
		}
	}
}
