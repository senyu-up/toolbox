package boot

import (
	"fmt"
	"github.com/senyu-up/toolbox/combz/facade"
	"github.com/senyu-up/toolbox/example/global"
	"github.com/senyu-up/toolbox/example/internal/dao/center_service"
)

func Consumer(confPath string) (err error) {
	conf, err := bootConfig(confPath)
	if err != nil {
		return err
	}

	if conf == nil {
		return fmt.Errorf("config is nil")
	}

	tbf, err = facade.InitApp(
		facade.ConfigOptionWithApp(conf.App),
		facade.ConfigOptionWithLog(conf.Log),
		facade.ConfigOptionWithRedis(conf.Redis),
		facade.ConfigOptionWithKafka(conf.Kafka),
		facade.ConfigOptionWithMysql(conf.Mysql),
		facade.ConfigOptionWithMongo(conf.Mongo),
		facade.ConfigOptionWithTrace(conf.Trace),
		facade.ConfigOptionWithHealth(conf.Health),
		facade.ConfigOptionWithGrpcClient(conf.GrpcClient),
		facade.ConfigOptionWithAwsKafka(conf.AwsKafka),
	)
	if err != nil {
		fmt.Printf("init app facade for script failed, err: %v", err)
		return
	}
	global.SetFacade(tbf)

	// 初始化额外的 db
	if cdb, err := tbf.NewAnotherMysql(conf.CenterDb); err != nil {
		fmt.Printf("init center db failed, err: %v", err)
		return err
	} else {
		center_service.Init(cdb)
	}

	return nil
}
