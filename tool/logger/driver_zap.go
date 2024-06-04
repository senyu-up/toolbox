package logger

import (
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 日志格式
func (z *Zap) getEncoder() zapcore.Encoder {
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
		EncodeTime:     zapcore.TimeEncoderOfLayout(LogTimeDefaultFormat),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, //zapcore.FullCallerEncoder zapcore.ShortCallerEncoder
		NewReflectedEncoder: func(w io.Writer) zapcore.ReflectedEncoder {
			encoder := jsoniter.NewEncoder(w)
			encoder.SetEscapeHTML(false)

			return encoder
		},
	}

	//return zapcore.NewJSONEncoder(encodingConfig)
	return zapcore.NewConsoleEncoder(encodingConfig)
}

// 日志写到哪里
func (z *Zap) getWriteSyncer() (zapcore.WriteSyncer, error) {
	if z.Output == "" || z.Output == "std" {
		return zapcore.AddSync(os.Stdout), nil
	} else {
		var (
			file *os.File
			err  error
		)
		if runtime.GOOS == "darwin" {
			file, err = os.OpenFile(z.Output, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0770) // mac
		} else {
			file, err = os.OpenFile(z.Output, os.O_CREATE|os.O_APPEND, 666)
		}

		if err != nil {
			return nil, err
		}
		return zapcore.AddSync(file), err
	}
}

type Zap struct {
	zapInst  *zap.Logger
	Level    string   `json:"level"`
	Colorful bool     `json:"color"`
	LogLevel LogLevel `json:"log_level"`
	//Writer io.Writer
	// std(默认) or 具体的文件路径
	Output string `json:"output"`
	// 调用栈往上走的层数
	CallerSkip int `json:"caller_skip"`
}

func (z *Zap) InitByConf(conf config.ZapConfig) (err error) {
	copier.CopyWithOption(z, conf, copier.Option{IgnoreEmpty: true})
	z.Colorful = conf.Colorful
	if runtime.GOOS == "windows" || !(z.Output == "std" || z.Output == "") {
		// 如果是 windows，或者输出到文件，则禁用颜色
		z.Colorful = false
	}

	if l, ok := LevelMap[z.Level]; ok {
		z.LogLevel = l
	} else {
		// 获取当前的stage, 并设置对应的日志级别
		stage := os.Getenv(enum.StageKey)
		if stage == enum.EvnStageProduction || stage == "master" {
			z.LogLevel = LevelInformational
		} else {
			z.LogLevel = LevelDebug
		}
	}

	encoder := z.getEncoder()
	writer, err := z.getWriteSyncer()
	if err != nil {
		return err
	}

	var lv zapcore.Level
	if z.LogLevel == LevelInformational {
		lv = zapcore.InfoLevel
	} else {
		lv = zapcore.DebugLevel
	}

	var skip = 2
	if z.CallerSkip > 0 {
		skip = z.CallerSkip
	}
	core := zapcore.NewCore(encoder, writer, lv)
	z.zapInst = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(skip))

	return err
}

func (z *Zap) LogWrite(when time.Time, msg string, level LogLevel, extras []Field) error {
	if level > z.LogLevel {
		return nil
	}
	var msgData = fieldsToZapFields(extras)
	switch level {
	case LevelEmergency:
		z.zapInst.Fatal(msg, msgData...)
	case LevelAlert:
		z.zapInst.Panic(msg, msgData...)
	case LevelCritical:
		z.zapInst.DPanic(msg, msgData...)
	case LevelError:
		z.zapInst.Error(msg, msgData...)
	case LevelWarning:
		z.zapInst.Warn(msg, msgData...)
	case LevelInformational:
		z.zapInst.Info(msg, msgData...)
	default:
		z.zapInst.Debug(msg, msgData...)
	}

	return nil
}

func (z *Zap) Destroy() {
	z.zapInst = nil
}

func (z *Zap) Name() string {
	return AdapterZap
}

func (z *Zap) CurrentLevel() LogLevel {
	return z.LogLevel
}

func toShortCaller(line string) string {
	list := strings.Split(line, "/")
	if len(list) > 2 {
		max := len(list)
		return strings.Join(list[max-3:max], "/")
	} else {
		return strings.Join(list, "/")
	}
}

// fieldsToZapFields
//
//	@Description: 将 Field 转换为 zap.Field
//	@param fields  body any true "-"
//	@return []zap.Field
func fieldsToZapFields(fields []Field) []zap.Field {
	if nil == fields || 1 > len(fields) {
		return nil
	}
	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		var zf = zap.Field{
			Key:       field.Key,
			Type:      zapcore.FieldType(field.Type),
			Integer:   field.Integer,
			Interface: field.Interface,
			String:    field.String,
		}
		if field.Type == BoolType {
			zf = zap.Bool(field.Key, field.Boolean)
		} else if field.Type == ByteStringType {
			zf = zap.ByteString(field.Key, field.Bytes)
		} else if field.Type == Int64Type || field.Type == Uint64Type {
			zf = zap.Int64(field.Key, field.Integer)
		} else if field.Type == Int32Type || field.Type == Uint32Type {
			zf = zap.Int32(field.Key, int32(field.Integer))
		} else if field.Type == Float64Type {
			zf = zap.Float64(field.Key, field.Float)
		} else if field.Type == Float32Type {
			zf = zap.Float32(field.Key, float32(field.Float))
		} else if field.Type == ReflectType {
			zf = zap.Any(field.Key, field.Interface)
		} else if field.Type == ErrorType {
			zf = zap.Any(field.Key, field.Interface)
		}
		zapFields = append(zapFields, zf)
	}

	return zapFields
}

// fieldsToZapFields
//
//	@Description: 将 Field 转换为 zap.Field
//	@param fields  body any true "-"
//	@return []zap.Field
func fieldsToZapFieldsV2(t time.Time, fields []Field, skip int) []zap.Field {
	// 时间点
	var zapFields = append(fieldsToZapFields(fields), zap.Field{
		Key:    "ts",
		Type:   zapcore.StringType,
		String: t.Format("2006-01-02 15:04:05.000"),
	})
	// 触发代码位置
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		caller := toShortCaller(file)
		zapFields = append(zapFields, zap.Field{
			Key:    "caller",
			Type:   zapcore.StringType,
			String: caller + ":" + cast.ToString(line),
		})
	}
	return zapFields
}
