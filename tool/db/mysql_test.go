package db

import (
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/trace"
	glogger "gorm.io/gorm/logger"
	"testing"
	"time"
)

func TestNewMysql(t *testing.T) {

	// 设置底层 logger dirver
	//logger.SwitchLogger(logger.AdapterConsole)
	logger.SwitchLogger(logger.AdapterZap)
	var logDriver = logger.GetLogger()

	// 初始化 mysqlLogger
	var mysqlLogger = NewGormLogger(
		GormLogOptWithLevel(glogger.Info),
		GormLogOptWithTraceOn(true),
		GormLogOptWithLogDriver(logDriver),
		GormLogOptWithSlowThreshold(time.Second),
	)

	// 初始化 mysql conf
	var mysqlConf = &config.MysqlConfig{
		Master: config.MysqlSingleConfig{
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "123456",
			Db:       "test",
		},
		Logger: mysqlLogger,
	}

	// 初始化 mysql
	db, err := NewMysql(mysqlConf)
	if err != nil {
		t.Errorf("new mysql err %v", err)
		return
	}

	// 一个trace ctx
	var ctx = trace.NewTrace()

	// query
	if rows, err := db.WithContext(ctx).Select("*").Table("test").Where("1=1").Rows(); err != nil {
		t.Errorf("select * from test err %v", err)
		return
	} else {
		cols, _ := rows.Columns()
		t.Log("get cols ", cols)
		type user struct {
			Id   int    `gorm:"column:id"`
			Name string `gorm:"column:name"`
			Sex  int    `gorm:"column:sex"`
		}

		if rows.Next() {
			var (
				rowMap = []user{}
				id     int
				name   string
				sex    int
			)

			if err = rows.Scan(&id, &name, &sex); err != nil {
				t.Errorf("scan row err %v", err)
			} else {
				t.Log("row: ", rowMap, id, name, sex)
			}
		}
	}
}
