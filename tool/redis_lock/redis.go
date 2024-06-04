package redis_lock

import "github.com/senyu-up/toolbox/tool/config"
import "github.com/go-redis/redis"

var Redis redis.UniversalClient

func InitRedisByConf(conf *config.RedisConfig) redis.UniversalClient {
	if conf.IsCluster {
		Redis = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:         conf.Addrs,
			Password:      conf.Password,
			PoolSize:      conf.PoolSize,
			MinIdleConns:  conf.MinIdleConn,
			RouteRandomly: conf.RouteRandomly,
		})
		return Redis
	} else {
		Redis = redis.NewClient(&redis.Options{
			Addr:         conf.Addrs[0],
			Password:     conf.Password,
			DB:           conf.DB,
			PoolSize:     conf.PoolSize,
			MinIdleConns: conf.MinIdleConn,
		})
		return Redis
	}
}
