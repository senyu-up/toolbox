package tests

import (
	"errors"
	"fmt"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/runtime"
	"github.com/senyu-up/toolbox/tool/trace"
	"github.com/senyu-up/toolbox/tool/wework/qwrobot"
	"testing"
	"time"
)

func SetUp() {

	var redConf = &config.RedisConfig{
		Addrs: []string{"127.0.0.1:6379"},
	}
	var redisCli = cache.InitRedisByConf(redConf)

	var hookUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=39eaf949-6bf4-4af8-9b67"

	var conf = &config.QwRobotConfig{
		Webhook:        hookUrl,
		InfoFreqLimit:  "1/s",
		WarnFreqLimit:  "1/s",
		ErrorFreqLimit: "1/s",
		MessageType:    "t_msg_type", // 消息类型
		Prefix:         "t_prefix",   // redis key 前缀
		RedisCli:       redisCli,
	}

	_, err := qwrobot.Init(conf, qwrobot.OptWithIp("127.0.0.1"),
		qwrobot.OptWithStage("local_test"),
		qwrobot.OptWithHostName("mac_arm64"))
	if err != nil {
		fmt.Printf("init qwrobot failed, err: %v", err)
	}
}

func TestNewQwRobot(t *testing.T) {
	SetUp()
	logger.SetCallBack(logger.NewQwRobot(qwrobot.Get()))

	var ext = logger.E()
	ext.Error(errors.New("Error Hooo!"))
	ext.String("set_str", "str_val")
	var ctx = trace.NewTrace()
	ext.Ctx(ctx)
	ext.String("app_name", "app_test")

	logger.Ctx(ctx).SetExtra(ext).Notify().Warn(" warn msg")

	ctx = trace.NewTrace()
	var ext2 = logger.E()
	ext2.String("meme", runtime.GetMemUsage())
	ext2.Span(trace.NewSpanID())

	logger.Ctx(ctx).SetErr(fmt.Errorf("This is err ")).Notify().Error("test err")

	logger.Info("just info log")
	time.Sleep(10 * time.Second)
}
