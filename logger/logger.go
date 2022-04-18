package logger

import (
	"io"
	"log"

	"github.com/comail/colog"
)

type multiLogger struct {
	std io.Writer
	sub io.Writer
}

func (m *multiLogger) Write(p []byte) (int, error) {
	if m.sub == nil {
		return m.std.Write(p)
	}

	_, err := m.std.Write(p)
	if err != nil {
		return -1, err
	}
	return m.sub.Write(p)
}

type Logger interface {
	Rotate()
	Init(*loggerOpt) error
	Write([]byte) (int, error)
}

var (
	loggers = make(map[string]Logger)
)

func addLogger(name string, l Logger) {
	loggers[name] = l
}

func Init(funcs ...LoggerOptFunc) error {
	opt := new(loggerOpt)
	for _, fn := range funcs {
		fn(opt)
	}

	for _, v := range loggers {
		if err := v.Init(opt); err != nil {
			return err
		}
	}
	return nil
}

func RotateLogs() {
	for _, v := range loggers {
		v.Rotate()
	}
}

func GetLogger(name string) Logger {
	return loggers[name]
}

func newLogger(opt *loggerOpt, out io.Writer, prefix string, flags int) *log.Logger {
	c := colog.NewCoLog(out, prefix, flags)
	c.SetDefaultLevel(colog.LDebug)
	c.SetMinLevel(colog.LInfo)
	if opt.debugMode {
		c.SetMinLevel(colog.LDebug)
		if opt.verboseMode {
			c.SetMinLevel(colog.LTrace)
		}
	}
	return c.NewLogger()
}
