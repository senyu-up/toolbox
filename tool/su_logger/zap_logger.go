package su_logger

import (
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"runtime"
)

var loggerWithCaller = NewZapLogger(`{"caller_skip":1}`)
var loggerWithoutCaller = NewZapLogger(`{}`)

const logTimeDefaultFormat = "2006-01-02 15:04:05.000"

// 日志格式
func (z *zapLogger) getEncoder() zapcore.Encoder {
	var encodeLevel zapcore.LevelEncoder
	if z.Colorful {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encodeLevel = zapcore.CapitalLevelEncoder
	}
	encodingConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey, //zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     zapcore.TimeEncoderOfLayout(logTimeDefaultFormat),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, //zapcore.FullCallerEncoder zapcore.ShortCallerEncoder
		NewReflectedEncoder: func(w io.Writer) zapcore.ReflectedEncoder {
			encoder := jsoniter.NewEncoder(w)
			encoder.SetEscapeHTML(false)

			return encoder
		},
	}

	return zapcore.NewConsoleEncoder(encodingConfig)
}

// 日志写到哪里
func (z *zapLogger) getWriteSyncer() (zapcore.WriteSyncer, error) {
	if z.Output == "" || z.Output == "std" {
		return zapcore.AddSync(os.Stdout), nil
	} else {
		file, err := os.OpenFile(z.Output, os.O_CREATE|os.O_APPEND, 666)
		if err != nil {
			return nil, err
		}
		return zapcore.AddSync(file), err
	}
}

type zapLogger struct {
	zapInst  *zap.Logger
	Level    string `json:"level"`
	Colorful bool   `json:"color"`
	//Writer io.Writer
	// std(默认) or 具体的文件路径
	Output string `json:"output"`
	// 调用栈往上走的层数
	CallerSkip int `json:"caller_skip"`
}

// 移动到 toolbox/tool/logger, 统一使用 logger 作为日志函数，支持zap作为底层日志实现
// deprecated
func NewZapLogger(jsonConfig string) (l *zap.Logger) {
	z := &zapLogger{}
	_ = z.init(jsonConfig)

	encoder := z.getEncoder()
	writer, _ := z.getWriteSyncer()

	core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)
	if z.CallerSkip > 0 {
		//  仅当设置了caller才会获取调用栈信息
		l = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(z.CallerSkip))
	} else {
		l = zap.New(core)
	}

	return l
}

func (z *zapLogger) init(jsonConfig string) error {
	if jsonConfig == "" {
		jsonConfig = "{}"
	}
	if jsonConfig != "{}" {
		_, _ = fmt.Fprintf(os.Stdout, "zap logger init with config:%s\n", jsonConfig)
	}

	err := json.Unmarshal([]byte(jsonConfig), z)
	if runtime.GOOS == "windows" {
		z.Colorful = false
	} else {
		z.Colorful = true
	}

	return err
}
