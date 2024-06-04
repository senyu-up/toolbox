package redis_lock

import (
	"github.com/senyu-up/toolbox/tool/logger"
	"testing"
	"time"
)

func TestNewDistributeLockRedis(t *testing.T) {
	//获得锁实例
	lock, err := NewDistributeLockRedis("test", 20, "1")
	if err != nil {
		//上锁失败
		logger.Warn(err.Error())
		return
	}
	task := func() {
		//测试锁竞争
		for i := 0; i < 20; i++ {
			time.Sleep(time.Second * 5)
			lock2, err := NewDistributeLockRedis("test", 20, "1")
			if err != nil {
				//上锁失败
				logger.Warn(err.Error())
				continue
			}
			err = lock2.Unlock()
			logger.Warn(err.Error())
		}
	}
	go task()
	//过期时间20s,设置任务执行时间为50s,测试看门狗是否生效
	time.Sleep(time.Second * 100)
	//释放锁
	err = lock.Unlock()
	logger.Warn(err.Error())
	return
}
