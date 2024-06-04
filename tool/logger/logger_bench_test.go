package logger

import (
	"errors"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/trace"
	"testing"
	"time"
)

func BenchmarkBaseLogger_ErrorForConsole(b *testing.B) {
	InitDefaultLoggerByConf(&config.LogConfig{
		AppName:    "test_bench",
		TimeFormat: time.ANSIC,
	})
	SwitchLogger(AdapterConsole)
	var logger = GetLogger()
	var err = errors.New("This is err ")
	var ctx = trace.NewTrace()
	type args struct {
		err    error
		field  *Extras
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			// BenchmarkBaseLogger_ErrorForConsole/1-8 	1000000000	         0.0000141 ns/op
			name: "1", args: args{err: err, format: "Got Err %v", args: []interface{}{err}, field: nil},
		},
		{
			// BenchmarkBaseLogger_ErrorForConsole/2-8 	1000000000	         0.0000108 ns/op
			name: "2", args: args{err: err, format: "Got Err "},
		},
		{
			// BenchmarkBaseLogger_ErrorForConsole/3-8 	1000000000	         0.0000076 ns/op
			name: "3", args: args{err: err, format: "Got Err %v", args: []interface{}{err}, field: E().String("name", "xh")},
		},
	}

	b.ResetTimer()
	for _, tt := range tests {
		b.Run(tt.name, func(t *testing.B) {
			var l = logger.Ctx(ctx).SetErr(err)
			if nil != tt.args.field {
				l.SetExtra(tt.args.field)
			}
			l.Error(tt.args.format, tt.args.args...)
		})
	}
	b.StopTimer()
}

func BenchmarkBaseLogger_ErrorForZap(b *testing.B) {
	SwitchLogger(AdapterZap)
	//SwitchLogger(AdapterConsole)
	var logger = GetLogger()
	var err = errors.New("This is err ")
	var ctx = trace.NewTrace()
	type args struct {
		err    error
		field  *Extras
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			// BenchmarkBaseLogger_ErrorForZap/1-8 	1000000000	         0.0000225 ns/op
			name: "1", args: args{err: err, format: "Got Err %v", args: []interface{}{err}, field: nil},
		},
		{
			// BenchmarkBaseLogger_ErrorForZap/2-8 	1000000000	         0.0000146 ns/op
			name: "2", args: args{err: err, format: "Got Err "},
		},
		{
			// BenchmarkBaseLogger_ErrorForZap/3-8 	1000000000	         0.0000206 ns/op
			name: "3", args: args{err: err, format: "Got Err %v", args: []interface{}{err}, field: E().String("name", "xh").Any("haha", map[string]interface{}{"a": 1, "b": 2})},
		},
	}

	b.ResetTimer()
	for _, tt := range tests {
		b.Run(tt.name, func(t *testing.B) {
			var l = logger.Ctx(ctx).SetErr(err)
			if nil != tt.args.field {
				l.SetExtra(tt.args.field)
			}
			l.Info(tt.args.format, tt.args.args...)
		})
	}
	b.StopTimer()
}

// best use
//
// BenchmarkBaseLogger_Error_Best
//
//	@Description: 基准测试，最好的情况，不要额外参数——不触发fmt.Sprintf(), 不传入ctx不获取trace
//	@param b  body any true "-"
func BenchmarkBaseLogger_Error_Best(b *testing.B) {
	SwitchLogger(AdapterZap)
	var logger = GetLogger()
	var err = errors.New("This is err ")
	type args struct {
		err    error
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			// BenchmarkBaseLogger_Error_Best/2-8 	1000000000	         0.0000118 ns/op
			name: "2", args: args{err: err, format: "Got Err "},
		},
	}

	b.ResetTimer()
	for _, tt := range tests {
		b.Run(tt.name, func(t *testing.B) {
			logger.SetErr(err).Error(tt.args.format)
		})
	}
	b.StopTimer()
}
