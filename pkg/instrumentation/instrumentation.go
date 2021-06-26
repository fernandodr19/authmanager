package instrumentation

import (
	"errors"

	"github.com/sirupsen/logrus"
)

type Instrumentation struct {
	Logger *logrus.Entry
	// Tracer *newrelic.Application
}

var instrumentation *Instrumentation

func Register(instr *Instrumentation) error {
	if instr == nil {
		return errors.New("tried to register nil instrumentation")
	}

	instrumentation = instr
	return nil
}

func Logger() *logrus.Entry {
	return instrumentation.Logger
}

// func Tracer() *newrelic.Application {
// 	return instrumentation.Tracer
// }
