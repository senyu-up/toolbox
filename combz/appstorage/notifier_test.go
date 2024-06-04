package appstorage

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"log"
	"os"
	"testing"
	"time"
)

var appDb *DBStorage
var redisCli redis.UniversalClient
var opDsn = "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=UTC&timeout=10s"
var opAppKey = "MlB0ZWQyYTRhYToxNjc4NTA0Njc0OmRldmVsb3A="

func SetUp() {
	if appDb != nil {
		return
	}
	var redConf = config.RedisConfig{
		Addrs: []string{"127.0.0.1:6379"},
	}
	var err error
	redisCli = cache.InitRedisByConf(&redConf)

	var dbConf = config.MysqlConfig{
		Master: config.MysqlSingleConfig{
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "123456",
			Db:       "test",
		}}
	dbCli, err := db.NewMysql(&dbConf)
	if err != nil {
		fmt.Printf("conn mysql err %v \n", err)
		return
	}

	appDb, err = NewAppStorageDB(StoreOptionWithGorm(dbCli),
		StoreOptionWithRedisCli(&redisCli),
		StoreOptionWithChannelID("xh_test"),
		StoreOptionWithAppStage(enum.EvnStageLocal))
	if err != nil {
		fmt.Printf("new app storage  err %v \n", err)
		return
	}

	RegisterAddEvent(func(appKey string) error {
		fmt.Printf("test func call , add event %s", appkey)
		return nil
	})

	RegisterUpdateEvent(func(appKey string) error {
		fmt.Printf("test func call , update event %s", appkey)
		return nil
	})

	RegisterDelEvent(func(appKey string) error {
		fmt.Printf("test func call , del event %s", appkey)
		return nil
	})

	time.Sleep(time.Second)
}

func AddDataBySql() {
	var myCli = mysql.Open(opDsn)
	db, err := gorm.Open(myCli, &gorm.Config{Logger: glogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer,
		glogger.Config{
			SlowThreshold: 5 * time.Second, //慢sql阀值
			Colorful:      true,
			LogLevel:      glogger.Info,
		})})
	if err != nil {
		fmt.Printf("new gorm  err %v", err)
		return
	}
	var sql = "INSERT INTO `app_dsns` (`id`, `app_key`, `dsn`, `dsn_slave`, `mongo_dsn`, `created_at`, `updated_at`, `xxx_name`, `is_delete`, `gateway`, `app_secret`, `icon`, `app_id`) VALUES " +
		" (200,'MlB0ZWQyYTRhYToxNjc4NTA0Njc0OmRldmVsb3A=','root:3x3O5QHmcNFiNSEeQ7kCSZBdsZfJU3kd@tcp(172.16.10.86:3306)/GP0102?charset=utf8mb4&parseTime=True&loc=UTC&timeout=10s','root:3x3O5QHmcNFiNSEeQ7kCSZBdsZfJU3kd@tcp(172.16.10.86:3306)/GP0102?charset=utf8mb4&parseTime=True&loc=UTC&timeout=10s', 'mongodb://root:gmdqw1yXdMKV3KBvzLbvdj4GEtJyPFe1@172.16.10.86:27017/?authMechanism=SCRAM-SHA-1', 1678504674,1678504674,'GP0102',0,'','279fbbb9566d06c90fe5566da7b57434','https://cdn6.platform-im-test.allxxx.com/2023/03/11/1f5424bae3e445748a3cbd0c10c621e9.jpg',2013)"
	if err = db.Exec(sql).Error; err != nil {
		fmt.Printf("exec sql err %v", err)
		return
	}

}

func DelDataBySql() {
	db, err := gorm.Open(mysql.Open(opDsn), &gorm.Config{Logger: glogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer,
		glogger.Config{
			SlowThreshold: 5 * time.Second, //慢sql阀值
			Colorful:      true,
			LogLevel:      glogger.Info,
		})})
	if err != nil {
		fmt.Printf("new gorm  err %v", err)
		return
	}
	var sql = "delete from `app_dsns` where `app_key` = \"" + opAppKey + "\""
	if err = db.Exec(sql).Error; err != nil {
		fmt.Printf("exec sql err %v", err)
		return
	}
}

func TestPushAppDSNChangeNotification(t *testing.T) {
	SetUp()
	type args struct {
		redisCli redis.UniversalClient
		msg      *DSNNotifyMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1", args: args{redisCli: redisCli, msg: &DSNNotifyMessage{
				AppKey:   opAppKey,
				Category: AddDataIns,
			}},
		},
		{
			name: "2", args: args{redisCli: redisCli, msg: &DSNNotifyMessage{
				AppKey:   opAppKey,
				Category: UpdateDataIns,
			}},
		},
		{
			name: "3", args: args{redisCli: redisCli, msg: &DSNNotifyMessage{
				AppKey:   opAppKey,
				Category: RemoveDataIns,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.msg.Category == AddDataIns {
				AddDataBySql()
			}
			if tt.args.msg.Category == RemoveDataIns {
				DelDataBySql()
				if gorm := appDb.GetDB(opAppKey); gorm != nil {
					if db, err := gorm.DB(); err == nil && db != nil {
						db.Close()
					}
				}
				if gorm := appDb.GetReadDB(opAppKey); gorm != nil {
					if db, err := gorm.DB(); err == nil && db != nil {
						db.Close()
					}
				}
			}
			if err := PushAppDSNChangeNotification(tt.args.redisCli, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("PushAppDSNChangeNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
			time.Sleep(time.Second)
			// 检查数量
			if gorm := appDb.GetDB(opAppKey); gorm == nil && (tt.args.msg.Category == AddDataIns || tt.args.msg.Category == UpdateDataIns) {
				t.Errorf("add dsn [NOT STORE] in map")
			} else if gorm != nil && tt.args.msg.Category == RemoveDataIns {
				t.Errorf("remove dsn [NOT DELETE] from map")
			}
		})
	}
}

func TestPushMongoAppDSNChangeNotification(t *testing.T) {
	SetUpForMongo()
	var ctx = context.Background()
	type args struct {
		redisCli redis.UniversalClient
		msg      *DSNNotifyMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1", args: args{redisCli: redisCli, msg: &DSNNotifyMessage{
				AppKey:   opAppKey,
				Category: AddDataIns,
			}},
		},
		{
			name: "2", args: args{redisCli: redisCli, msg: &DSNNotifyMessage{
				AppKey:   opAppKey,
				Category: UpdateDataIns,
			}},
		},
		{
			name: "3", args: args{redisCli: redisCli, msg: &DSNNotifyMessage{
				AppKey:   opAppKey,
				Category: RemoveDataIns,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.msg.Category == AddDataIns {
				AddDataBySql()
			}
			if tt.args.msg.Category == RemoveDataIns {
				DelDataBySql()
				if mgoCli := appMgoDb.GetWriteDB(ctx, opAppKey); mgoCli != nil {
					mgoCli.Client().Disconnect(ctx)
				}
				if mgoCli := appMgoDb.GetReadDB(ctx, opAppKey); mgoCli != nil {
					mgoCli.Client().Disconnect(ctx)
				}
			}
			if err := PushMongoAppDSNChangeNotification(tt.args.redisCli, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("PushMongoAppDSNChangeNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
			time.Sleep(time.Second)
			// 检查数量
			if mgoCli := appMgoDb.GetWriteDB(ctx, opAppKey); mgoCli == nil && (tt.args.msg.Category == AddDataIns || tt.args.msg.Category == UpdateDataIns) {
				t.Errorf("add mongo dsn [NOT STORE] in map")
			} else if mgoCli != nil && tt.args.msg.Category == RemoveDataIns {
				t.Errorf("remove mongo dsn [NOT DELETE] from map")
			}
		})
	}
}
