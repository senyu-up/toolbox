package timetask

import (
	"fmt"
	"github.com/roylee0704/gron"
	"github.com/senyu-up/toolbox/tool/cron"
	"testing"
	"time"
)

func TestCronTaskRegister(t *testing.T) {
	CronTaskRegister(gron.Every(1*time.Second), CronFunc)
	cron.Start()
	select {}
}

func CronFunc(appKey, appName string) error {
	fmt.Println("cron task run", appKey, appName)
	return nil
}
