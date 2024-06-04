package logger

import (
	"context"
	"errors"
	"runtime"
	"sync"

	"github.com/senyu-up/toolbox/tool/config"
)

var (
	logInst     *baseLogger = nil
	defaultConf             = config.LogConfig{}
	logLock                 = sync.RWMutex{}
)

// initDefaultLogger
//
//	@Description: 初始化一个默认 logger， 如果没有指定配置，就会初始化一个 console logger
//	@param opts  body any true "-"
//	@return Log
func initDefaultLogger(opts ...LogOption) *baseLogger {
	l := baseLogger{}

	// appName用于记录网络传输时标记的程序发送方，
	// 通过环境变量APPSN进行设置,默认为NONE,此时无法通过网络日志检索区分不同服务发送方
	if appSn == "" {
		appSn = "NONE"
	}
	l.appName = "[" + appSn + "]"
	l.callDepth = DefaultCallDepth
	l.ShowCallerLevel = LevelTrace
	l.adapterName = AdapterConsole
	l.timeFormat = LogTimeDefaultFormat
	l.extra = NewExtras()

	for _, opt := range opts {
		opt(&l)
	}

	var defaultAd = &Console{LogLevel: Level(), Colorful: runtime.GOOS != "windows"}
	register(AdapterConsole, defaultAd)
	l.adapter = defaultAd
	l.output = &nameLogger{name: AdapterConsole, config: ""}

	logLock.Lock()
	defer logLock.Unlock()
	logInst = &l
	return &l
}

// 通过 adapter name 切换全局 logger， 并返回当前 logger 对象
func SwitchLogger(adapter string) (Log, error) {
	var theLog = &baseLogger{}
	if logInst != nil {
		theLog = logInst
	}
	theLog.extra = NewExtras()
	theLog.adapterName = adapter
	if ad, err := GetAdapter(adapter); err != nil {
		// 如果没初始化，需要手动初始化个默认的
		switch adapter {
		case AdapterZap:
			var zapLogger = &Zap{LogLevel: Level(), CallerSkip: 3}
			zapLogger.InitByConf(config.ZapConfig{})
			register(adapter, zapLogger)
			theLog.output = &nameLogger{name: adapter, config: ""}
			theLog.adapter = zapLogger
		case AdapterConsole:
			var consoleLogger = &Console{
				LogLevel: Level(),
				Colorful: runtime.GOOS != "windows",
			}
			consoleLogger.InitByConf(config.ConsoleConfig{})
			register(AdapterConsole, consoleLogger)
			theLog.output = &nameLogger{name: adapter, config: ""}
			theLog.adapter = consoleLogger
		case AdapterFile:
			var fileLogger = &File{
				Daily:      true,
				MaxDays:    7,
				Append:     true,
				LogLevel:   Level(),
				PermitMask: "0777",
				MaxLines:   10,
				MaxSize:    10 * 1024 * 1024,
			}
			fileLogger.InitByConf(config.FileConfig{})
			register(AdapterFile, fileLogger)
			theLog.output = &nameLogger{name: adapter, config: ""}
			theLog.adapter = fileLogger
		default:
			// conn 必须要配置初始化, 没有缺省值
			return theLog, errors.New("Set logger adapter by name not found, and can not get default one ")
		}
	} else {
		// 如果找到了，切过去
		theLog.output = &nameLogger{name: adapter, config: ""}
		theLog.adapter = ad
	}

	logLock.Lock()
	defer logLock.Unlock()
	logInst = theLog
	return &LocalLogger{bl: theLog}, nil
}

// 通过注册第三方 log driver
func SetLogger(adapterName string, logDriver Driver) error {
	adapters[adapterName] = logDriver
	return nil
}

func SetCallBack(callBack Driver) error {
	if logInst == nil {
		initDefaultLogger()
	}
	logLock.Lock()
	defer logLock.Unlock()
	logInst.SetCallBack(callBack)
	return nil
}

// InitLoggerByConf
//
//	@Description: 通过配置初始化 logger
//	@param conf  body any true "-"
//	@return log
//	@return err
func InitLoggerByConf(conf *config.LogConfig) (log Log, err error) {
	var (
		adapter = ""
		driver  Driver
		theLog  = &baseLogger{}
		ok      bool
	)
	if conf.CallDepth < 1 {
		conf.CallDepth = DefaultCallDepth
	}
	if conf.Zap != nil {
		conf.Zap.CallerSkip = conf.CallDepth
		var zapLogger = &Zap{LogLevel: Level(), CallerSkip: DefaultCallDepth}
		if err = zapLogger.InitByConf(*conf.Zap); err != nil {
			return
		}
		register(AdapterZap, zapLogger)
		adapter = AdapterZap
	}
	if conf.File != nil {
		var fileLogger = &File{
			Daily:      true,
			MaxDays:    7,
			Append:     true,
			LogLevel:   Level(),
			Level:      "INFO",
			PermitMask: "0777",
			MaxLines:   10,
			MaxSize:    10 * 1024 * 1024,
		}
		if err = fileLogger.InitByConf(*conf.File); err != nil {
			return
		}
		register(AdapterFile, fileLogger)
		adapter = AdapterFile
	}
	if conf.Console != nil {
		var consoleLogger = &Console{LogLevel: Level(), Colorful: conf.Console.Colorful}
		if err = consoleLogger.InitByConf(*conf.Console); err != nil {
			return
		}
		register(AdapterConsole, consoleLogger)
		adapter = AdapterConsole
	}
	if conf.Conn != nil {
		var connLogger = &ConnLogger{LogLevel: Level()}
		if err = connLogger.InitByConf(*conf.Conn); err != nil {
			return
		}
		register(AdapterConn, connLogger)
		adapter = AdapterConn
	}

	if driver, ok = adapters[conf.DefaultLog]; ok {
		// 如果设置了 default，按照 default 得取
		adapter = conf.DefaultLog
	} else if 0 < len(adapter) {
		driver, ok = adapters[adapter]
	} else {
		for adapter, driver = range adapters {
			// 取第一个
			break
		}
	}

	// 初始化 baseLogger
	var opts = []LogOption{LogOptWithTimeFormat(conf.TimeFormat),
		LogOptWithAppName(conf.AppName),
		LogOptWithCallDepth(conf.CallDepth),
		LogOptWithUsePath(conf.UsePath)}
	if 0 < len(adapter) {
		// 如果有效
		theLog = newBaseLogger(adapter, driver, opts...)
	} else {
		theLog = initDefaultLogger(opts...)
	}
	return &LocalLogger{bl: theLog}, err
}

func InitDefaultLoggerByConf(conf *config.LogConfig) (log Log, err error) {
	theLog, err := InitLoggerByConf(conf)
	if err != nil {
		return nil, err
	}
	logLock.Lock()
	defer logLock.Unlock()
	if localLog, ok := theLog.(*LocalLogger); ok {
		logInst = localLog.bl
	} else {
		return nil, errors.New("init default logger failed")
	}
	return theLog, nil
}

func GetLogger() Log {
	if logInst == nil {
		initDefaultLogger()
	}
	return &LocalLogger{bl: logInst}
}

func cloneLogger() Log {
	if logInst == nil {
		initDefaultLogger()
	}
	return logInst.Clone()
}

func Ctx(c context.Context) Log {
	var newLogger = cloneLogger()
	newLogger.Ctx(c)

	return newLogger
}

func TraceId(t string) Log {
	var newLogger = cloneLogger()
	newLogger.TraceId(t)

	return newLogger
}

func SetExtra(e *Extras) Log {
	var newLogger = cloneLogger()
	newLogger.SetExtra(e)

	return newLogger
}

func SpanId(s string) Log {
	var newLogger = cloneLogger()
	newLogger.SpanId(s)

	return newLogger
}

func SetErr(e error) Log {
	var newLogger = cloneLogger()
	newLogger.SetErr(e)

	return newLogger
}

func Notify() Log {
	var newLogger = cloneLogger()
	newLogger.Notify()

	return newLogger
}

func Emer(format string, args ...interface{}) {
	GetLogger().Emer(format, args...)
}

func Alert(format string, args ...interface{}) {
	GetLogger().Alert(format, args...)
}

func Crit(format string, args ...interface{}) {
	GetLogger().Crit(format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Error(format, args...)
}

func Warn(format string, args ...interface{}) {
	GetLogger().Warn(format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Info(format, args...)
}

func Trace(format string, args ...interface{}) {
	GetLogger().Trace(format, args...)
}

func Debug(format string, args ...interface{}) {
	GetLogger().Debug(format, args...)
}
