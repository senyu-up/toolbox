package http_health

import "github.com/senyu-up/toolbox/tool/config"

type HealthOption func(box *config.HealthCheck)

func HealthOptionWithPprof(pprof bool) HealthOption {
	return func(option *config.HealthCheck) {
		if pprof {
			option.Pprof = pprof
		}
	}
}

func HealthOptionWithAddr(addr string) HealthOption {
	return func(option *config.HealthCheck) {
		if addr != "" {
			option.Addr = addr
		}
	}
}

func HealthOptionWithPort(port uint32) HealthOption {
	return func(option *config.HealthCheck) {
		if port != 0 {
			option.Port = port
		}
	}
}

func HealthOptionWithDisableLog(disable bool) HealthOption {
	return func(option *config.HealthCheck) {
		option.DisableLog = disable
	}
}
