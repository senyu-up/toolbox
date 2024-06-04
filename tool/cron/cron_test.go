package cron

import (
	"fmt"
	"github.com/roylee0704/gron"
	"testing"
	"time"
)

func Task() {
	fmt.Printf(" task run \n")
	time.Sleep(time.Second)
}

func Task2() {
	fmt.Printf(" task 2 run \n")
	time.Sleep(time.Second)
}

func Example() {
	// 注册 cron 函数
	Register(gron.Every(1*time.Second), Task)
	Register(gron.Every(2*time.Second), Task2)

	defer Stop()
	Start() // 非阻塞运行

	select {}
}

func TestCronTask(t *testing.T) {
	Example()
}
