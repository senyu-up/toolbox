package config

import "github.com/go-redis/redis"

type LeakyBucketConf struct {
	// redis key name
	Key string
	//
}

// StickyWindowConf
// @Description: 固定窗口配置
type StickyWindowConf struct {
	// 必填, redis key
	Key string
	// 必填, 每个周期产生的令牌数
	NumPerPeriod uint64
	// 可选,周期, 单位秒, 默认1秒
	Period uint32
	// 必填, redis cli
	Redis redis.UniversalClient
}

// TokenBucketConf
// @Description: 令牌桶的配置
type TokenBucketConf struct {
	// 必填, redis key
	Key string
	// 必填, 令牌桶的容量
	Capacity uint64
	// 必填, 每个周期产生的令牌数
	NumPerPeriod uint64
	// 可选,周期, 单位秒, 默认1秒
	Period uint32
	// 必填, redis cli
	Redis redis.UniversalClient
}
