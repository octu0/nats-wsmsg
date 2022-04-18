package clicommon

import (
	"context"

	"gopkg.in/urfave/cli.v1"

	"github.com/octu0/nats-wsmsg/logger"
)

func Prepare(c *cli.Context) (context.Context, error) {
	config := CreateConfig(c)
	if err := InitLogger(config); err != nil {
		return nil, err
	}

	return cliContext(
		withConfig(config),
		withGeneralLogger(logger.GetLogger("general").(*logger.GeneralLogger)),
		withHttpLogger(logger.GetLogger("http").(*logger.HttpLogger)),
		withNatsLogger(logger.GetLogger("nats").(*logger.NatsLogger)),
	), nil
}
