package logging

import (
	"os"
)

// these implement log.Logger for *Logger

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.dispatch(format, args, CRITICAL, baseDepth)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.dispatch(format, args, CRITICAL, baseDepth)
	os.Exit(1)
}

