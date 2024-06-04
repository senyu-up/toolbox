package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	athena2 "github.com/aws/aws-sdk-go/service/athena"
	"github.com/senyu-up/toolbox/tool/logger"
	athenadriver "github.com/uber/athenadriver/go"
)

// -------------------------------
var athenas sync.Map

type AthenaConfigs struct {
	LogPath   string
	Region    string
	AccessID  string
	AccessKey string
	Database  string
}

type Options func(cfg *AthenaConfigs)

func WithLogPath(path string) Options {
	return func(cfg *AthenaConfigs) {
		cfg.LogPath = path
	}
}

func WithRegion(region string) Options {
	return func(cfg *AthenaConfigs) {
		cfg.Region = region
	}
}

func WithAccessID(id string) Options {
	return func(cfg *AthenaConfigs) {
		cfg.AccessID = id
	}
}

func WithAccessKey(key string) Options {
	return func(cfg *AthenaConfigs) {
		cfg.AccessKey = key
	}
}

func WithDatabase(database string) Options {
	return func(cfg *AthenaConfigs) {
		cfg.Database = database
	}
}

func margeConfig(opts ...Options) *AthenaConfigs {
	cfg := &AthenaConfigs{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// InitAthena 初始化Athena数据库链接
func InitAthena(opts ...Options) (*sql.DB, error) {
	cfg := margeConfig(opts...)
	conf, err := athenadriver.NewDefaultConfig(
		cfg.LogPath,
		cfg.Region,
		cfg.AccessID,
		cfg.AccessKey)
	if err != nil {
		return nil, err
	}
	if cfg.Database != "" {
		conf.SetDB(cfg.Database)
	}
	// 目前只有data-operation和dw-agent在使用，设置这俩个参数来使空值返回
	// 的是对应类型的零值而不是统一返回空字符串
	conf.SetMissingAsEmptyString(false)
	conf.SetMissingAsDefault(true)

	athenaIns, err := sql.Open(athenadriver.DriverName, conf.Stringify())
	if err != nil {
		return nil, err
	}
	return athenaIns, nil
}

// InitAthenaClient 初始化Athena Client
func InitAthenaClient(opts ...Options) (*athena2.Athena, error) {
	cfg := margeConfig(opts...)
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessID, cfg.AccessKey, ""),
	})
	if err != nil {
		return nil, err
	}
	return athena2.New(sess), nil
}

func clientKey(appKey string) string {
	return fmt.Sprintf("%s.Client", appKey)
}

func StorageAthenaClient(appKey string, ops ...Options) error {
	cli, err := InitAthenaClient(ops...)
	if err != nil {
		return err
	}
	athenas.Store(clientKey(appKey), cli)
	return nil
}

// LoadAthenaClient 在内存中获取AthenaClient
// 如果内存中不存在,且Athena链接配置不为空,将会初始化链接并存储
// 如果内存中存在,且链接不为空的情况下,将内存中的Client返回
func LoadAthenaClient(appKey string, ops ...Options) (*athena2.Athena, error) {
	var err error
	v, ok := athenas.Load(clientKey(appKey))
	if !ok && len(ops) != 0 {
		v, err = InitAthena(ops...)
		if err != nil {
			logger.Error("*************init athena fail ! ! ! **************")
			return nil, err
		}
	}
	if v == nil {
		logger.Error("load athena database is nil")
		return nil, errors.New("load athena database is nil")
	}
	athenas.Store(clientKey(appKey), v)
	return v.(*athena2.Athena), nil
}

func DelAthenaClient(appKey string) {
	athenas.Delete(clientKey(appKey))
}

// StorageAthenaDatabase 初始化Athena数据库链接,并将数据库链接存储至内存
func StorageAthenaDatabase(appKey string, ops ...Options) error {
	db, err := InitAthena(ops...)
	if err != nil {
		return err
	}
	athenas.Store(appKey, db)
	return nil
}

// LoadAthenaDatabase 在内存中获取AthenaDatabase
// 如果内存中不存在,且Athena链接配置不为空,将会初始化链接并存储
// 如果内存中存在,且链接不为空的情况下,将内存中的Database返回
func LoadAthenaDatabase(appKey string, ops ...Options) (*sql.DB, error) {
	var err error
	v, ok := athenas.Load(appKey)
	if !ok && len(ops) != 0 {
		v, err = InitAthena(ops...)
		if err != nil {
			logger.Error("*************init athena fail ! ! ! **************")
			return nil, err
		}
	}
	if v == nil {
		logger.Error("load athena database is nil")
		return nil, errors.New("load athena database is nil")
	}
	athenas.Store(appKey, v)
	return v.(*sql.DB), nil
}

func DelAthenaDatabase(appKey string) {
	athena, exist := athenas.LoadAndDelete(appKey)
	if exist {
		athena.(*sql.DB).Close()
	}
}
