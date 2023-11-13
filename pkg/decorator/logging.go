package decorator

import (
	"context"
	"fmt"

	"github.com/htquangg/microservices-poc/pkg/logger"
)

type commandLoggingDecorator[C any] struct {
	base CommandHandler[C]
	log  logger.Logger
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	handlerType := generateActionName(cmd)

	fields := logger.Fields{
		"command":      handlerType,
		"command_body": fmt.Sprintf("%#v", cmd),
	}

	d.log.Debugw("Executing command", fields)
	defer func() {
		if err == nil {
			d.log.Debugw("Command executed successfully", fields)
		} else {
			d.log.Err("Failed to execute command", err)
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type queryLoggingDecorator[C any, R any] struct {
	base QueryHandler[C, R]
	log  logger.Logger
}

func (d queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	fields := logger.Fields{
		"query":      generateActionName(cmd),
		"query_body": fmt.Sprintf("%#v", cmd),
	}

	d.log.Debugw("Executing query", fields)
	defer func() {
		if err == nil {
			d.log.Info("Query executed successfully")
		} else {
			d.log.Err("Failed to execute query", err)
		}
	}()

	return d.base.Handle(ctx, cmd)
}
