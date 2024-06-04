package logger

type LogOption func(*baseLogger)

func LogOptWithAppName(ap string) LogOption {
	return func(option *baseLogger) {
		option.appName = ap
	}
}

func LogOptWithShowCallerLevel(scl LogLevel) LogOption {
	return func(option *baseLogger) {
		option.ShowCallerLevel = scl
	}
}

func LogOptWithCallDepth(cd int) LogOption {
	return func(option *baseLogger) {
		option.callDepth = cd
	}
}

func LogOptWithUsePath(up string) LogOption {
	return func(option *baseLogger) {
		option.usePath = up
	}
}

func LogOptWithTimeFormat(tf string) LogOption {
	return func(option *baseLogger) {
		if 0 < len(tf) {
			option.timeFormat = tf
		}
	}
}

func LogOptWithCallBack(cb Driver) LogOption {
	return func(option *baseLogger) {
		option.callBack = cb
	}
}

type QWRobotOption func(*QWRobot)

func QWRobotOptWithCallerSkip(skip int) QWRobotOption {
	return func(option *QWRobot) {
		option.callerSkip = skip
	}
}
