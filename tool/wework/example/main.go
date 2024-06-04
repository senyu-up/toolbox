package main

import (
	"fmt"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/wework"
	"github.com/senyu-up/toolbox/tool/wework/qwrobot"
	"net/http"
)

func Example1() {
	var conf = &config.WeWorkConfig{
		CorpId:          "wwef87055fdc2e7c",
		Secret:          "LBamk7d3XPiAMqoULNxSYEyi0a",
		AgentId:         "1231",
		Debug:           true,
		RefreshInterval: 7000, // less 2 hours
	}

	wcClient, err := wework.InitByConfig(conf, wework.OptWithHttpClient(&http.Client{}))
	if err != nil {
		fmt.Printf("init wework client err: %v \n", err)
		return
	}
	// 获取部门列表
	if depList, err := wcClient.GetDepartmentList(0); err != nil {
		fmt.Printf("get department list err: %v \n", err)
	} else {
		fmt.Printf("get department list: %v \n", depList)
		// 遍历部门
		for _, dep := range depList {
			// 获取部门成员, 只获取当前部门一级成员，子部门的成员不获取
			if employeeList, err := wcClient.GetDepartmentSimple(dep.ID, 0); err != nil {
				fmt.Printf("get department user list err: %v \n", err)
			} else {
				fmt.Printf("get department user list: %v \n", employeeList)
				for _, employee := range employeeList {
					// 获取成员详情
					if employeeDetail, err := wcClient.GetUserInfoRequest(employee.Userid); err != nil {
						fmt.Printf("get employee detail err: %v \n", err)
					} else {
						fmt.Printf("get employee detail: %v \n", employeeDetail)
					}
				}
			}
		}
	}
}

func ExampleRobot() {
	// 初始化 redis, 用于机器人发送频率控制
	var redisClient = cache.InitRedisByConf(&config.RedisConfig{
		Addrs: []string{"127.0.0.1:6379"},
	})

	// 初始化 qwrobot
	var qwRobotConf = config.QwRobotConfig{
		Webhook:     "https://qyapi.weixin.com/xx",
		MessageType: "AQ",
		Prefix:      "xh_qw_robot_", // 企微机器人 redis key的前缀
		RedisCli:    redisClient,

		InfoFreqLimit:  "10/m", // 1分钟内最多发送10条
		WarnFreqLimit:  "10/m", // 1分钟内最多发送10条
		ErrorFreqLimit: "20/m", // 1分钟内最多发送20条, 把额度留给Error
	}
	// 初始化 qwrobot
	qwr, err := qwrobot.Init(&qwRobotConf,
		qwrobot.OptWithIp("127.0.0.1"),  // 设置当前运行环境的ip
		qwrobot.OptWithHostName("auth"), // 当前运行环境的主机名
		qwrobot.OptWithStage("local"),   // 当前运行环境的环境名
	)
	if err != nil {
		logger.SetErr(err).Error("Init QWRobot err ")
	}

	// 获取 qwrobot
	qwr = qwrobot.Get()

	// 组织消息
	var msg = qwrobot.Message{
		Title:    "Note",
		Content:  "Hi, Welcome!",
		UserList: []string{}, // @的成员列表,
	}
	// 发消息
	qwr.Info(msg)
	qwr.Warn(msg)
	qwr.Error(msg)
}

func main() {
	//ExampleRobot()
	Example1()
}
