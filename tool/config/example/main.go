package main

import (
	"fmt"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/config/loader"
	"github.com/spf13/cast"
)

func main() {
	conf, err := loader.InitConf(&loader.File{},
		loader.ConfOptWithPath("./config.yaml"),
		loader.ConfOptWithType("yaml"))
	if err != nil {
		return
	}
	type AppConf struct {
		Etcd  config.Etcd
		Mysql config.MysqlConfig
		Redis config.RedisConfig
	}

	var app = AppConf{}
	if err = conf.Unmarshal(&app); err != nil {
		return
	}

	dsn, err := conf.Get("mysql.dsn")
	if err != nil {
		return
	}
	fmt.Printf("get mysql dsn %s", cast.ToString(dsn))
}
