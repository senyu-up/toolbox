package env

import (
	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/config"
	"sync"
)

// 应用的基本信息
type AppInfo struct {
	HostName string // 从系统获取的主机名
	Ip       string // 从系统获取的 pod ip，我们的运行环境是 k8s，所以这里是 pod ip
	Stage    string // 运行平台，从配置获取
	Name     string // 应用名，从配置获取
	Dev      bool   // 是否是开发环境
}

func init() {
	apLock.Lock()
	defer apLock.Unlock()
	appInfo = &AppInfo{
		Name:     "None",
		Stage:    "local",
		Ip:       GetIp(),
		HostName: GetHostName(),
	}
}

var (
	appInfo *AppInfo
	apLock  = sync.RWMutex{}
)

// InitAppInfo
func InitAppInfo(opts ...AppOption) *AppInfo {
	var app = &AppInfo{
		HostName: GetHostName(),
		Ip:       GetIp(),
		Stage:    enum.EvnStageDevelop,
		Name:     "None",
	}
	for _, opt := range opts {
		opt(app)
	}
	apLock.Lock()
	defer apLock.Unlock()
	appInfo = app
	return app
}

// InitAppInfoByConf
func InitAppInfoByConf(conf *config.App) *AppInfo {
	var app = &AppInfo{
		Ip:       GetIp(),
		HostName: GetHostName(),
		Name:     conf.Name,
		Stage:    conf.Stage,
		Dev:      conf.Dev,
	}
	apLock.Lock()
	defer apLock.Unlock()
	appInfo = app
	return app
}

func GetAppInfo() *AppInfo {
	apLock.RLock()
	defer apLock.RUnlock()
	return appInfo
}
