package config

import (
	"github.com/go-redis/redis"
)

type Broadcast struct {
	Redis redis.UniversalClient
	Topic string
}
