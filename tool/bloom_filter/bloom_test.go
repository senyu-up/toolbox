package bloom_filter

import (
	"fmt"
	"github.com/senyu-up/toolbox/tool/config"
	"testing"
	"time"
)

func TestMAdd(t *testing.T) {
	bf, _ := NewRedisBloom(&config.RedisBloomFilterConf{
		P:   0.001,
		Cap: 1000,
		Key: "test",
		//Cache: boot.RedisInst.GetInst().(redis.UniversalClient),
		Ttl: time.Minute,
	})
	bf.MAdd([]string{"test", "test1", "test2"})
	rs, err := bf.MExists([]string{"test", "test1", "aaa", "bbb"})
	if err != nil {
		t.Error(err)
		return
	}
	//4246724928
	//2147483648
	//
	//info, err := bf.Info()

	fmt.Println(rs, err)
}

func TestGoBloom(t *testing.T) {
	bf, _ := NewGoBloom(&config.GoBloomFilterConf{
		P:   0.001,
		Cap: 1000,
	})

	add, err2 := bf.MAdd([]string{"aaa", "bbb", "ccc"})
	fmt.Println(add, err2)

	exists, err := bf.MExists([]string{"aaa", "bbb", "ccc", "ddd", "eee"})
	fmt.Println(exists, err)
}
