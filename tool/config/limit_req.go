package config

import "github.com/go-redis/redis"

type RouteSource struct {
	// 资源的参数值
	Source string
	// 对应的路由
	Route string
}

type LimiterController string

const (
	//  日志
	LimiterControllerLogger LimiterController = "log"
	// 机器人告警
	LimiterControllerAlert LimiterController = "alert"
)

type LimiterType string

const (
	// 分布式限流, 采用 令牌桶 算法进行限流, 对精度控制较高, 依赖redis
	LimiterTypeDistributed LimiterType = "distributed"
	// 单机限流, 采用 滑动时间窗口 算法进行限流, 高性能, 推荐
	LimiterTypeLocal LimiterType = "local"
	LimiterTypeAli               = "ali"
)

type LimiterGroup struct {
	Name string
	Uris []string
}

type LimiterRule struct {
	Group []LimiterGroup
	// 最大请求数  1/s  1/m  1/h
	Burst string
	// 指定获取资源的key, 默认获取当前的uri
	SourceKey string
	// @description 当触发限流时如何处理
	Controller LimiterController
}

type Trust struct {
	// 子网掩码
	SubnetMask []string
}

type Limiter struct {
	AppName string // 限流器名字
	Trust   Trust
	Rules   []LimiterRule
	// 当 Category = LimiterTypeDistributed 时需要
	Redis redis.UniversalClient
	// distributed => 分布式    local => 单机
	Category LimiterType
	// 指定如何获取数据
	Taker func(input interface{}, key string) string
}

// 令牌桶限流器 tool/limit
type TokenLimiter struct {
	// 令牌桶容量
	Burst int
	// 令牌桶生产速率
	Rate int
}
