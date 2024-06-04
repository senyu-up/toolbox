package qwrobot

import (
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/logger"
	"testing"
)

func TestInit(t *testing.T) {
	// 初始化 redis
	var redisClient = cache.InitRedisByConf(&config.RedisConfig{
		Addrs: []string{"127.0.0.1:6379"},
	})

	// 初始化 qwrobot
	var qwRobotConf = config.QwRobotConfig{
		Webhook:     "https://qyapi.weixin.com/xx",
		MessageType: "AQ",
		Prefix:      "【】",
		RedisCli:    redisClient,
	}
	qwr, err := Init(&qwRobotConf,
		OptWithIp("127.0.0.1"),
		OptWithHostName("auth"),
		OptWithStage("local"),
	)
	if err != nil {
		logger.SetErr(err).Error("Init QWRobot err ")
	}

	// 获取 qwrobot
	qwr = Get()

	// 组织消息
	var msg = Message{
		Title:    "Note",
		Content:  "Hi, Welcome!",
		UserList: []string{}, // 成员列表
	}
	// 发消息
	qwr.Info(msg)
	qwr.Warn(msg)
	qwr.Error(msg)
}
