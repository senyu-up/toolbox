package timetask

import (
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// GetFuncName 获取函数名字符串
func GetFuncName(i any) string {
	nameFull := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	nameEnd := filepath.Ext(nameFull)        // .foo
	name := strings.TrimPrefix(nameEnd, ".") // foo
	return name
}
