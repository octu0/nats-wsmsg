package logger

import (
	"os"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
)

type GeneralLogger struct {
	r *rotatelogs.RotateLogs
	m *multiLogger
}

func (l *GeneralLogger) Rotate() {
	l.r.Rotate()
}

func (l *GeneralLogger) Write(p []byte) (int, error) {
	return l.m.Write(p)
}

func (l *GeneralLogger) Init(opt *loggerOpt) error {
	rotate, err := rotatelogs.New(
		opt.logDir+"/general_log.%Y%m%d",
		rotatelogs.WithRotationTime(1*time.Minute),
		rotatelogs.WithMaxAge(-1),
	)
	if err != nil {
		return err
	}

	l.r = rotate
	l.m = &multiLogger{
		std: rotate,
		sub: os.Stdout,
	}

	return nil
}

func init() {
	addLogger("general", new(GeneralLogger))
}
