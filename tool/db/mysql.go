package db

import (
	"fmt"
	"github.com/senyu-up/toolbox/tool/env"
	"net/url"
	"os"
	"time"

	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

func parseDsn(conf *config.MysqlSingleConfig) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=%s&timeout=30s",
		conf.User,
		conf.Password,
		conf.Addr,
		conf.Db,
		url.QueryEscape("UTC"))

	return dsn
}

func NewMysql(config *config.MysqlConfig) (dbInst *gorm.DB, err error) {
	dsn := config.Master.Dsn
	if dsn == "" {
		dsn = parseDsn(&config.Master)
	}

	level := glogger.Info
	slow := time.Second * 5
	if env.GetAppInfo().Stage != enum.EvnStageLocal {
		level = glogger.Error
	}

	if config.SlowThreshold > 0 {
		slow = time.Second * time.Duration(config.SlowThreshold)
	}

	if config.LogLevel > 0 {
		level = glogger.LogLevel(config.LogLevel)
	} else {
		stage := os.Getenv(enum.StageKey)
		if stage == enum.EvnStageProduction || stage == "master" {
			// 生产环境提升到warn级别
			level = glogger.Warn
		}
	}

	l := config.Logger // 如果外部传入则用外部的
	if l == nil {
		l = &GormLogger{
			level:         level,
			slowThreshold: slow,
			logDriver:     logger.GetLogger(),
		}
	}
	dbInst, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: !config.PluralTable}, Logger: l})
	if err != nil {
		return
	}

	var replicas = make([]gorm.Dialector, 0)
	if len(config.Slave) > 0 {
		for i, _ := range config.Slave {
			slaveDsn := parseDsn(&config.Slave[i])
			replicas = append(replicas, mysql.Open(slaveDsn))
		}
	}
	idleConn := compareIntWithDefault(config.MaxIdleConn, 5, 100, 20)
	openConn := compareIntWithDefault(config.MaxOpenConn, 1, 200, 10)
	idleTime := compareDurationWithDefault(config.MaxIdleTime, 10, 1800, 60)
	lifeTime := compareDurationWithDefault(config.MaxLifeTime, 30, 3600, 1800)

	err = dbInst.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(dsn)},
		Replicas: replicas,
	}).SetMaxIdleConns(idleConn).SetMaxOpenConns(openConn).SetConnMaxIdleTime(idleTime).SetConnMaxLifetime(lifeTime))

	sqlDB, err := dbInst.DB()
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()

	return
}

func compareDurationWithDefault(i time.Duration, min time.Duration, max time.Duration, dft time.Duration) (t time.Duration) {
	defer func() {
		t = t * time.Second
	}()
	if i < min {
		return dft
	} else if max > 0 && i > max {
		return dft
	}

	return i
}

func compareIntWithDefault(i int, min int, max int, dft int) int {
	if i < min {
		return dft
	} else if max > 0 && i > max {
		return dft
	}

	return i
}
