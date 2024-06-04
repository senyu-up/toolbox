package logger

import (
	"context"
	"errors"
	"fmt"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/trace"
	"runtime"
	"strings"
	"sync"
	"testing"
)

// Try each log level in decreasing order of priority.
func testZapCallTrace(bl Log) {
	var err = errors.New("This is err ")
	var ctx = trace.NewContextWithRequestIdAndSpanId(context.TODO(), "ctx_req1", "ctx_span1")

	bl.SetErr(err).Ctx(ctx).Error("error ctx trace %v, err is %v", trace.GetRequestId(ctx), err)
	bl.SpanId("set_spanId").TraceId("set_trace_id").Warn("warning with trace")
	bl.Debug("notice")
	bl.Info("informational1")
	bl.Trace("trace1")
}

// Try each log level in decreasing order of priority.
func testStaticZapTrace() {
	var err = errors.New("This is global err ")
	var ctx = trace.NewContextWithRequestIdAndSpanId(context.TODO(), "ctx_req2", "ctx_span2")

	Ctx(ctx).SetErr(err).Error("global error ctx trace %v, err is %v", trace.GetRequestId(ctx), err)
	SpanId("set_global_spanId").TraceId("set_global_trace_id").Info("global info with trace")
	Trace("trace2")
}

func testStaticZapTrace2() {
	var err = errors.New("This is global 3err ")
	var ctx = trace.NewContextWithRequestIdAndSpanId(context.TODO(), "ctx_req3", "ctx_span3")

	Ctx(ctx).SetErr(err).Error("error global3 trace %v, err is %v", trace.GetRequestId(ctx), err)
	SpanId("set_req3").TraceId("span3").Info("info 3 with trace")
	Trace("trace3")
}

func TestZapRace(t *testing.T) {
	var count = 100000
	var conf = &config.LogConfig{
		DefaultLog: AdapterZap,
		AppName:    "test_file",
		CallDepth:  4,

		Zap: &config.ZapConfig{
			Level:      "TRAC",
			Colorful:   true,
			CallerSkip: 4,
		},
	}
	InitDefaultLoggerByConf(conf)
	log1 := GetLogger()
	var wg = sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Logf("run testZapCallTrace")
		for i := 0; i < count; i++ {
			testZapCallTrace(log1)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Logf("run testZapCallTrace")
		for i := 0; i < count; i++ {
			testZapCallTrace(log1)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Logf("run testZapCallTrace")
		for i := 0; i < count; i++ {
			testStaticZapTrace()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Logf("run testZapCallTrace2")
		for i := 0; i < count; i++ {
			testStaticZapTrace2()
		}
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
