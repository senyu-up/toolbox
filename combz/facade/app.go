package facade

import (
	"context"
	"errors"
	"fmt"
	"github.com/senyu-up/toolbox/tool/http/gin_server"
	rc "github.com/senyu-up/toolbox/tool/redis_lock"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"

	"github.com/senyu-up/toolbox/combz/appstorage"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/cronv2"
	"github.com/senyu-up/toolbox/tool/db"
	"github.com/senyu-up/toolbox/tool/email"
	"github.com/senyu-up/toolbox/tool/env"
	"github.com/senyu-up/toolbox/tool/http/fiber"
	"github.com/senyu-up/toolbox/tool/http/http_health"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/mq/aws_kafka"
	"github.com/senyu-up/toolbox/tool/mq/kafka"
	"github.com/senyu-up/toolbox/tool/storage"
	"github.com/senyu-up/toolbox/tool/trace"
	"github.com/senyu-up/toolbox/tool/wework"
	"github.com/senyu-up/toolbox/tool/wework/qwrobot"
)

type ShutdownFunc func(context.Context)

type Config struct {
	App *config.App `yaml:"app" json:"app"`

	Log   *config.LogConfig   `yaml:"Log" json:"Log"`
	Trace *config.TraceConfig `yaml:"Trace" json:"Trace"`

	redis *config.RedisConfig `yaml:"redis" json:"redis"`

	mysql                *config.MysqlConfig     `yaml:"mysql"           json:"mysql"`
	appStorageMysql      *config.MysqlConfig     `yaml:"appstoragemysql" json:"appstoragemysql"`
	mongoAppStorageMysql *config.MysqlConfig     `yaml:"mongoappstoragemysql" json:"mongoappstoragemysql"`
	sqlRunner            *config.SqlRunnerConfig `yaml:"sqlrunner"       json:"sqlrunner"`
	mongo                *config.MongoConfig     `yaml:"mongo"           json:"mongo"`

	kafka    *config.KafkaConfig    `yaml:"kafka" json:"kafka"`
	awsKafka *config.AwsKafkaConfig `yaml:"aws_kafka" json:"aws_kafka"`
	cron     *config.Cron           `yaml:"cron" json:"cron"`

	awsS3 *config.Aws         `yaml:"awss3" json:"awss3"`
	email *config.EmailConfig `yaml:"email" json:"email"`

	weWork  *config.WeWorkConfig  `yaml:"wework" json:"wework"`
	qwRobot *config.QwRobotConfig `yaml:"qwrobot" json:"qwrobot"`

	fiber  *config.FiberConfig `yaml:"fiber" json:"fiber"`
	gin    *config.GinConfig   `yaml:"gin" json:"gin"`
	health *config.HealthCheck `yaml:"health" json:"health"`

	grpcClient *config.GrpcClientConfig `yaml:"grpcclient" json:"grpcclient"`
	grpcServer *config.GrpcServerConfig `yaml:"grpcserver" json:"grpcserver"`
}

type ToolFacade struct {
	configs Config

	env *env.AppInfo

	logger        *logger.Log // logger
	gormLogDriver *logger.Log // gorm log driver
	traceClient   *trace.SuTracer

	redisClient redis.UniversalClient

	mysqlClient *gorm.DB
	gormLogger  *db.GormLogger // gorm logger
	mongoClient *mongo.Client

	healthChecker *http_health.HealthChecker

	kafkaClient    *kafka.Kafka
	awsKafkaClient *aws_kafka.Kafka
	cronClient     *cronv2.Client

	emailSes      *email.AwsSes
	emailPinPoint *email.AwsPinPoint
	awsS3         *storage.S3Conn

	qwRobot *qwrobot.QWRobot
	weWork  *wework.WechatClient

	// 基础组件
	// ------------------
	// 业务组件
	fiber *fiber.App
	gin   *gin_server.App

	appStorage      *appstorage.DBStorage
	mongoAppStorage *appstorage.MongoStorage

	shutdown []ShutdownFunc
}

func InitApp(opts ...ConfigOption) (tb *ToolFacade, err error) {
	tb = &ToolFacade{}
	for _, opt := range opts {
		opt(tb)
	}
	// app Info
	if tb.configs.App != nil {
		tb.env = env.InitAppInfoByConf(tb.configs.App)
	} else {
		tb.env = env.InitAppInfo()
	}

	// init Log
	if tb.configs.Log != nil {
		tb.configs.Log.AppName = tb.env.Name
		logger.InitDefaultLoggerByConf(tb.configs.Log)
		if tb.qwRobot != nil {
			// 如果 企微机器人初始化了，则设置该 机器人
			logger.SetCallBack(logger.NewQwRobot(tb.qwRobot))
		}
		var log = logger.GetLogger()
		tb.logger = &log
	}

	// init cache
	// redis
	if tb.configs.redis != nil {
		tb.redisClient = rc.InitRedisByConf(tb.configs.redis)
	}

	// qw robot
	if tb.configs.qwRobot != nil {
		if tb.redisClient == nil {
			return tb, errors.New("QWRobot depends on redis, redis not init, please pass redis config ")
		} else {
			tb.configs.qwRobot.RedisCli = tb.redisClient
		}
		if tb.qwRobot, err = qwrobot.Init(tb.configs.qwRobot,
			qwrobot.OptWithHostName(tb.env.HostName),
			qwrobot.OptWithIp(tb.env.Ip),
			qwrobot.OptWithStage(tb.env.Stage)); err != nil {
			return tb, err
		} else {
			// 微信 机器人 初始化成功，设置回调
			logger.SetCallBack(logger.NewQwRobot(tb.qwRobot))
			var log = logger.GetLogger()
			tb.logger = &log
		}
	}

	// wework
	if tb.configs.weWork != nil {
		if tb.redisClient == nil {
			return tb, errors.New("WeWorkConfig depends on redis, redis not init, please pass redis config ")
		}
		if tb.weWork, err = wework.InitByConfig(tb.configs.weWork,
			wework.OptWithHttpClient(&http.Client{})); err != nil {
			return tb, err
		}
	}

	// init Trace
	if tb.configs.Trace != nil {
		tb.configs.Trace.Jaeger.AppName = tb.env.Name
		trace.Init(tb.configs.Trace) //  初始化并设置全局 trace
		tb.traceClient = trace.Tracer
	}

	if tb.configs.Log != nil && tb.configs.mysql != nil {
		if tb.configs.mysql.CallDepth == 0 {
			tb.configs.mysql.CallDepth = 3
		}
		// init gorm log driver
		if l, err := tb.getLoggerWithCallerSkip(tb.configs.Log.CallDepth + tb.configs.mysql.CallDepth); err != nil {
			return tb, err
		} else {
			tb.gormLogDriver = &l
		}
	}

	// init mysql
	if tb.configs.mysql != nil {
		// 初始化 mysqlLogger
		tb.gormLogger = db.NewGormLogger(
			db.GormLogOptWithLevel(glogger.Info),
			db.GormLogOptWithTraceOn(true),
			db.GormLogOptWithLogDriver(tb.GetGormLogDriver()),
			db.GormLogOptWithSlowThreshold(time.Second*time.Duration(tb.configs.mysql.SlowThreshold)),
		)
		tb.configs.mysql.Logger = tb.GetGormLogger()
		db, err := db.NewMysql(tb.configs.mysql)
		if err != nil {
			return tb, err
		}
		tb.mysqlClient = db
	}

	// mongo
	if tb.configs.mongo != nil {
		if tb.mongoClient, err = db.MongoDB(tb.configs.mongo); err != nil {
			return tb, err
		}
	}

	// kafka
	if tb.configs.kafka != nil {
		if tb.kafkaClient, err = kafka.New(tb.configs.kafka); err != nil {
			return tb, err
		}
	}
	// aws kafka
	if tb.configs.awsKafka != nil {
		if tb.awsKafkaClient, err = aws_kafka.New(context.Background(), tb.configs.awsKafka); err != nil {
			return tb, err
		}
	}

	// cron
	if tb.configs.cron != nil {
		tb.cronClient = cronv2.New(cronv2.CronOptionWithSecond(tb.configs.cron.AccurateSecond),
			cronv2.CronOptionWithTrace(tb.configs.cron.TraceOn),
			cronv2.CronOptionWithLogger(cronv2.NewCronLogger(*tb.GetLogger())))
	}

	// init email aws ses， pinpoint
	if tb.configs.email != nil {
		if tb.emailSes, err = email.InitAwsSes(tb.configs.email); err != nil {
			return tb, err
		}
		if tb.emailPinPoint, err = email.InitAwsPinPoint(tb.configs.email); err != nil {
			return tb, err
		}
	}

	// init s3
	if tb.configs.awsS3 != nil {
		if tb.awsS3 = storage.InitByConf(tb.configs.awsS3); err != nil {
			return tb, err
		}
	}

	// init fiber http server
	if tb.configs.fiber != nil {
		tb.configs.fiber.Name = env.GetAppInfo().Name
		if tb.fiber, err = fiber.NewApp(tb.configs.fiber); err != nil {
			return tb, err
		}
	}
	//init gin http server
	if tb.configs.gin != nil {
		tb.configs.gin.Name = env.GetAppInfo().Name
		if tb.gin, err = gin_server.NewApp(tb.configs.gin); err != nil {
			return tb, err
		}
	}

	// 初始化 app storage
	if tb.configs.appStorageMysql != nil && tb.redisClient != nil {
		if tb.configs.appStorageMysql.CallDepth == 0 {
			tb.configs.appStorageMysql.CallDepth = 3 // 如果没指定则定为3
		}
		log, err := tb.getLoggerWithCallerSkip(tb.configs.Log.CallDepth + tb.configs.appStorageMysql.CallDepth)
		if err != nil {
			return tb, err
		}
		appDb, err := tb.NewAnotherMysql(tb.configs.appStorageMysql)
		if err != nil {
			return tb, err
		}
		if tb.appStorage, err = appstorage.NewAppStorageDB(
			appstorage.StoreOptionWithLog(log),
			appstorage.StoreOptionWithGorm(appDb),
			appstorage.StoreOptionWithRedisCli(&tb.redisClient),
			appstorage.StoreOptionWithLogLevel(tb.configs.appStorageMysql.LogLevel),
			appstorage.StoreOptionWithChannelID(env.GetAppInfo().Name),
			appstorage.StoreOptionWithAppName(env.GetAppInfo().Name),
			appstorage.StoreOptionWithMySqlConf(tb.configs.appStorageMysql),
			appstorage.StoreOptionWithAppStage(env.GetAppInfo().Stage)); err != nil {
			return tb, err
		}
	}

	// 初始化 mongo app storage
	if tb.configs.mongoAppStorageMysql != nil && tb.redisClient != nil {
		if tb.configs.mongoAppStorageMysql.CallDepth == 0 {
			tb.configs.mongoAppStorageMysql.CallDepth = 3 // 如果没指定则定为3
		}
		log, err := tb.getLoggerWithCallerSkip(tb.configs.Log.CallDepth + tb.configs.appStorageMysql.CallDepth)
		if err != nil {
			return tb, err
		}
		appDb, err := tb.NewAnotherMysql(tb.configs.appStorageMysql)
		if err != nil {
			return tb, err
		}
		if tb.mongoAppStorage, err = appstorage.NewMongoStorageDB(
			appstorage.StoreOptionWithLog(log),
			appstorage.StoreOptionWithGorm(appDb),
			appstorage.StoreOptionWithRedisCli(&tb.redisClient),
			appstorage.StoreOptionWithLogLevel(tb.configs.appStorageMysql.LogLevel),
			appstorage.StoreOptionWithChannelID(env.GetAppInfo().Name),
			appstorage.StoreOptionWithAppName(env.GetAppInfo().Name),
			appstorage.StoreOptionWithAppStage(env.GetAppInfo().Stage)); err != nil {
			return tb, err
		}
	}

	//初始化 grpc client
	//if tb.configs.grpcClient != nil && tb.configs.App != nil {
	//	tb.grpcClientMan = serverhelper.InitGlobalConMan(tb.configs.grpcClient,
	//		serverhelper.ClientOptWithStage(env.GetAppInfo().Stage), // 运行环境
	//	)
	//}

	// 初始化 grpc server
	//if tb.configs.grpcServer != nil && tb.configs.App != nil {
	//	tb.configs.grpcServer.TraceOn = tb.traceClient != nil
	//	if tb.grpcServer, err = serverhelper.NewXHGrpc(tb.configs.grpcServer); err != nil {
	//		return tb, err
	//	}
	//}

	// health check
	if tb.configs.health != nil {
		if tb.healthChecker, err = http_health.NewHttpHealthCheckServer(
			http_health.HealthOptionWithDisableLog(tb.configs.health.DisableLog),
			http_health.HealthOptionWithPprof(tb.configs.health.Pprof),
			http_health.HealthOptionWithAddr(tb.configs.health.Addr),
			http_health.HealthOptionWithPort(tb.configs.health.Port)); err != nil {
			return tb, err
		}
	}
	return tb, err
}

func (a *ToolFacade) getLoggerWithCallerSkip(skip int) (logger.Log, error) {
	// 初始化 gormLogger
	var logConf = *a.configs.Log
	logConf.CallDepth = skip // caller skip
	return logger.InitLoggerByConf(&logConf)
}

// NewAnotherMysql
//
//	@Description: 初始化另一个 mysql client, 用于连接其他数据库
//	@receiver a
//	@param conf  body any true "-"
//	@return *gorm.DB
//	@return error
func (a *ToolFacade) NewAnotherMysql(conf *config.MysqlConfig) (*gorm.DB, error) {
	if conf == nil {
		return nil, errors.New("mysql config is nil")
	}
	if conf.CallDepth == 0 {
		conf.CallDepth = 3
	}
	// 没有则创建
	l, err := a.getLoggerWithCallerSkip(logger.DefaultCallDepth + conf.CallDepth)
	if err != nil {
		return nil, err
	}
	// 初始化 mysqlLogger
	var mysqlLogger = db.NewGormLogger(
		db.GormLogOptWithLevel(glogger.Info),
		db.GormLogOptWithTraceOn(true),
		db.GormLogOptWithLogDriver(l),
		db.GormLogOptWithSlowThreshold(time.Second*time.Duration(conf.SlowThreshold)),
	)
	conf.Logger = mysqlLogger
	return db.NewMysql(conf)
}

// NewAnotherRedis
//
//	@Description: 初始化另一个 Redis client, 用于连接其他Redis
//	@receiver a
//	@param conf  body any true "-"
//	@return *redis.UniversalClient
//	@return error
func (a *ToolFacade) NewAnotherRedis(conf *config.RedisConfig) (redis.UniversalClient, error) {
	if conf == nil {
		return nil, errors.New("mysql config is nil")
	}

	return rc.InitRedisByConf(conf), nil
}

// GetEnv
//
//	@Description: 获取当前环境信息，如果没有初始化，则返回空的，要求业务方必须设置 应用信息
//	@receiver a
//	@return *env.AppInfo
func (a *ToolFacade) GetEnv() *env.AppInfo {
	return a.env
}

// GetLogger
//
//	@Description: 获取日志对象，如果为空，则获取默认 logger
//	@receiver a
//	@return *logger.Log
func (a *ToolFacade) GetLogger() *logger.Log {
	if a.logger == nil {
		var l = logger.GetLogger()
		return &l
	}
	return a.logger
}

// GetGormLogDriver
//
//	@Description: 获取 gorm log driver，如果为空，则获取默认 logger
//	@receiver a
//	@return logger.Log
func (a *ToolFacade) GetGormLogDriver() logger.Log {
	if a.gormLogDriver == nil {
		// 没有则创建
		if l, err := a.getLoggerWithCallerSkip(logger.DefaultCallDepth + 3); err == nil {
			a.gormLogDriver = &l
			return l
		} else if a.logger != nil {
			return *a.logger
		}
	} else {
		return *a.gormLogDriver
	}
	return logger.GetLogger()
}

// GetTraceClient
//
//	@Description: 获取 trace client， 如果为空，则表示未设置 trace，其他组件无需开启 trace
//	@receiver a
//	@return *trace.SuTracer
func (a *ToolFacade) GetTraceClient() *trace.SuTracer {
	return a.traceClient
}

func (a *ToolFacade) GetRedisClient() redis.UniversalClient {
	return a.redisClient
}

func (a *ToolFacade) GetMysqlClient() *gorm.DB {
	return a.mysqlClient
}

func (a *ToolFacade) GetGormLogger() *db.GormLogger {
	return a.gormLogger
}

func (a *ToolFacade) GetMongoClient() *mongo.Client {
	return a.mongoClient
}

func (a *ToolFacade) GetKafkaClient() *kafka.Kafka {
	return a.kafkaClient
}

func (a *ToolFacade) GetAwsKafkaClient() *aws_kafka.Kafka {
	return a.awsKafkaClient
}

func (a *ToolFacade) GetCronClient() *cronv2.Client {
	return a.cronClient
}

func (a *ToolFacade) GetEmailSes() *email.AwsSes {
	return a.emailSes
}

func (a *ToolFacade) GetEmailPinPoint() *email.AwsPinPoint {
	return a.emailPinPoint
}

func (a *ToolFacade) GetAwsS3() *storage.S3Conn {
	return a.awsS3
}

func (a *ToolFacade) GetQwRobot() *qwrobot.QWRobot {
	return a.qwRobot
}

func (a *ToolFacade) GetWeWork() *wework.WechatClient {
	return a.weWork
}

func (a *ToolFacade) GetFiber() *fiber.App {
	return a.fiber
}
func (a *ToolFacade) GetGin() *gin_server.App {
	return a.gin
}

func (a *ToolFacade) GetHealthChecker() *http_health.HealthChecker {
	return a.healthChecker
}

//func (a *ToolFacade) GetGrpcServer() *serverhelper.XHGrpc {
//	return a.grpcServer
//}

//func (a *ToolFacade) GetGrpcClientMan() *serverhelper.ClientMan {
//	return a.grpcClientMan
//}

func (a *ToolFacade) GetAppStorage() *appstorage.DBStorage {
	return a.appStorage
}

func (a *ToolFacade) GetMongoAppStorage() *appstorage.MongoStorage {
	return a.mongoAppStorage
}

func (a *ToolFacade) StartFiber() error {
	if a.fiber != nil {
		a.shutdown = append(a.shutdown, func(ctx context.Context) {
			a.fiber.Fiber().Shutdown()
		})
		return a.fiber.Run()
	}
	return fmt.Errorf("fiber is nil, please init fiber first")
}

func (a *ToolFacade) StartGin() error {
	if a.gin != nil {
		return a.gin.Run()
	}
	return fmt.Errorf("fiber is nil, please init fiber first")
}

//func (a *ToolFacade) StartGrpc() error {
//	if a.grpcServer != nil {
//		a.shutdown = append(a.shutdown, func(ctx context.Context) {
//			a.grpcServer.Stop()
//		})
//		return a.grpcServer.Run()
//	}
//	return fmt.Errorf("grpc server is nil, please init grpc server first")
//}

func (a *ToolFacade) StartHealthChecker() error {
	if a.healthChecker != nil {
		a.shutdown = append(a.shutdown, func(ctx context.Context) {
			a.healthChecker.Close()
		})
		return a.healthChecker.Start()
	}
	return fmt.Errorf("health checker is nil, please init health checker first")
}

// StartKafkaConsumeAsync
//
//	@Description: 启动 kafka 异步
//	@receiver a
//	@return error
func (a *ToolFacade) StartKafkaConsumeAsync() error {
	if a.kafkaClient != nil {
		a.shutdown = append(a.shutdown, func(ctx context.Context) {
			a.kafkaClient.Close()
		})
		a.kafkaClient.StartConsume()
		return nil
	}
	return fmt.Errorf("kafka client is nil, please init kafka client first")
}

// StartAwsKafkaConsumeAsync
//
//	@Description: 启动 aws kafka 异步
//	@receiver a
//	@return error
func (a *ToolFacade) StartAwsKafkaConsumeAsync() error {
	if a.awsKafkaClient != nil {
		a.shutdown = append(a.shutdown, func(ctx context.Context) {
			a.awsKafkaClient.Close()
		})
		a.awsKafkaClient.StartConsume()
		return nil
	}
	return fmt.Errorf("aws kafka client is nil, please init aws kafka client first")
}

// StartCronAsync
//
//	@Description: 启动 cron 异步
//	@receiver a
//	@return error
func (a *ToolFacade) StartCronAsync() error {
	if a.cronClient != nil {
		a.shutdown = append(a.shutdown, func(ctx context.Context) {
			a.cronClient.Stop()
		})
		a.cronClient.Start()
		return nil
	}
	return fmt.Errorf("cron client is nil, please init cron client first")
}

// 关闭所有启动的服务
func (a *ToolFacade) Shutdown(ctx context.Context) {
	for _, shutdown := range a.shutdown {
		shutdown(ctx)
	}
}
