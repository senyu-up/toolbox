package main

import (
	"context"
	"fmt"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/db"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"time"
)

func NewMysqlWithLogger() {
	// 设置底层 logger driver
	//logger.SwitchLogger(logger.AdapterConsole)
	logger.SwitchLogger(logger.AdapterZap)
	var logDriver = logger.GetLogger()

	// 初始化 mysqlLogger
	var mysqlLogger = db.NewGormLogger(
		db.GormLogOptWithLevel(glogger.Info),
		db.GormLogOptWithTraceOn(true),
		db.GormLogOptWithLogDriver(logDriver),
		db.GormLogOptWithSlowThreshold(time.Second),
	)

	// 初始化 mysql conf
	var mysqlConf = &config.MysqlConfig{
		Master: config.MysqlSingleConfig{
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "12345678",
			Db:       "test",
		},
		Logger: mysqlLogger,
	}

	// 初始化 mysql
	dbClient, err := db.NewMysql(mysqlConf)
	if err != nil {
		logger.Error("new mysql err %v", err)
		return
	}

	// 一个trace ctx
	var ctx = trace.NewTrace()

	// query
	if rows, err := dbClient.WithContext(ctx).Select("*").Table("test").Where("1=1").Rows(); err != nil {
		logger.Error("select * from test err %v", err)
		return
	} else {
		cols, _ := rows.Columns()
		logger.Info("get cols ", cols)
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
				logger.Error("scan row err %v", err)
			} else {
				logger.Info("row: ", rowMap, id, name, sex)
			}
		}
	}
}

type Person struct {
	Id   uint   `gorm:"primary_key,column:id"`
	Name string `gorm:"column:name"`
	Sex  int    `gorm:"sex"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Person 表名
func (Person) TableName() string {
	return "person"
}

func Example1() {
	// 初始化 mysql conf
	var mysqlConf = &config.MysqlConfig{
		Master: config.MysqlSingleConfig{
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "123456",
			Db:       "test",
		},
		//Logger: mysqlLogger,
	}

	// 初始化 mysql
	dbClient, err := db.NewMysql(mysqlConf)
	if err != nil {
		fmt.Printf("new mysql err %v", err)
		return
	}

	// 以 person model 去查询数据
	var p = &Person{}
	if err := dbClient.Model(Person{}).Where("id = ?", 1).Find(&p).Error; err != nil {
		fmt.Printf("find err %v \n", err)
	} else {
		fmt.Printf("find success %+v\n", p)
	}
}

func SetMysqlLogger() {
	// 设置底层 logger driver
	logger.SwitchLogger(logger.AdapterZap)
	var logDriver = logger.GetLogger()

	// 初始化 mysqlLogger
	var mysqlLogger = db.NewGormLogger(
		db.GormLogOptWithLevel(glogger.Info),
		db.GormLogOptWithTraceOn(true),
		db.GormLogOptWithLogDriver(logDriver), // 设置 mysql logger 的 driver
		db.GormLogOptWithSlowThreshold(time.Second),
	)

	// 初始化 mysql conf
	var mysqlConf = &config.MysqlConfig{
		Master: config.MysqlSingleConfig{
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "123456",
			Db:       "test",
		},
		Logger: mysqlLogger, // 指定 mysql logger
	}

	// 初始化 mysql
	dbClient, err := db.NewMysql(mysqlConf)
	if err != nil {
		logger.Error("new mysql err %v", err)
		return
	}

	// 一个trace ctx
	var ctx = trace.NewTrace()

	// query, 把带有 trace 信息的 ctx 传入
	if rows, err := dbClient.WithContext(ctx).Select("*").Table("test").Where("1=1").Rows(); err != nil {
		logger.Error("select * from test err %v", err)
		logger.SetErr(err).Error("select * from test err 2 %v", err)
		logger.GetLogger().Error("select * from test err 3 %v", err)
		return
	} else {
		logger.Info("select success, get rows %v \n", rows)
	}
}

type User struct {
	Name      string    `bson:"name"`
	Age       uint32    `bson:"age"`
	Weight    float32   `bson:"weight"`
	Studying  bool      `bson:"studying"`
	Tag       []string  `bson:"tag"`
	CreatedAt time.Time `bson:"created_at"`
}

func InitMongo() {
	ctx, _ := context.WithTimeout(trace.NewTrace(), 5*time.Second) //创建一个带有trace信息，且5s超时的 context
	// 初始化 mongo conf
	var mongoConf = &config.MongoConfig{
		Dsn:        "mongodb://admin:123456@127.0.0.1:27017/?authSource=admin&authMechanism=SCRAM-SHA-1", // local
		Addr:       "127.0.0.1:27017",
		User:       "admin",
		Password:   "123456",
		AuthSource: "admin",
		IsSrv:      false,
	}
	mClient, err := db.MongoDB(mongoConf, db.MgoOptWithTraceOn(true), db.MgoOptWithContext(ctx), db.MgoOptWithLogDriver(logger.GetLogger()))
	if err != nil {
		logger.Error("init mongo err %v", err)
		return
	}
	if err = mClient.Ping(ctx, nil); err != nil {
		logger.Error("ping mongo err %v", err)
		return
	}

	// 查询条件
	var filter = bson.M{"name": "tom"}
	var user = &User{}
	// 指定 db ，collection
	var coll = mClient.Database("test").Collection("person")
	// 查询一行数据
	if coll.FindOne(ctx, filter).Decode(user); err != nil {
		logger.Error("find one err %v", err)
		return
	} else {
		logger.Info("find one success")
	}
}

func DBResolver(db *gorm.DB) {
	// 获取默认的db连接
	var defaultInst = db
	// 指定为写库， sql 请求会走主库
	var WriteInst = db.Clauses(dbresolver.Write)
	// 指定为读库，sql 请求会走从库
	var ReadInst = db.Clauses(dbresolver.Read)

	defaultInst.Exec("select * from test")
	WriteInst.Exec("select * from test")
	ReadInst.Exec("select * from test")
}

func NewMongo() *mongo.Client {
	var cnf = config.MongoConfig{
		Dsn: "mongodb://root:gmdqw1yXdMKV3KBvzLbvdj4GEtJyPFe1@172.16.10.86:27017/?authSource=admin&authMechanism=SCRAM-SHA-1", // dev
	}
	var ctx = context.Background()
	mgoDb, err := db.MongoDB(&cnf, db.MgoOptWithTraceOn(true), db.MgoOptWithContext(ctx), db.MgoOptWithLogDriver(logger.GetLogger()))
	if err != nil {
		logger.Error("new mongo err %v", err)
		return nil
	}
	return mgoDb
}

func TestMongoMigrate() {
	//var ctx = context.Background()
	//// 获取 mongo client
	//var mgo = NewMongo()
	//
	//// 声明一个 collection index，并制定 db
	//var dbName = "center_service2"
	//var collIdx = db.NewCollectionIndex(mgo.Database(dbName))
	//collIdx.RegisterTable(TestLog{})
	//if errMsg := collIdx.Migrate(ctx); errMsg != "" {
	//	logger.Error("migrate err %v", errMsg)
	//}
}

func main() {
	//Example1()

	//SetMysqlLogger()

	//InitMongo()

	//TestMongoMigrate()
	NewMysqlWithLogger()
}
