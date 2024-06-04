package event

import (
	"github.com/panjf2000/ants/v2"
	"github.com/senyu-up/toolbox/tool/logger"
	"runtime/debug"
)

var pool *ants.Pool

func init() {
	pool, _ = ants.NewPool(10000, ants.WithPanicHandler(func(i interface{}) {
		logger.Error(string(debug.Stack()))
	}))
}
