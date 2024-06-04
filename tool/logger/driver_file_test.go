package logger

import (
	"fmt"
	"github.com/senyu-up/toolbox/tool/config"
	"testing"
)

func TestFile_InitByConf(t *testing.T) {
	var conf = &config.LogConfig{
		DefaultLog: AdapterFile,
		AppName:    "test_file",
		CallDepth:  4,
		UsePath:    "/Users/th/xh/",
		File: &config.FileConfig{
			Filename:   "test_file.log",
			Level:      "DEBG",
			Append:     true,
			Daily:      true,
			PermitMask: "0770",
			MaxLines:   1000,
			MaxDays:    1,
		},
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

	conf.File.Level = "EROR"
	InitDefaultLoggerByConf(conf)
	log2 := GetLogger()
	testConsoleCalls(log2)
	fmt.Print("\n\n")
	testStaticCalls()

}
