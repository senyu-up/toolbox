package driver

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/senyu-up/toolbox/tool/config"
	"math"
)

var luaExistsCheckScript = redis.NewScript(`
local exists = redis.call("EXISTS", KEYS[1])
if exists == 0 then
	return 0 
end

for k, v in ipairs(ARGV) do 
	if redis.call("GETBIT", KEYS[1], v) == 0 then
		exists = 0
		break
	end
end
return exists
`)

var luaAddScript = redis.NewScript(`
local keyExists = redis.call("EXISTS", KEYS[1])
local exists = 1
for k, v in ipairs(ARGV) do
	if k > 1 and redis.call("GETBIT", KEYS[1], v) == 0 then
		exists = 0
		break
	end
end
if exists == 1 then
	return 0
end
for k, v in ipairs(ARGV) do 
	redis.call("SETBIT", KEYS[1], v, 1)
end
local ttl = tonumber(ARGV[1])
if keyExists == 0 and ttl > 0 then
	redis.call("EXPIRE", KEYS[1], ttl)
end
return 1
`)

type BitsetBloom struct {
	key   string
	cap   uint
	p     float64
	ttl   int
	redis redis.UniversalClient
	m     uint
	k     uint
}

func (b *BitsetBloom) init() error {
	b.cap = getCap(b.cap)
	b.p = getErrorRatio(b.p)

	b.m = uint(math.Ceil(-1 * float64(b.cap) * math.Log(b.p) / math.Pow(math.Log(2), 2)))
	b.k = uint(math.Ceil(math.Log(2) * float64(b.m) / float64(b.cap)))

	if b.redis == nil {
		return errors.New("redis instance is nil")
	}

	return nil
}

func NewBitsetBloom(conf *config.BitsetBloomConf) (*BitsetBloom, error) {
	b := &BitsetBloom{
		key:   conf.Key,
		cap:   conf.Cap,
		redis: conf.Cache,
		ttl:   int(conf.Ttl.Seconds()),
		p:     conf.P,
	}

	err := b.init()
	return b, err
}

// location returns the ith hashed location using the four base hash values
func location(h [4]uint64, i uint) uint64 {
	ii := uint64(i)
	return h[ii%2] + ii*h[2+(((ii+(ii%2))%4)/2)]
}

// location returns the ith hashed location using the four base hash values
func (b *BitsetBloom) location(h [4]uint64, i uint) uint64 {
	return location(h, i) % uint64(b.m)
}

func (b *BitsetBloom) getLocation(key []byte) []interface{} {
	var d digest128 // murmur hashing
	hash1, hash2, hash3, hash4 := d.sum256(key)
	hashList := [4]uint64{hash1, hash2, hash3, hash4}
	rs := make([]interface{}, b.k, b.k)
	var i uint = 0
	for ; i < b.k; i++ {
		rs[i] = b.location(hashList, i)
	}

	return rs
}

func (b *BitsetBloom) Exists(key string) (exists bool, err error) {
	locations := b.getLocation([]byte(key))
	rs := luaExistsCheckScript.Eval(b.redis, []string{b.key}, locations...)
	if rs.Err() != nil {
		return false, rs.Err()
	}

	v, err := rs.Int()
	if err != nil {
		return false, err
	}

	return v == 1, nil
}

func (b *BitsetBloom) MExists(keys []string) (statusList []bool, err error) {
	//TODO implement me
	for i, _ := range keys {
		exists, err := b.Exists(keys[i])
		if err != nil {
			return nil, err
		}
		statusList = append(statusList, exists)
	}

	return statusList, nil
}

func (b *BitsetBloom) Add(key string) (ok bool, err error) {
	//TODO implement me
	locations := make([]interface{}, 0, b.k+1)
	locations = append(locations, b.ttl)
	locations = append(locations, b.getLocation([]byte(key))...)
	rs := luaAddScript.Eval(b.redis, []string{b.key}, locations...)
	if rs.Err() != nil {
		return false, rs.Err()
	}

	v, err := rs.Int()
	if err != nil {
		return false, err
	}

	return v == 1, nil
}

func (b *BitsetBloom) MAdd(keys []string) (statusList []bool, err error) {
	//TODO implement me
	for i, _ := range keys {
		ok, err := b.Add(keys[i])
		if err != nil {
			return nil, err
		}
		statusList = append(statusList, ok)
	}

	return statusList, nil
}

func (b *BitsetBloom) Info() (info map[string]interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *BitsetBloom) Insert(key string) (ok bool, err error) {
	//TODO implement me
	return b.Add(key)
}

func (b *BitsetBloom) MInsert(keys []string) (statusList []bool, err error) {
	//TODO implement me
	return b.MExists(keys)
}
