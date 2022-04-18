package logger

import (
	"fmt"
	"log"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/nats-io/nats-server/v2/server"
)

type NatsLogger struct {
	r  *rotatelogs.RotateLogs
	cl *log.Logger
}

func (l *NatsLogger) ServerLogger() server.Logger {
	return l
}

func (l *NatsLogger) Rotate() {
	l.r.Rotate()
}

func (l *NatsLogger) Write(p []byte) (int, error) {
	return l.r.Write(p)
}

func (l *NatsLogger) Init(opt *loggerOpt) error {
	rotate, err := rotatelogs.New(
		opt.logDir+"/nats_log.%Y%m%d",
		rotatelogs.WithRotationTime(1*time.Minute),
		rotatelogs.WithMaxAge(-1),
	)
	if err != nil {
		return err
	}

	l.r = rotate
	l.cl = newLogger(opt, rotate, "nats ", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}

func (l *NatsLogger) Noticef(format string, v ...interface{}) {
	l.cl.Printf("info: %s", fmt.Sprintf(format, v...))
}

func (l *NatsLogger) Warnf(format string, v ...interface{}) {
	l.cl.Printf("warn: %s", fmt.Sprintf(format, v...))
}

func (l *NatsLogger) Fatalf(format string, v ...interface{}) {
	l.cl.Printf("error: %s", fmt.Sprintf(format, v...))
}

func (l *NatsLogger) Errorf(format string, v ...interface{}) {
	l.cl.Printf("error: %s", fmt.Sprintf(format, v...))
}

func (l *NatsLogger) Debugf(format string, v ...interface{}) {
	l.cl.Printf("debug: %s", fmt.Sprintf(format, v...))
}

func (l *NatsLogger) Tracef(format string, v ...interface{}) {
	l.cl.Printf("trace: %s", fmt.Sprintf(format, v...))
}

func init() {
	addLogger("nats", new(NatsLogger))
}
