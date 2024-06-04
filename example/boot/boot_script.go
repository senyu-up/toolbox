package boot

import (
	"fmt"
	"github.com/senyu-up/toolbox/combz/facade"
	"github.com/senyu-up/toolbox/example/global"
)

// 初始化 cron_job, 一次性脚本，都用这个初始化
func Script(confPath string) (err error) {
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
		facade.ConfigOptionWithMysql(conf.Mysql),
		facade.ConfigOptionWithKafka(conf.Kafka),
		facade.ConfigOptionWithMongo(conf.Mongo),
		facade.ConfigOptionWithTrace(conf.Trace),
		facade.ConfigOptionWithHealth(conf.Health),
		facade.ConfigOptionWithCron(conf.Cron),
		facade.ConfigOptionWithGrpcClient(conf.GrpcClient),
	)
	if err != nil {
		fmt.Printf("init app facade for script failed, err: %v", err)
		return
	}
	global.SetFacade(tbf)
	return nil
}
