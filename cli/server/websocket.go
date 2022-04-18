package server

import (
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"gopkg.in/urfave/cli.v1"

	"github.com/octu0/nats-wsmsg"
	"github.com/octu0/nats-wsmsg/cli/clicommon"
	"github.com/octu0/nats-wsmsg/http"
)

func websocketServerAction(c *cli.Context) error {
	parent, err := clicommon.Prepare(c)
	if err != nil {
		return err
	}
	config := clicommon.ValueConfig(parent)
	httpLogger := clicommon.ValueHttpLogger(parent)
	natsLogger := clicommon.ValueNatsLogger(parent)

	listenAddr := net.JoinHostPort(c.String("ip"), c.String("port"))
	listener, err := reuseport.Listen("tcp4", listenAddr)
	if err != nil {
		return err
	}

	log.Printf("info: http server %s", listenAddr)

	ctx, stop := signal.NotifyContext(parent, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ns := server.New(&server.Options{
		Host:         "127.0.0.1",
		Port:         -1,
		HTTPPort:     -1,
		Cluster:      server.ClusterOpts{Port: -1},
		NoLog:        true,
		NoSigs:       true,
		Debug:        config.DebugMode,
		Trace:        config.VerboseMode,
		MaxPayload:   int32(c.Int("max-payload")),
		PingInterval: 5 * time.Millisecond,
		MaxPingsOut:  10,
	})
	ns.SetLogger(natsLogger.ServerLogger(), config.DebugMode, config.VerboseMode)

	go ns.Start()

	if ns.ReadyForConnections(10*time.Second) != true {
		return fmt.Errorf("unable to start a NATS Server")
	}
	defer ns.Shutdown()

	log.Printf("info: local nats started: %s", ns.Addr().String())

	handler := http.Handler(
		fmt.Sprintf("nats://%s", ns.Addr().String()),
		httpLogger,
	)
	s := &fasthttp.Server{
		ReadTimeout:  time.Duration(c.Int("http-read-timeout")) * time.Second,
		WriteTimeout: time.Duration(c.Int("http-write-timeout")) * time.Second,
		IdleTimeout:  time.Duration(c.Int("http-idle-timeout")) * time.Second,
		Handler:      handler,
		Logger:       httpLogger.Logger(),
		Name:         wsmsg.UA,
	}

	go func() {
		<-ctx.Done()

		if err := s.Shutdown(); err != nil {
			log.Fatalf("error: http server shutdown err: %+v", err)
		}
	}()
	return s.Serve(listener)
}

func init() {
	addCommand(cli.Command{
		Name:   "websocket",
		Usage:  "run websocket server",
		Action: websocketServerAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "i, ip",
				Usage:  "server bind-ip",
				Value:  "[0.0.0.0]",
				EnvVar: "WSMSG_BIND_IP",
			},
			cli.StringFlag{
				Name:   "p, port",
				Usage:  "server bind-port",
				Value:  "8080",
				EnvVar: "WSMSG_BIND_PORT",
			},
			cli.IntFlag{
				Name:   "max-payload",
				Usage:  "nats msg max payload size(byte)",
				Value:  1024 * 1024, // 1MB
				EnvVar: "WSMSG_MAX_PAYLOAD",
			},
			cli.IntFlag{
				Name:  "http-read-timeout",
				Usage: "http server read timeout(seconds)",
				Value: 10,
			},
			cli.IntFlag{
				Name:  "http-write-timeout",
				Usage: "http server write timeout(seconds)",
				Value: 10,
			},
			cli.IntFlag{
				Name:  "http-idle-timeout",
				Usage: "http server idle timeout(seconds)",
				Value: 15,
			},
		},
	})
}
