package main

import (
	"log"
	"os"

	"github.com/comail/colog"
	"gopkg.in/urfave/cli.v1"

	"github.com/octu0/nats-wsmsg"
	"github.com/octu0/nats-wsmsg/cli/clicommon"
	"github.com/octu0/nats-wsmsg/cli/server"
)

func mergeCommand(values ...[]cli.Command) []cli.Command {
	merged := make([]cli.Command, 0, 0xff)
	for _, commands := range values {
		merged = append(merged, commands...)
	}
	return merged
}

func main() {
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LInfo)

	colog.SetFormatter(&colog.StdFormatter{
		Flag: log.Ldate | log.Ltime | log.Lshortfile,
	})
	colog.Register()

	app := cli.NewApp()
	app.Version = wsmsg.Version
	app.Name = wsmsg.AppName
	app.Author = ""
	app.Email = ""
	app.Usage = ""
	app.Commands = mergeCommand(
		server.Command(),
	)
	app.Flags = clicommon.GlobalFlag()
	if err := app.Run(os.Args); err != nil {
		log.Printf("error: %+v", err)
		cli.OsExiter(1)
	}
}
