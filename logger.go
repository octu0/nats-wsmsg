package wsmsg

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/comail/colog"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/nats-io/nats-server/v2/server"
)

type MultiLogger struct {
	std io.Writer
	sub io.Writer
}

func (m *MultiLogger) Write(p []byte) (int, error) {
	if m.sub == nil {
		return m.std.Write(p)
	}

	_, err := m.std.Write(p)
	if err != nil {
		return -1, err
	}
	return m.sub.Write(p)
}

const TAB string = "\t"

type HttpLogger interface {
	Write(host string, method string, uri string, status int, ua string)
}

type DefLogger struct {
	m  *MultiLogger
	r  *rotatelogs.RotateLogs
	cl *log.Logger
}

func NewHttpLogger(config Config) HttpLogger {
	rotate, err := rotatelogs.New(
		config.LogDir+"/http_log.%Y%m%d", // TODO
		rotatelogs.WithRotationTime(1*time.Minute),
		rotatelogs.WithMaxAge(-1),
	)
	if err != nil {
		log.Fatalf("error: http file logger creation failed: %s", err.Error())
	}

	multi := new(MultiLogger)
	multi.std = rotate
	if config.HttpLogStdout {
		multi.sub = os.Stdout
	}

	c := colog.NewCoLog(multi, "http ", log.Ldate|log.Ltime|log.Lshortfile)
	c.SetDefaultLevel(colog.LDebug)
	c.SetMinLevel(colog.LInfo)
	if config.DebugMode {
		c.SetMinLevel(colog.LDebug)
		if config.VerboseMode {
			c.SetMinLevel(colog.LTrace)
		}
	}

	l := new(DefLogger)
	l.m = multi
	l.r = rotate
	l.cl = c.NewLogger()
	return l
}
func (l *DefLogger) Write(host string, method string, uri string, status int, ua string) {
	msg := []string{
		"host:", host,
		TAB,
		"method:", method,
		TAB,
		"uri:", uri,
		TAB,
		"status:", strconv.Itoa(status),
		TAB,
		"ua:", ua,
	}
	m := strings.Join(msg, "")
	l.cl.Printf("info: %s", m)
}

type NatsLogger struct {
	m  *MultiLogger
	r  *rotatelogs.RotateLogs
	cl *log.Logger
}

func NewNatsLogger(config Config) server.Logger {
	rotate, err := rotatelogs.New(
		config.LogDir+"/nats_log.%Y%m%d", // TODO
		rotatelogs.WithRotationTime(1*time.Minute),
		rotatelogs.WithMaxAge(-1),
	)
	if err != nil {
		log.Fatalf("error: nats file logger creation failed: %s", err.Error())
	}

	multi := new(MultiLogger)
	multi.std = rotate
	if config.NatsLogStdout {
		multi.sub = os.Stdout
	}

	c := colog.NewCoLog(multi, "nats ", log.Ldate|log.Ltime|log.Lshortfile)
	c.SetDefaultLevel(colog.LDebug)
	c.SetMinLevel(colog.LInfo)
	if config.DebugMode {
		c.SetMinLevel(colog.LDebug)
		if config.VerboseMode {
			c.SetMinLevel(colog.LTrace)
		}
	}

	l := new(NatsLogger)
	l.m = multi
	l.r = rotate
	l.cl = c.NewLogger()
	return l
}
func (n *NatsLogger) Noticef(format string, v ...interface{}) {
	label := "info: " + format
	n.cl.Printf(label, v...)
}
func (n *NatsLogger) Warnf(format string, v ...interface{}) {
	label := "warn: " + format
	n.cl.Printf(label, v...)
}
func (n *NatsLogger) Fatalf(format string, v ...interface{}) {
	label := "error: fatal " + format
	n.cl.Printf(label, v...)
}
func (n *NatsLogger) Errorf(format string, v ...interface{}) {
	label := "error: " + format
	n.cl.Printf(label, v...)
}
func (n *NatsLogger) Debugf(format string, v ...interface{}) {
	label := "debug: " + format
	n.cl.Printf(label, v...)
}
func (n *NatsLogger) Tracef(format string, v ...interface{}) {
	label := "trace: " + format
	n.cl.Printf(label, v...)
}
