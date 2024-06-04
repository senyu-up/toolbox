package limiter

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/senyu-up/toolbox/tool/logger"
	"testing"
	"time"
)

var client redis.UniversalClient

func initRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 50,
	})
}

func TestLimiter(t *testing.T) {
	initRedis()
	key := "test"
	limiter := NewTokenLimiter(OptWithRate(1), OptWithBurst(1), OptWithRedisClient(client))
	for i := 0; i < 6; i++ {
		logger.Info("%d, %v", i, limiter.ReserveNow(key))
		logger.Warn("key1:", client.Get(fmt.Sprintf(tokenFormat, key)))
	}
	fmt.Println()
	time.Sleep(3 * time.Second)
	for i := 0; i < 5; i++ {
		logger.Info("%d, %v", i, limiter.ReserveNow(key))
		logger.Warn("key1:", client.Get(fmt.Sprintf(tokenFormat, key)))
	}
}
