package env

type AppOption func(*AppInfo)

func OptWithIp(ip string) AppOption {
	return func(option *AppInfo) {
		option.Ip = ip
	}
}

func OptWithHostName(hn string) AppOption {
	return func(option *AppInfo) {
		option.HostName = hn
	}
}

func OptWithStage(stage string) AppOption {
	return func(option *AppInfo) {
		option.Stage = stage
	}
}

func OptWithAppName(name string) AppOption {
	return func(option *AppInfo) {
		option.Name = name
	}
}
