package logger

type LoggerOptFunc func(*loggerOpt)

type loggerOpt struct {
	logDir        string
	debugMode     bool
	verboseMode   bool
	stdoutHttpLog bool
	stdoutNatsLog bool
}

func LogDir(path string) LoggerOptFunc {
	return func(opt *loggerOpt) {
		opt.logDir = path
	}
}

func DebugMode(enable bool) LoggerOptFunc {
	return func(opt *loggerOpt) {
		opt.debugMode = enable
	}
}

func VerboseMode(enable bool) LoggerOptFunc {
	return func(opt *loggerOpt) {
		opt.verboseMode = enable
	}
}

func StdoutHttpLog(enable bool) LoggerOptFunc {
	return func(opt *loggerOpt) {
		opt.stdoutHttpLog = enable
	}
}

func StdoutNatsLog(enable bool) LoggerOptFunc {
	return func(opt *loggerOpt) {
		opt.stdoutNatsLog = enable
	}
}
