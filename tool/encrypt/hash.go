package encrypt

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"strings"
)

// MD5
// @description 大写形式的MD5
func MD5(str string) string {
	s := Md5(str)
	return strings.ToUpper(s)
}

// Md5
// @description 小写形式的md5
func Md5(str string) string {
	data := []byte(str) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

func CRC32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}
