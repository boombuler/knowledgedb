package log

import (
	"fmt"
)

type Logger interface {
	Fatal(v ...interface{})
	Error(v ...interface{})
	Warning(v ...interface{})
	Info(v ...interface{})
}

var current Logger

func Fatal(v ...interface{}) {
	if current != nil {
		current.Fatal(v...)
	}
}

func Error(v ...interface{}) {
	if current != nil {
		current.Fatal(v...)
	}
}
func Warning(v ...interface{}) {
	if current != nil {
		current.Warning(v...)
	}
}

func Info(v ...interface{}) {
	if current != nil {
		current.Info(v...)
	}
}

func Fatalf(f string, v ...interface{}) {
	Fatal(fmt.Sprintf(f, v...))
}

func Errorf(format string, a ...interface{}) {
	Error(fmt.Sprintf(format, a...))
}

func Warningf(format string, a ...interface{}) {
	Warning(fmt.Sprintf(format, a...))
}

func Infof(format string, a ...interface{}) {
	Info(fmt.Sprintf(format, a...))
}
