package config

import "time"

type BreakerConf struct {
	// 单次操作的超时时间, 超过设定设定时间直接返回错误
	Timeout time.Duration
	// 相同命令(name)再同一时间的最大并发数
	MaxConcurrentRequests int
	// 触发熔断器检测的最小处理数量
	RequestVolumeThreshold int
	// 错误率阈值, 达到阈值, 启动熔断, n%
	ErrorPercentThreshold int
	// 当熔断开始后, 尝试恢复正常的探测周期
	// 放开对接口的限制（断路器状态 HALF-OPEN），然后尝试使用 1 个请求去调用接口，如果调用成功，则恢复正常（断路器状态 CLOSED），如果调用失败或出现超时等待，就需要再重新等待circuitBreakerSleepWindowInMilliseconds 的时间
	SleepWindow time.Duration
}
