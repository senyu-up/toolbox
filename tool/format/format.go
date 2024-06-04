package format

import "fmt"

// Byte2Human
// @description 字节转人类友好的格式
func Byte2Human(byteSize int64) (size string) {
	if byteSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(byteSize)/float64(1))
	} else if byteSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(byteSize)/float64(1024))
	} else if byteSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(byteSize)/float64(1024*1024))
	} else if byteSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(byteSize)/float64(1024*1024*1024))
	} else if byteSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(byteSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(byteSize)/float64(1024*1024*1024*1024*1024))
	}
}
