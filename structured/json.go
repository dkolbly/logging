package structured

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/dkolbly/logging"
	"github.com/mattn/go-isatty"
)

type StructuredLogRecord struct {
	Timestamp time.Time `json:"@timestamp"`
	Message   string    `json:"message"`
	Time      time.Time `json:"time"`
	Module    string    `json:"module"`
	File      string    `json:"file"`
	Line      int       `json:"line"`
	Seq       uint64    `json:"seq"`
	Level     string    `json:"level"`
}

type StructuredFormatter struct {
}

func Format(r *logging.Record, skip int) []byte {
	var lr StructuredLogRecord

	lr.Message = fmt.Sprintf(r.Format, r.Args...)
	lr.Seq = r.ID
	t := r.Timestamp.UTC()
	lr.Timestamp = t
	lr.Time = t
	lr.Module = r.Module
	lr.Level = r.Level.String()

	_, file, line, ok := runtime.Caller(skip)
	if ok {
		lr.File = path.Base(file)
		lr.Line = line
	}

	buf, err := json.Marshal(&lr)
	if err != nil {
		return nil
	}
	return buf
}

func (lf *StructuredFormatter) Format(r *logging.Record, nocolor bool, skip int) []byte {
	return append(Format(r, skip+1), '\n')
}

// usage:
//    logging.DefaultBackend.Target = structured.AutoWriter()

func NewWriter() logging.Writer {
	return logging.NewTextWriterUsing(os.Stdout, &StructuredFormatter{})
}

func NewPretty(color bool) logging.Writer {
	w := logging.MustTextWriter(os.Stdout, prettyFormat)
	w.NoColor = !color
	return w
}

func AutoWriter() *logging.LevelFilter {
	var color bool
	var json bool

	switch os.Getenv("LOGGING_FORMAT") {
	case "json":
		json = true
	case "color":
		color = true
		json = false
	case "no-color":
		color = false
		json = false
	default:
		color = isatty.IsTerminal(os.Stdout.Fd())
		json = !color
	}
	var writer logging.Writer
	if json {
		writer = NewWriter()
	} else {
		writer = NewPretty(color)
	}
	lf := logging.MustFilter(writer)
	lf.SetLevel(logging.INFO, "*")
	return lf
}

const prettyFormat = "%{color}%{time:15:04:05.000} %{level:-8s} [%{module}|%{shortfile:%s:%d}]%{/color} %{leftmargin}%{message}\n"
