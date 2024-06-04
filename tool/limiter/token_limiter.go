package limiter

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/senyu-up/toolbox/tool/logger"
	xrate "golang.org/x/time/rate"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	tokenFormat  = "Token:Limiter:Key:%s" // redis key，存储当前桶内令牌数
	pingInterval = time.Millisecond * 100 // 监测 redis 时间间隔
)

// @Deprecated 迁移到 traffic_limit 包
type TokenLimiter struct {
	// 每秒生产速率
	rate int
	// 桶容量
	burst int
	// 存储容器
	store redis.UniversalClient
	// lock
	rescueLock sync.Mutex
	// redis 健康标识
	redisAlive uint32
	// redis 故障时采用进程内 令牌桶限流器
	rescueLimiter *xrate.Limiter
	// redis 监控探测任务标识
	monitorStarted bool
	// lua 脚本
	script *redis.Script
}

// @Deprecated 建议使用 trafic_limit 包
func NewTokenLimiter(opts ...LimitOption) *TokenLimiter {
	var limiter = &TokenLimiter{redisAlive: 0, rate: 1, burst: 1}
	for _, opt := range opts {
		opt(limiter)
	}
	limiter.rescueLimiter = xrate.NewLimiter(xrate.Every(time.Second/time.Duration(limiter.rate)), limiter.burst)
	limiter.script = createScript()
	return limiter
}

// ReserveNow 根据当前时间获取一个令牌
func (lim *TokenLimiter) ReserveNow(key string) bool {
	return lim.ReserveN(time.Now(), key, 1)
}

// Reserve 根据时间获取 1 个令牌
func (lim *TokenLimiter) Reserve(time time.Time, key string) bool {
	return lim.ReserveN(time, key, 1)
}

// ReserveNowN 根据当前时间获取 n 个令牌
func (lim *TokenLimiter) ReserveNowN(key string, n int) bool {
	return lim.ReserveN(time.Now(), key, n)
}

// ReserveN 根据时间获取 n 个令牌
func (lim *TokenLimiter) ReserveN(time time.Time, key string, n int) bool {
	// 判断redis是否健康
	// redis故障时采用进程内限流器
	// 兜底保障
	if atomic.LoadUint32(&lim.redisAlive) == 0 {
		logger.Info("use local limiter")
		return lim.rescueLimiter.AllowN(time, n)
	}
	// 执行脚本获取令牌
	resp, err := lim.evalScript(
		lim.store,
		[]string{
			fmt.Sprintf(tokenFormat, key),
		},
		[]string{
			strconv.Itoa(lim.rate),
			strconv.Itoa(lim.burst),
			strconv.FormatInt(time.Unix(), 10),
			strconv.Itoa(n),
		})
	// redis allowed == false
	// Lua boolean false -> r Nil bulk reply
	// 特殊处理key不存在的情况
	if err == redis.Nil {
		return false
	} else if err != nil {
		logger.Error("fail to use rate limiter: %s, use in-process limiter for rescue", err)
		// 执行异常，开启redis健康探测任务
		// 同时采用进程内限流器作为兜底
		lim.startMonitor()
		return lim.rescueLimiter.AllowN(time, n)
	}

	code, ok := resp.(int64)
	if !ok {
		logger.Error("fail to eval redis script: %v, use in-process limiter for rescue", resp)
		lim.startMonitor()
		return lim.rescueLimiter.AllowN(time, n)
	}

	// redis allowed == true
	// Lua boolean true -> r integer reply with value of 1
	return code == 1
}

// 执行 lua 脚本
func (lim *TokenLimiter) evalScript(client redis.UniversalClient, keys []string, args ...interface{}) (interface{}, error) {
	sha, err := lim.script.Load(client).Result()
	if err != nil {
		logger.Error("script load err: %v", err.Error())
		return nil, err
	}
	ret := client.EvalSha(sha, keys, args...)
	result, err := ret.Result()
	if err != nil {
		if err != redis.Nil {
			logger.Error("Execute Redis fail: %v", err.Error())
		}
		return nil, err
	}
	return result, nil
}

// 开启redis健康探测
func (lim *TokenLimiter) startMonitor() {
	lim.rescueLock.Lock()
	defer lim.rescueLock.Unlock()
	// 防止重复开启
	if lim.monitorStarted {
		return
	}

	// 设置任务和健康标识
	lim.monitorStarted = true
	atomic.StoreUint32(&lim.redisAlive, 1)
	// 健康探测
	go lim.waitForRedis()
}

// redis健康探测定时任务
func (lim *TokenLimiter) waitForRedis() {
	ticker := time.NewTicker(pingInterval)
	// 健康探测成功时回调此函数
	defer func() {
		ticker.Stop()
		lim.rescueLock.Lock()
		lim.monitorStarted = false
		lim.rescueLock.Unlock()
	}()

	for range ticker.C {
		// ping属于redis内置健康探测命令
		pingRs := lim.store.Ping()
		if pingRs != nil && pingRs.Val() != "" {
			// 健康探测成功，设置健康标识
			atomic.StoreUint32(&lim.redisAlive, 1)
			return
		} else {
			atomic.StoreUint32(&lim.redisAlive, 0)
		}
	}
}

// 返回基于令牌桶实现的 lua 脚本
func createScript() *redis.Script {
	script := redis.NewScript(`
				-- 每秒生成token数量即token生成速度
				local rate = tonumber(ARGV[1])
				-- 桶容量
				local capacity = tonumber(ARGV[2])
				-- 当前时间戳
				local now = tonumber(ARGV[3])
				-- 当前请求token数量
				local requested = tonumber(ARGV[4])
				-- 需要多少秒才能填满桶
				local fill_time = capacity/rate
				-- 向下取整,ttl为填满时间的2倍
				local ttl = math.floor(fill_time*2)
				-- 获取当前桶容量及上次刷新时间   tokens:timestamp
				local value = redis.call("get", KEYS[1])
				if value == false then
				value = capacity..":".."0"
                end
				local index,_ = string.find(value, ":")
                -- 当前时间桶容量
                local last_tokens = tonumber(string.sub(value,1,index-1))
                -- 上一次刷新的时间
                local last_refreshed = tonumber(string.sub(value,index+1,-1))
                -- 距离上次请求的时间跨度
                local delta = math.max(0, now-last_refreshed)
                -- 距离上次请求的时间跨度,总共能生产token的数量,如果超多最大容量则丢弃多余的token
                local filled_tokens = math.min(capacity, last_tokens+(delta*rate))
                -- 本次请求token数量是否足够
                local allowed = filled_tokens >= requested
                -- 桶剩余数量
                local new_tokens = filled_tokens
                -- 允许本次token申请,计算剩余数量
                if allowed then
                new_tokens = filled_tokens - requested
                end
                -- 设置剩余token数量和刷新时间
                redis.call("setex", KEYS[1], ttl, new_tokens..":"..now)

				return allowed

	`)
	return script
}
