package runtime

import (
	"fmt"
	"runtime"
	"strings"
)

// GetString 获取堆string信息
func GetString() string {
	buf := make([]byte, 4096)
	bufLen := runtime.Stack(buf[:], false)
	return string(buf[:bufLen])
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func GetMemUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	var str strings.Builder
	str.WriteString(fmt.Sprintf("Alloc = %v MiB", bToMb(m.Alloc)))
	str.WriteString(fmt.Sprintf("\tHeapAlloc = %v MiB", bToMb(m.HeapAlloc)))
	str.WriteString(fmt.Sprintf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc)))
	str.WriteString(fmt.Sprintf("\tSys = %v MiB", bToMb(m.Sys)))
	str.WriteString(fmt.Sprintf("\tNumGC = %v\n", m.NumGC))
	str.WriteString(fmt.Sprintf("\tGoroutine = %v\n", runtime.NumGoroutine()))
	return str.String()
}
