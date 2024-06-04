package facade

import "github.com/senyu-up/toolbox/tool/config"

type ConfigOption func(box *ToolFacade)

// Env
func ConfigOptionWithApp(app *config.App) ConfigOption {
	return func(option *ToolFacade) {
		if nil != app {
			option.configs.App = app
		}
	}
}

// health
func ConfigOptionWithHealth(health *config.HealthCheck) ConfigOption {
	return func(option *ToolFacade) {
		if nil != health {
			option.configs.health = health
		}
	}
}

// Log
func ConfigOptionWithLog(log *config.LogConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != log {
			option.configs.Log = log
		}
	}
}

// redis
func ConfigOptionWithRedis(redis *config.RedisConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != redis {
			option.configs.redis = redis
		}
	}
}

// mysql
func ConfigOptionWithMysql(mysql *config.MysqlConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != mysql {
			option.configs.mysql = mysql
		}
	}
}

// app storage mysql
func ConfigOptionWithAppStorageMysql(mysql *config.MysqlConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != mysql {
			option.configs.appStorageMysql = mysql
		}
	}
}

// mongo app storage mysql
func ConfigOptionWithMongoAppStorageMysql(mysql *config.MysqlConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != mysql {
			option.configs.mongoAppStorageMysql = mysql
		}
	}
}

// sql runner
func ConfigOptionWithSqlRunner(src *config.SqlRunnerConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != src {
			option.configs.sqlRunner = src
		}
	}
}

// mongo
func ConfigOptionWithMongo(mongo *config.MongoConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != mongo {
			option.configs.mongo = mongo
		}
	}
}

// kafka
func ConfigOptionWithKafka(kafka *config.KafkaConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != kafka {
			option.configs.kafka = kafka
		}
	}
}

// aws kafka
func ConfigOptionWithAwsKafka(kafka *config.AwsKafkaConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != kafka {
			option.configs.awsKafka = kafka
		}
	}
}

// cron
func ConfigOptionWithCron(cron *config.Cron) ConfigOption {
	return func(option *ToolFacade) {
		if nil != cron {
			option.configs.cron = cron
		}
	}
}

// aws S3
func ConfigOptionWithAwsS3(awsS3 *config.Aws) ConfigOption {
	return func(option *ToolFacade) {
		if nil != awsS3 {
			option.configs.awsS3 = awsS3
		}
	}
}

// email
func ConfigOptionWithEmail(email *config.EmailConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != email {
			option.configs.email = email
		}
	}
}

// Trace
func ConfigOptionWithTrace(trace *config.TraceConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != trace {
			option.configs.Trace = trace
		}
	}
}

// weWork
func ConfigOptionWithWeWork(weWork *config.WeWorkConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != weWork {
			option.configs.weWork = weWork
		}
	}
}

// qwRobot
func ConfigOptionWithQwRobot(qwRobot *config.QwRobotConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != qwRobot {
			option.configs.qwRobot = qwRobot
		}
	}
}

// fiber http server
func ConfigOptionWithFiber(conf *config.FiberConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != conf {
			option.configs.fiber = conf
		}
	}
}

// ConfigOptionWithGin gin http server
func ConfigOptionWithGin(conf *config.GinConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != conf {
			option.configs.gin = conf
		}
	}
}

// grpc client
func ConfigOptionWithGrpcClient(clientConf *config.GrpcClientConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != clientConf {
			option.configs.grpcClient = clientConf
		}
	}
}

// grpc server
func ConfigOptionWithGoGrpcServer(serverConf *config.GrpcServerConfig) ConfigOption {
	return func(option *ToolFacade) {
		if nil != serverConf {
			option.configs.grpcServer = serverConf
		}
	}
}
