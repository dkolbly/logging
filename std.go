package logging

import (
	"bufio"
	"io"
	"log"
)

// StdLogger returns a true *log.Logger which feeds the output
// to our interface
func (l *Logger) StdLogger() *log.Logger {
	rd, wr := io.Pipe()
	go l.slurp(rd)
	return log.New(wr, "", 0)
}

func (l *Logger) slurp(rd io.Reader) {
	buf := bufio.NewReader(rd)
	for {
		line, err := buf.ReadBytes('\n')
		if err != nil {
			return
		}
		l.dispatch("%s", []interface{}{string(line[:len(line)-1])}, NOTICE, baseDepth)
	}
}
