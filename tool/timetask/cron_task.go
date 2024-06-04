package timetask

import (
	"fmt"
	"github.com/roylee0704/gron"
	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/cron"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/redis_lock"
	"runtime/debug"
)

// CronTaskFunc func待执行的定时任务 appKey：定时任务标识key  appName：定时任务名称
type CronTaskFunc func(appKey, appName string) error

func CronTaskRegister(schedule gron.AtSchedule, task CronTaskFunc) {
	key := GetFuncName(task)
	//定时任务一分钟一次
	cron.Register(schedule, func() {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				logger.Error("task panic:", panicErr, string(debug.Stack()))
			}
		}()
		CronTask(key, task, 0)
	})
}

// CronTask 通用定时处理任务方法
func CronTask(key string, task CronTaskFunc, expire int64) {
	if expire == 0 {
		expire = enum.ActionLockTime
	}
	//redis分布式锁（多台机器一个定时任务）
	lock, err := redis_lock.NewDistributeLockRedis(key, expire, "1")
	if err != nil {
		logger.Warn(key, " err:", err)
		return
	}
	defer lock.Unlock()
	logger.Info("start ", key)
	err = task("", key)
	if err != nil {
		logger.Warn(fmt.Sprintf("key: %s err:%v", key, err))
		return
	}

	logger.Info("start end ", key)
	return
}
