# 功能点

- 支持markdown,text 两种企微卡片消息
- 支持消息限流, 限流配置 n/S (n次每秒), n/M (n次每分), n/H (n次每小时)
- 支持 Info, Error, Warn 三种消息级别

# 使用

1. 引入包

```shell
go get -u "toolbox/utils/qwrobot"
```

2. 配置文件 `config.yaml` 新增配置

```yaml
qwRobot:
  webhook: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=25e1862e-d324-4a27-b3a8-beb71dae1597
  infoFreqLimit: 10/H
  warnFreqLimit: 10/H
  errorFreqLimit: 10/H
  messageType: data-operation
  prefix: do
```

4. 基于log方式触发

需要boot.go中对qwrobot进行初始化

```go
qwrobotCnf := config.GetConf().QwRobot
qwrobot.InitDefaultQWRobot(qwrobot.Config{
    Webhook:        qwrobotCnf.Webhook,
    InfoFreqLimit:  qwrobotCnf.InfoFreqLimit,
    WarnFreqLimit:  qwrobotCnf.WarnFreqLimit,
    ErrorFreqLimit: qwrobotCnf.ErrorFreqLimit,
    MessageType:    qwrobotCnf.MessageType,
    Prefix:         qwrobotCnf.Prefix,
    RedisCli:       MysqlInst.Redis,
})
```

7. 常规方式调用

```go
rob := qwrobot.New(qwrobot.Config{
    Webhook:        "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=25e1862e-d324-4a27-b3a8-beb71dae1597",
    InfoFreqLimit:  "1/M",
    WarnFreqLimit:  "2/M",
    ErrorFreqLimit: "10/H",
    MessageType:    "test",
    Prefix:         "qw",
	RedisCli: NewRedisClient(),
})

// 错误
rob.Error(qwrobot.Message{
    Title:    "error title",
    Content:  "error content",
    UserList: nil,
})

// 常规消息
rob.Info(qwrobot.Message{
    Title:    "info title",
    Content:  "info content",
    UserList: nil,
})

// 警告
rob.Warn(qwrobot.Message{
    Title:    "warning title",
    Content:  "warning content",
    UserList: nil,
})
```

# 注意事项

- 当UserList不为空时, 会将消息类型转成文本类型, 因为 markdown 类型不支持 @某某某
- 消息体企微限制4096, 超过部分会被截取
- 频率限制仅支持 n/S (n次每秒), n/M (n次每分), n/H (n次每小时), 匹配失败会使用默认值 100/H