package clicommon

import (
	"github.com/comail/colog"

	"github.com/octu0/nats-wsmsg/logger"
)

func InitLogger(config Config) error {
	if config.DebugMode {
		colog.SetMinLevel(colog.LDebug)
		if config.VerboseMode {
			colog.SetMinLevel(colog.LTrace)
		}
	}
	err := logger.Init(
		logger.LogDir(config.LogDir),
		logger.DebugMode(config.DebugMode),
		logger.VerboseMode(config.VerboseMode),
	)
	if err != nil {
		return err
	}

	colog.SetOutput(logger.GetLogger("general"))
	return nil
}
