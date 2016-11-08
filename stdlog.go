package logging

import (
	"bytes"
	"os"
)

// these implement log.Logger for *Logger

func (l *Logger) Fatal(args ...interface{}) {
	l.dispatch(synthesizeFormat(args), args, CRITICAL, baseDepth)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.dispatch(format, args, CRITICAL, baseDepth)
	os.Exit(1)
}

func synthesizeFormat(args []interface{}) string {
	if len(args) == 1 {
		return "%v"
	}
	var buf bytes.Buffer
	// pretend like it was a string so we don't add a space
	wasString := true
	for _, a := range args {
		_, isStr := a.(string)
		if !(wasString && isStr) {
			buf.WriteByte(' ')
		}
		buf.WriteByte('%')
		buf.WriteByte('v')
		wasString = isStr
	}
	return buf.String()
}
