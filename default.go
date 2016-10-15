package logging

// A MutableWriter is the one log writer that is designed to be mutable,
// and is primarily used as the default log destination so that modules
// can configure their logging using:
//
//    var log = logging.MustGetLogger("me")
//
// and the main application can configure where those log messages go
// to by resetting the Target of the DefaultBackend
type MutableWriter struct {
	Target Writer
}

func (m *MutableWriter) Write(rec *Record, skip int) {
	m.Target.Write(rec, skip+1)
}

var DefaultBackend = &MutableWriter{Stdout{}}
