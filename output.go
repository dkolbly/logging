package logging

import (
	"fmt"
	"io"
	"os"
)

type Writer interface {
	Write(*Record, int)
}

type TextWriter struct {
	dest    io.Writer
	format  Formatter
	NoColor bool
}

func MustTextWriter(dest io.Writer, format string) *TextWriter {
	tw, err := NewTextWriter(dest, format)
	if err != nil {
		panic(err)
	}
	return tw
}

func NewTextWriterUsing(dest io.Writer, f Formatter) *TextWriter {
	return &TextWriter{
		dest:   dest,
		format: f,
	}
}

func NewTextWriter(dest io.Writer, format string) (*TextWriter, error) {
	f, err := PatternFormatter(format)
	if err != nil {
		return nil, err
	}
	return &TextWriter{
		dest:   dest,
		format: f,
	}, nil
}

func (t *TextWriter) Write(rec *Record, skip int) {
	t.dest.Write(t.format.Format(rec, t.NoColor, skip+1))
}

type Stdout struct{}

func (std Stdout) Write(rec *Record, skip int) {
	buf := fmt.Sprintf(rec.Format, rec.Args...)
	os.Stdout.Write([]byte(buf))
	os.Stdout.Write([]byte{'\n'})
}
