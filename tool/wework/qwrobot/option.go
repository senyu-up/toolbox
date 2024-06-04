package qwrobot

type QwrOption func(*QWRobot)

func OptWithIp(ip string) QwrOption {
	return func(option *QWRobot) {
		option.ip = ip
	}
}

func OptWithHostName(hn string) QwrOption {
	return func(option *QWRobot) {
		option.hostName = hn
	}
}

func OptWithStage(stage string) QwrOption {
	return func(option *QWRobot) {
		option.stage = stage
	}
}

/*func OptWithRedis(redis redis.UniversalClient) QwrOption {
	return func(option *QWRobot) {
		option.redisCli = redis
	}
}*/
