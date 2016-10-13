package logging

import (
	"context"
	//"fmt"
	//"io"
	//"os"
	"time"
)

type Logger struct {
	annot   map[string]interface{}
	module  string
	outputs []Writer
}

func New(module string) *Logger {
	return &Logger{
		module:  module,
		outputs: []Writer{Stdout{}},
	}
}

func (l *Logger) To(wr ...Writer) *Logger {
	return &Logger{
		module: l.module,
		outputs: wr,
	}
}

func (l *Logger) Tee(wr Writer) *Logger {
	o := make([]Writer, len(l.outputs)+1)
	copy(o[1:], l.outputs)
	o[0] = wr
	return &Logger{
		module: l.module,
		outputs: o,
	}
}


/*
type Tee struct {
	first, remainder Output
}

func (t *Tee) Write(m *Record) {
	t.first.Write(m)
	t.remainder.Write(m)
}

type File struct {
	w io.Writer
}

func (m *Record) Formatted() string {
	if !m.isFormatted {
		m.formatted = fmt.Sprintf(m.Format, m.Args...)
		m.isFormatted = true
	}
	return m.formatted
}

func (f *File) Write(m *Record) {
	f.w.Write([]byte(m.Formatted()))
	f.w.Write([]byte{'\n'})
}

func New(ctx context.Context, module string) *Logger {
	return &Logger{
		Module: module,
		outputs: &outputChain{
			output: &File{w: os.Stdout},
			remainder: nil,
		},
	}
}

func (l *Logger) dispatch(rec *Record) {
	for a := l.annotators; a != nil; a = a.remainder {
		a.annotator.Annotate(rec)
	}
	for o := l.outputs; o != nil; o = o.remainder {
		o.output.Write(rec)
	}
}
*/

func (l *Logger) Info(format string, args ...interface{}) {
	msg := &Record{
		Module:      l.module,
		Annotations: l.annot,
		Level:       INFO,
		Timestamp:   time.Now(),
		Format:      format,
		Args:        args,
	}
	l.dispatch(msg)
}

func (l *Logger) dispatch(rec *Record) {
	for _, wr := range l.outputs {
		wr.Write(rec)
	}
}

type logger int
const CurrentLogger = logger(0)

func (l *Logger) In(ctx context.Context) context.Context {
	return context.WithValue(ctx, CurrentLogger, l)
}

func In(ctx context.Context) *Logger {
	return ctx.Value(CurrentLogger).(*Logger)
}
