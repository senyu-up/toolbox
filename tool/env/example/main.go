package main

import (
	"fmt"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/env"
)

func main() {
	// 初始化 app 信息
	var appConf = &config.App{
		Name:  "test",
		Stage: "local",
	}
	env.InitAppInfoByConf(appConf)

	fmt.Printf("get app info, name: %s, stage: %s, ip: %s, hostname: %s\n",
		env.GetAppInfo().Name,
		env.GetAppInfo().Stage,
		env.GetAppInfo().Ip,
		env.GetAppInfo().HostName)

	// 如果你没有通过 env.InitAppInfo 初始化app信息，可以直接调用如下方法获取 ip，hostname

	// 获取 容器/本机 ip
	env.GetIp()

	// 获取 hostname
	env.GetHostName()
}
