package driver

import (
	"errors"
	"github.com/senyu-up/toolbox/tool/logger"
)

// ErrRecordNotFound
var (
	// 定义的key已存在
	ErrItemAlreadyExists = errors.New("ERR item exists")
)

func getErrorRatio(p float64) float64 {
	if p > 0.01 {
		p = 0.001
	} else if p <= 0 {
		p = 0.001
		logger.Warn("RedisBloomFilter", "init", "p is too low, set to 0.001")
	}

	return p
}

func getCap(cap uint) uint {
	if cap == 0 {
		cap = 2 << 16
	}

	return cap
}
