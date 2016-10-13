package logging

import (
	"time"
)

type Record struct {
	Id uint64
	Module      string
	Level       Level
	Annotations map[string]interface{}
	Timestamp   time.Time
	Format      string
	Args        []interface{}
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
