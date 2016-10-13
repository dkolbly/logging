package logging

import (
	"io"
	"fmt"
	"os"
)

type Writer interface {
	Write(*Record)
}

type TextWriter struct {
	dest   io.Writer
	format Formatter
}

func NewTextWriter(dest io.Writer, format string) (Writer, error) {
	f, err := PatternFormatter(format)
	if err != nil {
		return nil, err
	}
	return &TextWriter{
		dest:   dest,
		format: f,
	}, nil
}

func (t *TextWriter) Write(rec *Record) {
	t.dest.Write(t.format.Format(rec))
}

type Stdout struct{}

func (std Stdout) Write(rec *Record) {
	buf := fmt.Sprintf(rec.Format, rec.Args...)
	os.Stdout.Write([]byte(buf))
	os.Stdout.Write([]byte{'\n'})
}
