package logging

import (
	"context"
	"sync/atomic"
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
		module:  l.module,
		outputs: wr,
	}
}

func (l *Logger) Tee(wr Writer) *Logger {
	o := make([]Writer, len(l.outputs)+1)
	copy(o[1:], l.outputs)
	o[0] = wr
	return &Logger{
		module:  l.module,
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

var seq uint64 = 1

func (l *Logger) dispatch(format string, args []interface{}, v Level, s int) {
	rec := &Record{
		ID:          atomic.AddUint64(&seq, 1),
		Module:      l.module,
		Annotations: l.annot,
		Level:       v,
		Timestamp:   time.Now(),
		Format:      format,
		Args:        args,
	}
	for _, wr := range l.outputs {
		wr.Write(rec, s+1)
	}
}

type logger int

const CurrentLoggerKey = logger(0)

func (l *Logger) In(ctx context.Context) context.Context {
	return context.WithValue(ctx, CurrentLoggerKey, l)
}

func In(ctx context.Context) *Logger {
	return ctx.Value(CurrentLoggerKey).(*Logger)
}

func (l *Logger) Info(format string, args ...interface{}) {
	// this 2 accounts for our depth and theirs
	l.dispatch(format, args, INFO, 2)
}

