package loader

import (
	"errors"
)

var (
	ErrConfigEmpty     = errors.New("Config Empty ")
	ErrConfigKeyNotSet = errors.New("Config key not set ")
)

type Config struct {
	loader Loader

	configType string // 配置类型
	path       string // 配置路径, 可以是filepath，或者 uri

}

type ConfVal struct {
	Val interface{}
}

type Loader interface {
	Init(...ConfOption) error            // 初始化配置信息
	Get(key string) (interface{}, error) // 传入路径，获取对应配置，
	Unmarshal(dst interface{}) error     // 把配置信息反序列化到传入的结构体上
}

type ConfOption func(*Config)

func ConfOptWithType(confType string) ConfOption {
	return func(option *Config) {
		option.configType = confType
	}
}

func ConfOptWithPath(p string) ConfOption {
	return func(option *Config) {
		option.path = p
	}
}

// 初始化配置对象，同时调用 loader.init 方法
func InitConf(loader Loader, param ...ConfOption) (Config, error) {
	return Config{loader: loader}, loader.Init(param...)
}

// 把配置序列化到传入结构体
// 如果结构体和配置有差异，会尽量赋值，不会报错！
func (c *Config) Unmarshal(dst interface{}) error {
	return c.loader.Unmarshal(dst)
}

func (c *Config) Get(key string) (interface{}, error) {
	return c.loader.Get(key)
}
