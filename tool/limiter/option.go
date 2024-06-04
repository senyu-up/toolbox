package limiter

import (
	"github.com/go-redis/redis"
)

type LimitOption func(*TokenLimiter)

func OptWithRedisClient(red redis.UniversalClient) LimitOption {
	return func(obj *TokenLimiter) {
		obj.store = red
		obj.redisAlive = 1 // 设置了redis，则认为redis是健康的
	}
}

// burst 令牌桶容量
func OptWithBurst(burst int) LimitOption {
	return func(obj *TokenLimiter) {
		if burst > 0 {
			obj.burst = burst
		}
	}
}

// rate 令牌桶生产速率
func OptWithRate(rate int) LimitOption {
	return func(obj *TokenLimiter) {
		if rate > 0 {
			obj.rate = rate
		}
	}
}
