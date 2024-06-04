# toolbox.config 配置

## 使用

```go
package main

import (
	"fmt"

	"github.com/spf13/cast"
	"toolbox/tool/config"
	"toolbox/tool/config/loader"
)


func main() () {
    // 用 文件loader 加载配置，配置信息： ./config.yml,  配置格式 yaml
    conf, err := loader.InitConf(&loader.File{}, "./config.yml", "yaml")
    // 用 文件loader 加载配置，配置信息： ./config.yaml, 配置格式会从文件拓展名自动判断
    conf, err = loader.InitConf(&loader.File{}, "./config.yaml")
    // 或者不传入参数，自动使用默认配置
    conf, err = loader.InitConf(&loader.File{})
    if err != nil {
        return
    }

    // 读取某个变量
    re, err := conf.Get("app.dev")
    if err != nil {
        return
    } else {
        fmt.Printf("get config val %v by app.dev", cast.ToString(re))
    }

	// toolbox 不再提供 config.Conf 配置结构体，各个项目按照自己所需，引入个组件的配置
	// 例如项目中，使用 Etcd, Kafka, Mysql, Redis 组件，那么在项目中创建一个 MyConf 结构体就能满足需要，如下：
	type MyConf struct {
		App        config.App
		Etcd       config.Etcd
		Kafka      config.Kafka
		Mysql      config.MysqlConfig
		Redis      config.RedisConfig
	}

	// 配置序列化到结构体
	var appConf = MyConf{}
	if err := conf.Unmarshal(&appConf); err == nil {
		fmt.Printf("get config obj %+v", appConf)
	}
}
```