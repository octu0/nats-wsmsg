package logger

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
)

type HttpLogger struct {
	r    *rotatelogs.RotateLogs
	cl   *log.Logger
	pool *sync.Pool
}

func (l *HttpLogger) Logger() *log.Logger {
	return l.cl
}

func (l *HttpLogger) Rotate() {
	l.r.Rotate()
}

func (l *HttpLogger) Write(p []byte) (int, error) {
	return l.r.Write(p)
}

func (l *HttpLogger) ltsv(sb *strings.Builder, key string, value []byte) {
	sb.WriteString(key)
	sb.WriteString(":")
	sb.Write(value)
	sb.WriteString("\t")
}

func (l *HttpLogger) ltsvEnd(sb *strings.Builder, key string, status int) {
	sb.WriteString(key)
	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(status))
	sb.WriteString("\n")
}

func (l *HttpLogger) Accesslog(host, uri, method, ua []byte, status int) {
	sb := l.pool.Get().(*strings.Builder)
	defer func() {
		sb.Reset()
		l.pool.Put(sb)
	}()

	l.ltsv(sb, "host", host)
	l.ltsv(sb, "uri", uri)
	l.ltsv(sb, "method", method)
	l.ltsv(sb, "ua", ua)
	l.ltsvEnd(sb, "status", status)

	l.cl.Printf("info: %s", sb.String())
}

func (l *HttpLogger) Init(opt *loggerOpt) error {
	rotate, err := rotatelogs.New(
		opt.logDir+"/http_log.%Y%m%d",
		rotatelogs.WithRotationTime(1*time.Minute),
		rotatelogs.WithMaxAge(-1),
	)
	if err != nil {
		return err
	}

	l.r = rotate
	l.cl = newLogger(opt, rotate, "http ", log.Ldate|log.Ltime|log.Lshortfile)
	l.pool = &sync.Pool{
		New: func() interface{} {
			return new(strings.Builder)
		},
	}
	return nil
}

func init() {
	addLogger("http", new(HttpLogger))
}
