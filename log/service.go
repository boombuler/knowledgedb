package log

import (
	"github.com/kardianos/service"
	"os"
)

type serviceLogger struct {
	service.Logger
}

func SetServiceLogger(l service.Logger) {
	if l != nil {
		current = serviceLogger{l}
	} else {
		current = nil
	}
}

func init() {
	SetServiceLogger(service.ConsoleLogger)
}

func (l serviceLogger) Fatal(v ...interface{}) {
	l.Logger.Error(v...)
	os.Exit(1)
}
func (l serviceLogger) Error(v ...interface{}) {
	l.Logger.Error(v...)
}
func (l serviceLogger) Warning(v ...interface{}) {
	l.Logger.Warning(v...)
}
func (l serviceLogger) Info(v ...interface{}) {
	l.Logger.Info(v...)
}
