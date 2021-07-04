package logger

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

type ctxKey string

const LoggerCtxKey ctxKey = "logget-ctx-key"

var defaultLogger *logrus.Entry

func init() {
	defaultLogger = logrus.NewEntry(logrus.New())
}

func Register(logger *logrus.Entry) error {
	if logger == nil {
		return errors.New("tried to register nil logger")
	}

	defaultLogger = logger
	return nil
}

// Default returns the default logger
func Default() *logrus.Entry {
	return defaultLogger
}

// FromCtx retrives logger from context (inserted by api middleware)
func FromCtx(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(LoggerCtxKey)

	if logger == nil {
		// Fallback to deafault if nil
		return Default()
	}

	return logger.(*logrus.Entry)
}
