package kafka

import (
	"github.com/IBM/sarama"
	"strconv"
)

// GetStringFromHeader
// @description 从header中获取string类型数据
func GetStringFromHeader(header []*sarama.RecordHeader, key string) string {
	if header == nil {
		return ""
	}

	for _, h := range header {
		if string(h.Key) == key {
			return string(h.Value)
		}
	}
	return ""
}

// GetIntFromHeader
// @description 从header中获取int类型数据
func GetIntFromHeader(header []*sarama.RecordHeader, key string) int {
	if header == nil {
		return 0
	}

	for _, h := range header {
		if string(h.Key) == key {
			times, _ := strconv.Atoi(string(h.Value))
			return times
		}
	}
	return 0
}

// AddIntToHeader
// @description 向header添加item, 如果item已存在会替换成当前值
func AddIntToHeader(header *[]*sarama.RecordHeader, key string, val int) {
	if header == nil {
		*header = append(*header, &sarama.RecordHeader{Key: []byte(key), Value: []byte(strconv.Itoa(val))})
	} else {
		var find bool
		for i, recordHeader := range *header {
			if string(recordHeader.Key) == key {
				find = true
				(*header)[i].Value = []byte(strconv.Itoa(val))
				return
			}
		}

		if find == false {
			*header = append(*header, &sarama.RecordHeader{Key: []byte(key), Value: []byte(strconv.Itoa(val))})
		}
	}
}

// AddStringToHeader
// @description 向header添加item, 如果item已存在会替换成当前值
func AddStringToHeader(header *[]*sarama.RecordHeader, key string, val string) {
	if header == nil {
		*header = append(*header, &sarama.RecordHeader{Key: []byte(key), Value: []byte(val)})
	} else {
		var find bool
		for i, recordHeader := range *header {
			if string(recordHeader.Key) == key {
				find = true
				(*header)[i].Value = []byte(val)
				return
			}
		}

		if find == false {
			*header = append(*header, &sarama.RecordHeader{Key: []byte(key), Value: []byte(val)})
		}
	}
}
