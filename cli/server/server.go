package server

import (
	"gopkg.in/urfave/cli.v1"
)

var commands []cli.Command

func addCommand(cmd cli.Command) {
	commands = append(commands, cmd)
}

func Command() []cli.Command {
	return commands
}
