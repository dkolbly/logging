package logging

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

type Record struct {
	ID          uint64                 `json:"id"`
	Module      string                 `json:"module"`
	Level       Level                  `json:"level"`
	Timestamp   time.Time              `json:"timestamp"`
	Format      string                 `json:"format"`
	Args        []interface{}          `json:"args"`
	Annotations map[string]interface{} `json:"annotations,omitempty"`
}

type formattedRecord struct {
	Record
	/*	Id          uint64                 `json:"id"`
		Module      string                 `json:"module"`
		Level       Level                  `json:"level"`
		Timestamp   time.Time              `json:"timestamp"`
		Format      string                 `json:"format"`
		Args        []interface{}          `json:"args"`
		Annotations map[string]interface{} `json:"annotations,omitempty"`
	*/
	Message string `json:"message"`
	File    string `json:"file"`
	Line    int    `json:"line"`
}

func (r *Record) JSON(skip int) ([]byte, error) {
	fr := &formattedRecord{
		Record: *r,
		/*
			Id:          r.Id,
			Module:      r.Module,
			Level:       r.Level,
			Timestamp:   r.Timestamp,
			Format:      r.Format,
			Args:        r.Args,
			Annotations: r.Annotations,
		*/
		Message: fmt.Sprintf(r.Format, r.Args...),
	}
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		fr.File = file
		fr.Line = line
	}
	return json.Marshal(fr)
}

// Level is a small integer (0-7) that identifies the severity or
// importance of the information contained in the message.  Taken
// from the definitions for syslog, the levels are:
//
// 0 EMERGENCY  The system is unusable
// 1 ALERT      Should be corrected immediately
// 2 CRITICAL   Critical conditions; a failure of the main application
// 3 ERROR      Error conditions; something is actually broken
// 4 WARNING    Something might be broken, or might break soon
// 5 NOTICE     Something unusual, but probably not an error
// 6 INFO       Normal operation; no action required
// 7 DEBUG      Useful to developers for debugging the app
//
type Level uint8

const (
	EMERGENCY = Level(iota)
	ALERT
	CRITICAL
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

var levelNames = [8]string{
	EMERGENCY: "EMERGENCY",
	ALERT:     "ALERT",
	CRITICAL:  "CRITICAL",
	ERROR:     "ERROR",
	WARNING:   "WARNING",
	NOTICE:    "NOTICE",
	INFO:      "INFO",
	DEBUG:     "DEBUG",
}

// String returns the string representation of a logging level.
func (p Level) String() string {
	return levelNames[p]
}
