package db

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoOpts struct {
	logDriver     logger.Log // logger 接口
	ctx           context.Context
	traceOn       bool
	slowThreshold time.Duration // 慢日志阈值
	// 优先模式, 选项有： PrimaryMode (只查主), PrimaryPreferredMode (优先请求主),
	// SecondaryMode(只请求副), SecondaryPreferredMode（优先请求副）, NearestMode（主、副 哪个近请求哪个）
	preferMode readpref.Mode
}

var MongoCli *mongo.Client

func MongoDB(conf *config.MongoConfig, opts ...MgoOption) (*mongo.Client, error) {
	var opt *options.ClientOptions
	if conf.Dsn != "" {
		opt = options.Client().ApplyURI(conf.Dsn)
	} else {
		if conf.IsCluster {
			opt = options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", strings.Join(conf.Addrs, ",")))
		} else if conf.IsSrv {
			opt = options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s", conf.Addr))
		} else {
			opt = options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", conf.Addr))
		}
	}

	if conf.Password != "" {
		opt.SetAuth(options.Credential{
			AuthSource: conf.AuthSource,
			Username:   url.QueryEscape(conf.User),
			Password:   url.QueryEscape(conf.Password),
		})
	}

	// opts
	var optObj = &MongoOpts{
		traceOn:       conf.TraceOn,
		slowThreshold: time.Duration(conf.SlowThreshold) * time.Second,
	}
	for _, o := range opts {
		o(optObj)
	}

	// ctx
	var ctx = context.Background()
	if optObj.ctx != nil {
		ctx = optObj.ctx
	}

	// log driver
	var logM = MongoMonitor{
		logLevel:      int(conf.LogLevel),
		slowThreshold: optObj.slowThreshold,
		traceOn:       optObj.traceOn,
		split:         4,
		splitDuration: 10 * time.Second,

		monitorLogMap: sync.Map{},
		seqLogId:      sync.Map{},
		gcCount:       1000, // 必须大于0
	}
	if optObj.logDriver != nil {
		opt.SetMonitor(logM.GetMonitor(optObj.logDriver))
	} else if optObj.logDriver == nil && optObj.traceOn == false {
		opt.SetMonitor(GetDefaultMonitor(optObj.logDriver)) // 都不设置，则使用默认的
	} else {
		opt.SetMonitor(logM.GetMonitor(logger.GetLogger()))
	}

	// set prefer
	opt.SetReadPreference(readpref.PrimaryPreferred())
	if optObj.preferMode != 0 {
		if pref, err := readpref.New(optObj.preferMode); err == nil {
			opt.SetReadPreference(pref)
		}
	}

	cli, err := mongo.Connect(ctx, opt)
	if err != nil {
		logger.Error("connect mongo err:", err)
		return cli, err
	}
	MongoCli = cli
	return cli, err
}
