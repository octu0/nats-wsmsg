package main

import(
  "log"
  "fmt"
  "runtime"
  "context"
  "time"
  "os"
  "os/signal"
  "syscall"

  "github.com/comail/colog"
  "gopkg.in/urfave/cli.v1"
  natsd "github.com/nats-io/gnatsd/server"

  "github.com/octu0/nats-wsmsg"
)

var (
  Commands = make([]cli.Command, 0)
)
func AddCommand(cmd cli.Command){
  Commands = append(Commands, cmd)
}

func action(c *cli.Context) error {
  config := wsmsg.Config{
    DebugMode:    c.Bool("debug"),
    VerboseMode:  c.Bool("verbose"),
    Procs:        c.Int("procs"),
    LogDir:       c.String("log-dir"),
    BindIP:       c.String("ip"),
    BindPort:     c.Int("port"),
  }
  if config.Procs < 1 {
    config.Procs = 1
  }

  if config.DebugMode {
    colog.SetMinLevel(colog.LDebug)
    if config.VerboseMode {
      colog.SetMinLevel(colog.LTrace)
    }
  }

  opts := &natsd.Options{
    Host:         "127.0.0.1",
    Port:         -1,
    HTTPPort:     -1,
    Cluster:      natsd.ClusterOpts{Port: -1},
    NoLog:        true,
    NoSigs:       true,
    Debug:        config.DebugMode,
    Trace:        config.VerboseMode,
    MaxPayload:   c.Int("max-payload"),
    PingInterval: time.Millisecond * time.Duration(wsmsg.DEFAULT_PING_INTERVAL),
    MaxPingsOut:  wsmsg.DEFAULT_PING_OUT,
  }
  ns, err := natsd.NewServer(opts)
  if err != nil {
    log.Printf("error: nats server start failure: %s", err.Error());
    return err
  }
  ns.SetLogger(wsmsg.NewNatsLogger(config), opts.Debug, opts.Trace)

  go ns.Start()

  if ns.ReadyForConnections(10 * time.Second) != true {
    return fmt.Errorf("error: unable to start a NATS Server on %s:%d", opts.Host, opts.Port)
  }
  log.Printf("info: local nats started: %s", ns.Addr().String())

  ctx := context.Background()
  ctx  = context.WithValue(ctx, "config", config)
  ctx  = context.WithValue(ctx, "logger.http", wsmsg.NewHttpLogger(config))
  ctx  = context.WithValue(ctx, "logger.nats", wsmsg.NewNatsLogger(config))
  ctx  = context.WithValue(ctx, "nats.url", fmt.Sprintf("nats://%s", ns.Addr().String()))

  http        := wsmsg.NewHttpServer(ctx)
  error_chan  := make(chan error, 0)
  stopService := func() error {
    sctx, cancel := context.WithTimeout(ctx, 10 * time.Second);
    defer cancel()

    if err := http.Stop(sctx); err != nil {
      log.Printf("error: %s", err.Error())
      return err
    }

    ns.Shutdown()
    return nil
  }

  go func(){
    if err := http.Start(context.TODO()); err != nil {
      error_chan <- err
    }
  }()

  signal_chan := make(chan os.Signal, 10)
  signal.Notify(signal_chan, syscall.SIGTERM)
  signal.Notify(signal_chan, syscall.SIGHUP)
  signal.Notify(signal_chan, syscall.SIGQUIT)
  signal.Notify(signal_chan, syscall.SIGINT)
  running := true
  var lastErr error
  for running {
    select {
    case err := <-error_chan:
      log.Printf("error: error has occurred: %s", err.Error())
      lastErr = err
      if e := stopService(); e != nil {
        lastErr = e
      }
      running = false
    case sig := <-signal_chan:
      log.Printf("info: signal trap(%s)", sig.String())
      if err := stopService(); err != nil {
        lastErr = err
      }
      running = false
    }
  }
  if lastErr == nil {
    log.Printf("info: shutdown successful")
    return nil
  }
  return lastErr
}

func main(){
  colog.SetDefaultLevel(colog.LDebug)
  colog.SetMinLevel(colog.LInfo)

  colog.SetFormatter(&colog.StdFormatter{
    Flag: log.Ldate | log.Ltime | log.Lshortfile,
  })
  colog.Register()

  app         := cli.NewApp()
  app.Version  = wsmsg.Version
  app.Name     = wsmsg.AppName
  app.Author   = ""
  app.Email    = ""
  app.Usage    = ""
  app.Action   = action
  app.Commands = Commands
  app.Flags    = []cli.Flag{
    cli.StringFlag{
      Name: "i, ip",
      Usage: "server bind-ip",
      Value: wsmsg.DEFAULT_BIND_IP,
      EnvVar: "WSMSG_BIND_IP",
    },
    cli.IntFlag{
      Name: "p, port",
      Usage: "server bind-port",
      Value: wsmsg.DEFAULT_BIND_PORT,
      EnvVar: "WSMSG_BIND_PORT",
    },
    cli.IntFlag{
      Name: "max-payload",
      Usage: "msg max payload size",
      Value: wsmsg.DEFAULT_MSG_MAX_PAYLOAD,
      EnvVar: "WSMSG_MAX_PAYLOAD",
    },
    cli.StringFlag{
      Name: "log-dir",
      Usage: "/path/to/log directory",
      Value: wsmsg.DEFAULT_LOG_DIR,
      EnvVar: "WSMSG_LOG_DIR",
    },
    cli.IntFlag{
      Name: "procs, P",
      Usage: "attach cpu(s)",
      Value: runtime.NumCPU(),
      EnvVar: "WSMSG_PROCS",
    },
    cli.BoolFlag{
      Name: "debug, d",
      Usage: "debug mode",
      EnvVar: "WSMSG_DEBUG",
    },
    cli.BoolFlag{
      Name: "verbose, V",
      Usage: "verbose. more message",
      EnvVar: "WSMSG_VERBOSE",
    },
  }
  if err := app.Run(os.Args); err != nil {
    log.Printf("error: %s", err.Error())
    cli.OsExiter(1)
  }
}
