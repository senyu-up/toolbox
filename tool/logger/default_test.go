package logger

import (
	"fmt"
	"github.com/senyu-up/toolbox/tool/config"
	"testing"
	"time"
)

type MyNotify struct{}

func (m MyNotify) LogWrite(when time.Time, msg string, level LogLevel, extra []Field) error {
	fmt.Printf("Hi it's MyNotify, now is %s level %d Msg: %s \n extra %v\n", when.String(), level, msg, extra)
	return nil
}

func (m MyNotify) Destroy() {
}

func (m MyNotify) Name() string {
	return "my_notify"
}

func (m MyNotify) CurrentLevel() LogLevel {
	return LevelTrace
}

func TestSetLogger(t *testing.T) {
	SetLogger("my_logger", &MyNotify{})
	log, err := SwitchLogger("my_logger")
	if err != nil {
		t.Errorf("%v", err)
	}
	log.Info(" [this is my logger log] ")

	SwitchLogger(AdapterConsole)
}

func Test_initDefaultLogger(t *testing.T) {
	log := initDefaultLogger(LogOptWithShowCallerLevel(LevelTrace), LogOptWithUsePath("./"),
		LogOptWithCallBack(MyNotify{}))
	log.Notify().Warn("this is notify Warn")
}

func SetUp_Init() {
	logInst = nil
	adapters = make(map[string]Driver)
}

func TestInitLoggerByConf(t *testing.T) {
	SetUp_Init()
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
	if _, err := GetAdapter(AdapterZap); err != nil {
		t.Errorf("%v", err)
	}
	if _, err := GetAdapter(AdapterFile); err == nil {
		t.Errorf("%v", err)
	}
}

func TestInitLoggerByConf2(t *testing.T) {
	SetUp_Init()
	var conf = &config.LogConfig{
		AppName:   "test_file",
		CallDepth: 4,

		Zap: &config.ZapConfig{
			Level:      "TRAC",
			Colorful:   true,
			CallerSkip: 4,
		},
		File: &config.FileConfig{
			Filename:   "test_file.log",
			Level:      "DEBG",
			Append:     true,
			Daily:      true,
			PermitMask: "0770",
			MaxLines:   1000,
			MaxDays:    1,
		},
	}
	InitDefaultLoggerByConf(conf)
	if _, err := GetAdapter(AdapterZap); err != nil {
		t.Errorf("%v", err)
	}
	if _, err := GetAdapter(AdapterFile); err != nil {
		t.Errorf("%v", err)
	}
	if _, err := GetAdapter(AdapterConn); err == nil {
		t.Errorf("%v", err)
	}
}

func TestInitLoggerByConf3(t *testing.T) {
	SetUp_Init()
	var conf = &config.LogConfig{
		AppName:   "test_file",
		CallDepth: 4,

		File: &config.FileConfig{
			Filename:   "test_file.log",
			Level:      "DEBG",
			Append:     true,
			Daily:      true,
			PermitMask: "0770",
			MaxLines:   1000,
			MaxDays:    1,
		},
	}
	InitDefaultLoggerByConf(conf)
	if _, err := GetAdapter(AdapterZap); err == nil {
		t.Errorf("%v", err)
	}
	if _, err := GetAdapter(AdapterFile); err != nil {
		t.Errorf("%v", err)
	}
	if _, err := GetAdapter(AdapterConn); err == nil {
		t.Errorf("%v", err)
	}
}
