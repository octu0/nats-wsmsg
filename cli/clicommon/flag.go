package clicommon

import (
	"gopkg.in/urfave/cli.v1"
)

func GlobalFlag() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "log-dir",
			Usage:  "/path/to/log directory",
			Value:  "/tmp",
			EnvVar: "WSMSG_LOG_DIR",
		},
		cli.BoolFlag{
			Name:   "debug, d",
			Usage:  "debug mode",
			EnvVar: "WSMSG_DEBUG",
		},
		cli.BoolFlag{
			Name:   "verbose, V",
			Usage:  "verbose. more message",
			EnvVar: "WSMSG_VERBOSE",
		},
	}
}
