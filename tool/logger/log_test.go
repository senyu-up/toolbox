package logger

import (
	"context"
	"errors"
	"fmt"
	"github.com/senyu-up/toolbox/tool/trace"
	"os"
	"testing"
	"time"
)

var p = `{
	"Console": {
		"level": "DEBG",
		"color": true
	},
	"File": {
		"filename": "app.log",
		"level": "EROR",
		"daily": true,
		"maxlines": 1000000,
		"maxsize": 256,
		"maxdays": -1,
		"append": true,
		"permit": "0660"
	}
}`

func logAllLevel(l Log) {
	l.Trace("this is Trace")
	l.Debug("this is Debug")
	l.Info("this is Info")
	l.Warn("this is Warn")
	l.Error("this is Error")
	l.Crit("this is Critical")
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "logger panic err %v", r)
		}
	}()
	l.Alert("this is Alert")
	l.Emer("this is Emergency")
	l.Panic("this is Panic")
}

func logExtraFields(l Log) {
	var ctx = trace.NewTrace()
	l.Ctx(ctx).Info("log ctx")

	traceId, spanId := trace.NewTraceID(), trace.NewSpanID()
	l.TraceId(traceId).SpanId(spanId).Trace("log traceId, spanId")

	var err = errors.New("This is err ")
	l.SetErr(err).Error("I got a err", err)

	l.SetExtra(E().String("name", "xh").String("level", "2")).Info("log extra fields")

	l.SetExtra(E().Bytes("name", []byte{'x', 'h'})).Info("log extra fields") // ?
}

func TestLogDrivers(t *testing.T) {
	var adapters = []string{AdapterZap, AdapterFile, AdapterConsole}
	for _, ad := range adapters {
		SwitchLogger(ad)
		logAllLevel(GetLogger())
		logExtraFields(GetLogger())
	}

	if _, err := SwitchLogger(AdapterConn); err == nil {
		t.Errorf("SwitchLogger should return err")
	} else {
		t.Logf("SwitchLogger conn failed: %v", err)
	}
	SwitchLogger(AdapterConsole)
}

func TestSwitchLogger(t *testing.T) {
	adapters = make(map[string]Driver)
	var adapterNames = []string{AdapterZap, AdapterFile, AdapterConsole}
	for _, ad := range adapterNames {
		SwitchLogger(ad)
		logAllLevel(GetLogger())
		logExtraFields(GetLogger())
	}
}

func TestLogOut(t *testing.T) {
	//setLogger(p)
	Trace("this is Trace")
	Debug("this is Debug")
	Info("this is Info")
	Warn("this is Warn")
	Error("this is Error")
	Crit("this is Critical")
	Alert("this is Alert")
	Emer("this is Emergency")
}

func TestLogConfigReload(t *testing.T) {

	for {
		Trace("this is Trace")
		Debug("this is Debug")
		Info("this is Info")
		Warn("this is Warn")
		Error("this is Error")
		Crit("this is Critical")
		Alert("this is Alert")
		Emer("this is Emergency")
		fmt.Println()

		time.Sleep(time.Millisecond)
		break
	}

}

func TestRequestId(t *testing.T) {
	SwitchLogger(AdapterZap)
	var traceCtx = trace.NewTrace()
	Ctx(traceCtx).Info("print log with new trace ctx")

	var newCtx = trace.NewContextWithRequestIdAndSpanId(context.TODO(), "req_111", "span_222")
	Ctx(newCtx).Info("print log with set trace,span ctx")
}

func TestLogger(t *testing.T) {
	//INFO
	err := errors.New("OOM")
	Info("this", err)
}

func Test_formatLog(t *testing.T) {
	type args struct {
		f interface{}
		v []interface{}
	}

	arr := []string{"foo", "bar", "baz"}
	var ifaceArr []interface{}
	for _, s := range arr {
		ifaceArr = append(ifaceArr, s)
	}
	var iterArr2 []interface{}
	iterArr2 = append(iterArr2, arr)
	var iterArr3 []interface{}
	iterArr3 = append(iterArr3, 3.0)
	var iterArr4 []interface{}
	iterArr4 = append(iterArr4, 3, 4)
	var iterArr5 []interface{}
	iterArr5 = append(iterArr5, 'a', 'b', 'c')
	var iterArr6 []interface{}
	iterArr6 = append(iterArr6, "str")

	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "1", args: args{"test", ifaceArr}, want: "test foo bar baz"},
		{name: "2", args: args{"test %v", iterArr2}, want: "test [foo bar baz]"},
		{name: "3", args: args{"f %.3f", iterArr3}, want: "f 3.000"},
		{name: "4", args: args{"int %d", iterArr4}, want: "int 3 4"},
		{name: "5", args: args{"byte ", iterArr5}, want: "byte  97 98 99"},
		{name: "5-1", args: args{"byte %c", iterArr5}, want: "byte a 98 99"},
		{name: "6", args: args{"string %v %s", iterArr6}, want: "string str %!s(MISSING)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatLog(tt.args.f, tt.args.v...); got != tt.want {
				t.Errorf("formatLog() = %v, want %v", got, tt.want)
			}
		})
	}
}
