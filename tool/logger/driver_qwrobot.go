package logger

import (
	"github.com/senyu-up/toolbox/tool/wework/qwrobot"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

var fieldEncoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
	TimeKey:        "ts",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	FunctionKey:    zapcore.OmitKey,
	MessageKey:     "msg",
	StacktraceKey:  "stack",
	LineEnding:     "",
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     zapcore.TimeEncoderOfLayout(LogTimeDefaultFormat),
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.FullCallerEncoder,
})

type QWRobot struct {
	robot      *qwrobot.QWRobot
	LogLevel   LogLevel
	callerSkip int
}

func NewQwRobot(r *qwrobot.QWRobot, opts ...QWRobotOption) *QWRobot {
	var qwr = &QWRobot{robot: r, callerSkip: DefaultCallDepth + 1} // qw robot call skip 加一层
	for _, opt := range opts {
		opt(qwr)
	}
	return qwr
}

func (Q *QWRobot) LogWrite(when time.Time, msg string, level LogLevel, extras []Field) error {
	fields := fieldsToZapFieldsV2(when, extras, Q.callerSkip)

	var qMsg = qwrobot.Message{
		Title:   msg,
		Content: fieldsToString(fields),
	}
	switch level {
	case LevelError, LevelCritical, LevelEmergency, LevelAlert:
		Q.robot.Error(qMsg)
	case LevelWarning:
		Q.robot.Warn(qMsg)
	case LevelInformational:
		Q.robot.Info(qMsg)
	}
	return nil
}

func (Q *QWRobot) CallerSkip(skip int) {
	Q.callerSkip = skip
	return
}

func (Q *QWRobot) Destroy() {
}

func (Q *QWRobot) Name() string {
	return "qw_robot"
}

func (Q *QWRobot) CurrentLevel() LogLevel {
	// 在 callback 中，目测没用到这个方法
	return Q.LogLevel
}

func fieldsToString(fields []zap.Field) string {
	if fields == nil || len(fields) == 0 {
		return ""
	}
	en := fieldEncoder.Clone()
	entry := zapcore.Entry{}
	en.AddTime("ts", time.Now())
	buff, _ := en.EncodeEntry(entry, fields)

	s := buff.String()
	if pos := strings.Index(s, "\t\t"); pos != -1 {
		s = s[pos+2 : len(s)-1]
	}

	return s
}
