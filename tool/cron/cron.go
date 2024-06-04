package cron

import (
	"github.com/roylee0704/gron"
)

var client = gron.New()

func Register(schedule gron.AtSchedule, cmd func()) {
	client.AddFunc(schedule, cmd)
}

func Start() {
	client.Start()
}

func Stop() {
	client.Stop()
}
