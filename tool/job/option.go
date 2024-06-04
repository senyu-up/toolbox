package job

import "github.com/senyu-up/toolbox/tool/logger"

type JobOption func(*AsyncJob)

func OptWithMaxJob(max int64) JobOption {
	return func(option *AsyncJob) {
		option.maxJobs = max
	}
}

func OptWithLogger(logger logger.Log) JobOption {
	return func(option *AsyncJob) {
		option.log = logger
	}
}
