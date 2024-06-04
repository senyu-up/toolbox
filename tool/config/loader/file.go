package loader

import (
	"github.com/spf13/viper"
	"path"
	"strings"
)

type File struct {
	driver *viper.Viper
}

// 接受最多两个参数
// 第一个参数指定文件路径，
// 第二个参数指定配置类型，例如：yaml，toml，json，如果不传第二个，则通过文件拓展名自动判断
func (f *File) Init(opts ...ConfOption) error {
	var (
		c = &Config{
			path: "config.yaml",
		}
	)
	for _, opt := range opts {
		opt(c)
	}
	if 1 > len(c.configType) {
		c.configType = strings.Trim(path.Ext(c.path), ".") // 获取的ext带点！
	}

	f.driver = viper.New()
	f.driver.SetConfigType(c.configType)
	f.driver.SetConfigFile(c.path)

	return f.driver.ReadInConfig()
}

func (f *File) Get(key string) (interface{}, error) {
	var err error = nil
	if !f.driver.IsSet(key) {
		err = ErrConfigKeyNotSet
	}
	return f.driver.Get(key), err
}

func (f *File) Unmarshal(dst interface{}) error {
	keys := f.driver.AllKeys()
	if 1 > len(keys) {
		// 如果配置为空，报错
		return ErrConfigEmpty
	}
	return f.driver.Unmarshal(dst)
}
