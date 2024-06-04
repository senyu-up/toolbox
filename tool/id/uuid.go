package id

import (
	uuid "github.com/satori/go.uuid"
	"github.com/senyu-up/toolbox/tool/snowflake"
	"github.com/yitter/idgenerator-go/idgen"
	"sync"
	"time"
)

// GetWithUUID
// @description
// Version 1,基于 timestamp 和 MAC address (RFC 4122)
// Version 2,基于 timestamp, MAC address 和 POSIX UID/GID (DCE 1.1)
// Version 3, 基于 MD5 hashing (RFC 4122)
// Version 4, 基于 random numbers (RFC 4122)
// Version 5, 基于 SHA-1 hashing (RFC 4122)
func GetWithUUID() string {
	return uuid.NewV4().String()
}

var nodeMap = make(map[int64]*snowflake.Node, 16)
var nodeMapLock = sync.Mutex{}

func getSnowflakeNode(node int64) (n *snowflake.Node, err error) {
	if _, ok := nodeMap[node]; !ok {
		nodeMapLock.Lock()
		defer nodeMapLock.Unlock()
		// 双重锁判断
		if _, ok = nodeMap[node]; !ok {
			nodeMap[node], err = snowflake.NewNode(node)
		}
	}

	return nodeMap[node], err
}

// GetStringWithSnowflake
// @description获取string类型的雪花id
func GetStringWithSnowflake(node int64) (string, error) {
	newNode, err := getSnowflakeNode(node)
	if err != nil {
		return "", err
	}
	s := newNode.Generate().String()

	return s, nil
}

// GetIntWithSnowflake
// @description 获取int64类型的雪花id
func GetIntWithSnowflake(node int64) (int64, error) {
	newNode, err := getSnowflakeNode(node)
	if err != nil {
		return 0, err
	}
	s := newNode.Generate().Int64()

	return s, nil
}

// InitIdGen
// options.WorkerIdBitLength = 6  // 默认值6，限定 WorkerId 最大值为2^6-1，即默认最多支持64个节点。
// options.SeqBitLength = 6; // 默认值6，限制每毫秒生成的ID个数。若生成速度超过5万个/秒，建议加大 SeqBitLength 到 10。
// options.BaseTime = Your_Base_Time // 如果要兼容老系统的雪花算法，此处应设置为老系统的BaseTime。
type GenParam struct {
	BitLength      int
	BaseTimeMs     int64
	WorkerIdBitLen int
}

func InitIdGen(workId uint16, param ...GenParam) {
	// NewIdGeneratorOptions
	opt := idgen.NewIdGeneratorOptions(workId)
	opt.WorkerIdBitLength = 6
	opt.SeqBitLength = 6
	setBaseTimeFlag := false
	if len(param) > 0 {
		if param[0].BitLength > 0 {
			opt.SeqBitLength = byte(param[0].BitLength)
		}
		if param[0].WorkerIdBitLen > 0 {
			opt.SeqBitLength = byte(param[0].WorkerIdBitLen)
		}

		if param[0].BaseTimeMs > 0 {
			opt.BaseTime = param[0].BaseTimeMs
			setBaseTimeFlag = true
		}
	}

	if !setBaseTimeFlag {
		opt.BaseTime = time.Date(2023, 3, 8, 12, 13, 14, 15, time.UTC).UnixNano() / 1e6
	}

	idgen.SetIdGenerator(opt)
}

// GetIdByIdGen 之前需要调用 InitIdGen() 进行初始化
func GetIdByIdGen() int64 {
	return idgen.NextId()
}
