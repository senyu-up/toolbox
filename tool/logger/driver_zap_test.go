package logger

import (
	"context"
	"fmt"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/trace"
	"testing"
)

// Try each log level in decreasing order of priority.
func testZapCalls(bl Log) {
	//bl.Emer("emergency")
	//bl.Alert("alert")
	bl.Crit("critical")
	bl.Error("error")
	bl.Warn("warning")
	bl.Debug("notice")
	bl.Info("informational")
	bl.Trace("trace")
}

func testZapExtra(bl Log) {
	var ctx = trace.NewTrace()
	var err = fmt.Errorf("test error")
	bl.Ctx(ctx).SetErr(err).Error("error %v", err)
	bl.Ctx(ctx).Warn("Warn", err)

	ctxc, cancel := context.WithCancel(context.TODO())
	cancel()
	bl.Ctx(ctxc).Info("print canceled context")

	bl.Ctx(nil).Info("print nil context")
}

// Try each log level in decreasing order of priority.
func testStaticZapCalls() {
	//Emer("emergency")
	//Alert("alert")
	Crit("critical")
	Error("error")
	Warn("warning")
	Debug("notice")
	Info("informational")
	Trace("trace")
}

func TestZapStd(t *testing.T) {
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
	testZapCalls(log1)
	fmt.Print("\n\n")
	testZapExtra(log1)
	fmt.Print("\n\n")
	testStaticZapCalls()
	fmt.Print("\n\n")

	conf.Zap.Level = "EROR"
	InitDefaultLoggerByConf(conf)
	log2 := GetLogger()
	testZapCalls(log2)
	fmt.Print("\n\n")
	testStaticZapCalls()

}

func TestZap_Init(t *testing.T) {
	var conf = &config.LogConfig{
		DefaultLog: AdapterFile,
		AppName:    "test_file",
		CallDepth:  4,

		Zap: &config.ZapConfig{
			Output:     "test_zap.log",
			Level:      "TRAC",
			Colorful:   true,
			CallerSkip: 4,
		},
	}

	InitDefaultLoggerByConf(conf)
	log1 := GetLogger()
	testZapCalls(log1)
	fmt.Print("\n\n")
	testStaticZapCalls()
	fmt.Print("\n\n")

	conf.Zap.Level = "EROR"
	InitDefaultLoggerByConf(conf)
	log2 := GetLogger()
	testZapCalls(log2)
	fmt.Print("\n\n")
	testStaticZapCalls()

}

func TestZapNilExtra(t *testing.T) {
	SwitchLogger(AdapterZap)
	var nilStr string
	SetExtra(E().String("nil", nilStr)).Info("print nil extra")

	type nilStruct struct {
		Name string
		tie  *nilStruct
	}
	type GetStreamPrefixParam struct {
		// appId appKey stream_prefix 三选一
		AppId        int32  `json:"app_id"`
		AppKey       string `json:"app_key"`
		StreamPrefix string `json:"stream_prefix"`
	}
	type ReportParam struct {
		getStreamPrefixParam *GetStreamPrefixParam `json:"get_stream_prefix_param"`

		EventName string                 `json:"event_name"`
		Data      map[string]interface{} `json:"data"`
	}
	var re = &ReportParam{getStreamPrefixParam: &GetStreamPrefixParam{}}

	SetExtra(E().Interface("inter1", "1").Interface("inter2", re)).Ctx(nil).
		Info("print empty param  %+v", re)

	re.getStreamPrefixParam.AppId = 1
	SetExtra(E().Interface("inter1", "1").Interface("inter2", re)).Ctx(nil).
		Info("print empty param  %+v", re)

	ctxc, cancel := context.WithCancel(context.TODO())
	cancel()
	Ctx(ctxc).Info("print canceled context")

	Ctx(nil).Info("print nil context")
}

func Test_toShortCaller(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{"github.com/abc/def/ghi/jkl"}, "def/ghi/jkl"},
		{"test2", args{"github.com/abc/def/ghi/jkl.go"}, "def/ghi/jkl.go"},
		{"test3", args{"github.com/abc/def/ghi/jkl.go:123"}, "def/ghi/jkl.go:123"},
		{"test3-1", args{"github.com/abc/def/ghi/jkl.go@ccc"}, "def/ghi/jkl.go@ccc"},
		{"test4", args{"github.com/abc/def/ghi/jkl.go:123:456"}, "def/ghi/jkl.go:123:456"},

		{"testf1", args{"jkl"}, "jkl"},
		{"testf2", args{"ghi/jkl"}, "ghi/jkl"},
		{"testf3", args{"abc/ghi/jkl"}, "abc/ghi/jkl"},
		{"testf4", args{"github.com/abc/ghi/jkl"}, "abc/ghi/jkl"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toShortCaller(tt.args.line); got != tt.want {
				t.Errorf("toShortCaller() = %v, want %v", got, tt.want)
			}
		})
	}
}
