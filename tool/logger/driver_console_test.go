package logger

import (
	"fmt"
	"github.com/senyu-up/toolbox/tool/config"
	"testing"
)

// Try each log level in decreasing order of priority.
func testConsoleCalls(bl Log) {
	bl.Emer("emergency")
	bl.Alert("alert")
	bl.Crit("critical")
	bl.Error("error")
	bl.Warn("warning")
	bl.Debug("notice")
	bl.Info("informational")
	bl.Trace("trace")
}

// Try each log level in decreasing order of priority.
func testStaticCalls() {
	Emer("emergency")
	Alert("alert")
	Crit("critical")
	Error("error")
	Warn("warning")
	Debug("notice")
	Info("informational")
	Trace("trace")
}

func TestConsole(t *testing.T) {
	var conf = &config.LogConfig{
		DefaultLog: AdapterConsole,
		AppName:    "test_console",
		CallDepth:  4,
		Console: &config.ConsoleConfig{
			Level:    "DEBG",
			Colorful: true,
		},
	}
	InitDefaultLoggerByConf(conf)
	log1 := GetLogger()
	testConsoleCalls(log1)
	fmt.Print("\n\n")
	testStaticCalls()
	fmt.Print("\n\n")

	conf.Console.Level = "EROR"
	InitDefaultLoggerByConf(conf)
	log2 := GetLogger()
	testConsoleCalls(log2)
	fmt.Print("\n\n")
	testStaticCalls()
}

// Test console without color
func TestNoColorConsole(t *testing.T) {
	var conf = &config.LogConfig{
		DefaultLog: AdapterConsole,
		AppName:    "test_console",
		CallDepth:  3,
		UsePath:    "/Users/th/xh/",
		Console: &config.ConsoleConfig{
			Level:    "DEBG",
			Colorful: false,
		},
	}
	InitDefaultLoggerByConf(conf)
	log := GetLogger()
	testConsoleCalls(log)
	fmt.Print("\n\n")
	testStaticCalls()

}
