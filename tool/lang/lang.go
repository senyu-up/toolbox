package lang

import "strings"

// 转换zh开头的语言为zh
func ToZhStr(str string) string {
	if strings.HasPrefix(str, "zh") {
		return "zh"
	}
	return str
}
