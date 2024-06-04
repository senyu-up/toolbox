package env

import "os"

// GetHostName
//
//	@Description: 获取 hostname
//	@return string
func GetHostName() string {
	hn, _ := os.Hostname()
	return hn
}
