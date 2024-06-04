package bloom_filter

import (
	"github.com/senyu-up/toolbox/tool/bloom_filter/driver"
	"github.com/senyu-up/toolbox/tool/config"
)

type BloomInterface interface {
	// 判断key是否存在
	Exists(key string) (exists bool, err error)
	// 批量判断key是否存在, 返回值通知索引下标关联key是否存在
	MExists(keys []string) (statusList []bool, err error)
	// 添加key, 如果key存在返回false, 反之为true
	Add(key string) (ok bool, err error)
	// 批量添加key, 返回值通知索引下标关联key是否存在
	MAdd(keys []string) (statusList []bool, err error)
	// 查看当前布隆的状态, 比如, 容量, key 数量, ...
	Info() (info map[string]interface{}, err error)
	// 插入一个key, 如果key存在返回false, 反之为true, 与Add类似, 区别在redis会判断当前布隆是否存在, 不存在会主动创建
	Insert(key string) (ok bool, err error)
	// 批量一个key, 如果key存在返回false, 反之为true, 与Add类似, 区别在redis会判断当前布隆是否存在, 不存在会主动创建
	MInsert(keys []string) (statusList []bool, err error)
}

func NewRedisBloom(conf *config.RedisBloomFilterConf) (BloomInterface, error) {
	return driver.NewRedisBloom(conf)
}

func NewGoBloom(conf *config.GoBloomFilterConf) (BloomInterface, error) {
	return driver.NewGoBloomFilter(conf)
}

func NewBitsetBloom(conf *config.BitsetBloomConf) (BloomInterface, error) {
	return driver.NewBitsetBloom(conf)
}
