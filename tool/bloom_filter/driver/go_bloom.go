package driver

import (
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/senyu-up/toolbox/tool/config"
)

type GoBloomFilter struct {
	conf *config.GoBloomFilterConf
	inst *bloom.BloomFilter
}

func NewGoBloomFilter(conf *config.GoBloomFilterConf) (*GoBloomFilter, error) {
	b := &GoBloomFilter{
		conf: conf,
	}
	_ = b.init()
	b.inst = bloom.NewWithEstimates(conf.Cap, conf.P)

	return b, nil
}

func (g *GoBloomFilter) init() error {
	g.conf.Cap = getCap(g.conf.Cap)
	g.conf.P = getErrorRatio(g.conf.P)

	return nil
}

func (g *GoBloomFilter) Exists(key string) (exists bool, err error) {
	exists = g.inst.TestString(key)

	return
}

func (g *GoBloomFilter) MExists(keys []string) (statusList []bool, err error) {
	//TODO implement me
	statusList = make([]bool, len(keys))
	for i, _ := range keys {
		statusList[i] = g.inst.TestString(keys[i])
	}

	return
}

func (g *GoBloomFilter) Add(key string) (ok bool, err error) {
	ok = g.inst.TestOrAddString(key)

	return
}

func (g *GoBloomFilter) MAdd(keys []string) (statusList []bool, err error) {
	//TODO implement me
	statusList = make([]bool, len(keys))
	for i, _ := range keys {
		statusList[i] = g.inst.TestOrAddString(keys[i])
	}

	return
}

func (g *GoBloomFilter) Insert(key string) (ok bool, err error) {
	ok = g.inst.TestOrAddString(key)

	return
}

func (g *GoBloomFilter) MInsert(keys []string) (statusList []bool, err error) {
	//TODO implement me
	statusList = make([]bool, len(keys))
	for i, _ := range keys {
		statusList[i] = g.inst.TestOrAddString(keys[i])
	}

	return
}

func (g *GoBloomFilter) Info() (info map[string]interface{}, err error) {
	//TODO implement me
	info = map[string]interface{}{
		"hash_count": g.inst.K(),
		"cap":        g.inst.Cap(),
	}

	return
}
