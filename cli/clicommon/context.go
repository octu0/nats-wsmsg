package clicommon

import (
	"context"

	"github.com/octu0/nats-wsmsg/logger"
)

type withFunc func(ctx context.Context) context.Context

func cliContext(funcs ...withFunc) context.Context {
	ctx := context.Background()
	for _, fn := range funcs {
		ctx = fn(ctx)
	}
	return ctx
}

func withConfig(config Config) withFunc {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, "config", config)
	}
}

func ValueConfig(ctx context.Context) Config {
	return ctx.Value("config").(Config)
}

func withGeneralLogger(lg *logger.GeneralLogger) withFunc {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, "logger.general", lg)
	}
}

func ValueGeneralLogger(ctx context.Context) *logger.GeneralLogger {
	return ctx.Value("logger.general").(*logger.GeneralLogger)
}

func withHttpLogger(lg *logger.HttpLogger) withFunc {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, "logger.http", lg)
	}
}

func ValueHttpLogger(ctx context.Context) *logger.HttpLogger {
	return ctx.Value("logger.http").(*logger.HttpLogger)
}

func withNatsLogger(lg *logger.NatsLogger) withFunc {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, "logger.nats", lg)
	}
}

func ValueNatsLogger(ctx context.Context) *logger.NatsLogger {
	return ctx.Value("logger.nats").(*logger.NatsLogger)
}
