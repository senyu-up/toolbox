package su_logger

import (
	"context"
	"errors"
	"fmt"
	"github.com/senyu-up/toolbox/tool/trace"
	"runtime"
	"strings"
	"sync"
	"testing"
)

// Try each log level in decreasing order of priority.
func testZapCallTrace() {
	var err = errors.New("This is err ")
	var ctx = trace.NewContextWithRequestIdAndSpanId(context.TODO(), "ctx_req1", "ctx_span1")

	Error(ctx, err, "error ctx trace %v, err is %v", E().Ctx(ctx).Error(err).Trace("req11"))
	Warn(ctx, "warning with trace")
	Debug(ctx, "notice")
	Info(ctx, "informational1")
}

func testZapCallTrace2() {
	var err = errors.New("new err 2")
	var ctx = trace.NewContextWithRequestIdAndSpanId(context.TODO(), "req2", "span2")

	Error(ctx, err, "error 2ctx trace %v, err2 is %v", E().Ctx(ctx).Error(err).Trace("req22"))
	Warn(ctx, "warning2 with trace")
	Debug(ctx, "notice2")
	Info(ctx, "informational2")
}

func TestZapRace(t *testing.T) {
	var count = 100000

	var wg = sync.WaitGroup{}
	wg.Add(1)
	go func() {
		t.Logf("run testZapCallTrace")
		for i := 0; i < count; i++ {
			testZapCallTrace()
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		t.Logf("run testZapCallTrace")
		for i := 0; i < count; i++ {
			testZapCallTrace2()
		}
		wg.Done()
	}()

	wg.Wait()

	t.Logf("done %v", GetMemUsage())
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
