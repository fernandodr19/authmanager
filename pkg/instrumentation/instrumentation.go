package instrumentation

import (
	"errors"

	"github.com/sirupsen/logrus"
)

type Instrumentation struct {
	Logger *logrus.Entry
}

var instrumentation *Instrumentation

func init() {
	instrumentation = &Instrumentation{
		Logger: logrus.NewEntry(logrus.New()),
	}
}

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
