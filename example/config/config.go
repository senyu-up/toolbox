package config

import "github.com/senyu-up/toolbox/tool/config"

type Config struct {
	App *config.App `yaml:"app" json:"app"`

	Log   *config.LogConfig   `yaml:"Log" json:"Log"`
	Trace *config.TraceConfig `yaml:"Trace" json:"Trace"`

	Redis *config.RedisConfig `yaml:"redis" json:"redis"`

	Mysql           *config.MysqlConfig     `yaml:"mysql" json:"mysql"`
	AppStorageMysql *config.MysqlConfig     `yaml:"appstoragemysql" json:"appstoragemysql"`
	SqlRunner       *config.SqlRunnerConfig `yaml:"sqlrunner" json:"sqlrunner"`
	CenterDb        *config.MysqlConfig     `yaml:"centerdb" json:"centerdb"`
	Mongo           *config.MongoConfig     `yaml:"mongo" json:"mongo"`

	Kafka    *config.KafkaConfig    `yaml:"kafka" json:"kafka"`
	AwsKafka *config.AwsKafkaConfig `yaml:"aswkafka" json:"aswkafka"`

	AwsS3 *config.Aws         `yaml:"awss3" json:"awss3"`
	Email *config.EmailConfig `yaml:"email" json:"email"`

	Cron    *config.Cron          `yaml:"cron" json:"cron"`
	QwRobot *config.QwRobotConfig `yaml:"qwrobot" json:"qwrobot"`

	Fiber  *config.FiberConfig `yaml:"fiber" json:"fiber"`
	Gin    *config.GinConfig   `yaml:"gin" json:"gin"`
	Health *config.HealthCheck `yaml:"health" json:"health"`

	GrpcClient *config.GrpcClientConfig `yaml:"grpcclient" json:"grpcclient"`
	GrpcServer *config.GrpcServerConfig `yaml:"grpcserver" json:"grpcserver"`
}
