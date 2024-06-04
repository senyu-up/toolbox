package driver

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/senyu-up/toolbox/tool/config"
	"strings"
)

// RedisBloomFilter
// @Description: 布隆过滤器
type RedisBloomFilter struct {
	conf *config.RedisBloomFilterConf
}

func NewRedisBloom(conf *config.RedisBloomFilterConf) (*RedisBloomFilter, error) {
	b := &RedisBloomFilter{
		conf: conf,
	}

	err := b.init()
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			return b, ErrItemAlreadyExists
		}
	}

	return b, err
}

func (b *RedisBloomFilter) pipe(pipeArgs ...[]interface{}) ([]redis.Cmder, error) {
	if len(pipeArgs) == 1 {
		pipe := b.conf.Cache.Pipeline()
		defer pipe.Close()
		for i, _ := range pipeArgs {
			pipe.Do(pipeArgs[i]...)
		}
		return pipe.Exec()
	} else {
		tx := b.conf.Cache.TxPipeline()
		defer tx.Close()
		for i, _ := range pipeArgs {
			tx.Do(pipeArgs[i]...)
		}
		return tx.Exec()
	}
}

func (b *RedisBloomFilter) init() (err error) {
	b.conf.Cap = getCap(b.conf.Cap)
	b.conf.P = getErrorRatio(b.conf.P)

	cmdList := make([][]interface{}, 0, 2)
	if b.conf.Cap < 1 {
		cmdList = append(cmdList, []interface{}{"BF.RESERVE", b.conf.Key, b.conf.P})
	} else {
		cmdList = append(cmdList, []interface{}{"BF.RESERVE", b.conf.Key, b.conf.P, b.conf.Cap})
	}
	if b.conf.Ttl > 0 {
		ttl := b.conf.Ttl.Seconds()
		cmdList = append(cmdList, []interface{}{"EXPIRE", b.conf.Key, ttl})
	}

	_, err = b.pipe(cmdList...)

	return err
}

// Exists
// @description 判断当前key是否存在
func (b *RedisBloomFilter) Exists(key string) (exists bool, err error) {
	cmderList, err := b.pipe([]interface{}{"BF.EXISTS", b.conf.Key, key})
	if err != nil {
		return false, err
	}

	exists, err = cmderList[0].(*redis.Cmd).Bool()

	return
}

// MExists
// @description 批量查询key是否存在
func (b *RedisBloomFilter) MExists(keys []string) (statusList []bool, err error) {
	cmds := make([]interface{}, 0, len(keys)+2)
	cmds = append(cmds, "BF.MEXISTS", b.conf.Key)
	for i, _ := range keys {
		cmds = append(cmds, keys[i])
	}
	cmderList, err := b.pipe(cmds)
	if err != nil {
		return nil, err
	}
	rs, err := cmderList[0].(*redis.Cmd).Result()
	if err != nil {
		return nil, err
	}
	rsList, ok := rs.([]interface{})
	if !ok {
		return nil, fmt.Errorf("BF.MADD result is not list")
	}
	statusList = make([]bool, 0, len(rsList))
	for i, _ := range rsList {
		statusList = append(statusList, rsList[i].(int64) == 1)
	}

	return
}

// Insert
// @description 插入, 如果缓存的key不存在会自动创建
func (b *RedisBloomFilter) Insert(key string) (ok bool, err error) {
	rsList, err := b.MInsert([]string{key})
	if err != nil {
		return
	}

	return rsList[0], nil
}

// MInsert
// @description 批量插入, 如果缓存的key不存在会自动创建
func (b *RedisBloomFilter) MInsert(keys []string) (statusList []bool, err error) {
	var cmds []interface{}
	if b.conf.Cap < 1 {
		cmds = make([]interface{}, 0, len(keys)+5)
		cmds = append(cmds, "BF.INSERT", b.conf.Key, "error", b.conf.P, "ITEMS")
	} else {
		cmds = make([]interface{}, 0, len(keys)+7)
		cmds = append(cmds, "BF.INSERT", b.conf.Key, "CAPACITY", b.conf.Cap, "error", b.conf.P, "ITEMS")
	}
	for i, _ := range keys {
		cmds = append(cmds, keys[i])
	}

	var cmderList []redis.Cmder
	if b.conf.Ttl > 0 {
		ttl := b.conf.Ttl.Seconds()
		ttlCmds := append(cmds, "EXPIRE", b.conf.Key, ttl)
		cmderList, err = b.pipe(cmds, ttlCmds)
	} else {
		cmderList, err = b.pipe(cmds)
	}
	if err != nil {
		return nil, err
	}

	rs, err := cmderList[0].(*redis.Cmd).Result()
	if err != nil {
		return nil, err
	}
	rsList, ok := rs.([]interface{})
	if !ok {
		return nil, fmt.Errorf("BF.MADD result is not list")
	}
	statusList = make([]bool, 0, len(rsList))
	for i, _ := range rsList {
		statusList = append(statusList, rsList[i].(int64) == 1)
	}

	return
}

// Add
// @description 添加元素, 如果存在则返回false, 反之true
func (b *RedisBloomFilter) Add(key string) (ok bool, err error) {
	cmderList, err := b.pipe([]interface{}{"BF.ADD", b.conf.Key, key})
	if err != nil {
		return false, err
	}

	ok, err = cmderList[0].(*redis.Cmd).Bool()

	return
}

// MAdd
// @description 批量添加
func (b *RedisBloomFilter) MAdd(keys []string) (statusList []bool, err error) {
	cmds := make([]interface{}, 0, len(keys)+2)
	cmds = append(cmds, "BF.MADD", b.conf.Key)
	for i, _ := range keys {
		cmds = append(cmds, keys[i])
	}
	cmderList, err := b.pipe(cmds)
	if err != nil {
		return nil, err
	}

	rs, err := cmderList[0].(*redis.Cmd).Result()
	if err != nil {
		return nil, err
	}
	rsList, ok := rs.([]interface{})
	if !ok {
		return nil, fmt.Errorf("BF.MADD result is not list")
	}
	statusList = make([]bool, 0, len(rsList))
	for i, _ := range rsList {
		statusList = append(statusList, rsList[i].(int64) == 1)
	}

	return
}

// Info
// @description 获取BloomFilter信息
func (b *RedisBloomFilter) Info() (info map[string]interface{}, err error) {
	cmderList, err := b.pipe([]interface{}{"BF.INFO", b.conf.Key})
	if err != nil {
		return nil, err
	}

	rs, err := cmderList[0].(*redis.Cmd).Result()
	if err != nil {
		return nil, err
	}
	rsList, ok := rs.([]interface{})
	if !ok {
		return nil, fmt.Errorf("BF.INFO result is not list")
	}
	l := len(rsList)
	info = make(map[string]interface{}, l>>1)
	for i := 0; i < l; i += 2 {
		info[rsList[i].(string)] = rsList[i+1]
	}

	return
}
