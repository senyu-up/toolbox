package env

import (
	"github.com/senyu-up/toolbox/tool/geoip"
	"sync"
)

var (
	ip   string = ""
	lock        = sync.RWMutex{}
)

// GetIp
//
//	@Description: 获取本地Ip
//	@return string IP字符
func GetIp() string {
	if 0 < len(ip) {
		lock.RLock()
		defer lock.RUnlock()
		return ip
	}

	lock.Lock()
	defer lock.Unlock()
	ip, _ = geoip.GetLocalIPV4()
	return ip
}
