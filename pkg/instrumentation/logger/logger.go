package logger

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

type ctxKey string

const loggerCtxKey ctxKey = "logger-ctx-key"

var defaultLogger *logrus.Entry

func init() {
	defaultLogger = logrus.NewEntry(logrus.New())
}

// Register overrides the default logger
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
	logger := ctx.Value(loggerCtxKey)

	if logger == nil {
		// Fallback to deafault if nil
		return Default()
	}

	return logger.(*logrus.Entry)
}

// ToCtx returns a new context with the provided logger
func ToCtx(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}
