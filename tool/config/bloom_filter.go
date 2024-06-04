package config

import (
	"github.com/go-redis/redis"
	"time"
)

type RedisBloomFilterConf struct {
	// 误判率, 如: 0.00000001
	//1% error rate requires 7 hash functions and 10.08 bits per item.
	//0.1% error rate requires 10 hash functions and 14.4 bits per item.
	//0.01% error rate requires 14 hash functions and 20.16 bits per item
	//https://redis.io/commands/bf.add/
	P float64
	// 容量, -1 => 弹性伸缩, 0 => 2 << 31, >0 => Cap
	Cap uint
	// 缓存key名
	Key string
	// redis 实例
	Cache redis.UniversalClient
	// if 0 => no exipire, else expire time
	Ttl time.Duration
}

type BitsetBloomConf struct {
	// redis的key
	Key string
	// redis cli 实例
	Cache redis.UniversalClient
	// 过期时间, 可不填, 不填即不过期
	Ttl time.Duration
	// 容量
	Cap uint
	// 误判率, 如: 0.00000001
	//1% error rate requires 7 hash functions and 10.08 bits per item.
	//0.1% error rate requires 10 hash functions and 14.4 bits per item.
	//0.01% error rate requires 14 hash functions and 20.16 bits per item
	//细节参考: https://redis.io/commands/bf.add/
	P float64
}

type GoBloomFilterConf struct {
	// 容错率
	P float64
	// 容量
	Cap uint
}
