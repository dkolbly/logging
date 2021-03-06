// Package pretty provides a simple method to reconfigure the default backend
// for pretty, colored log messages that respects whether the output is a
// tty or not (turning off coloring in that case)
package pretty

import (
	"os"
	
	"github.com/dkolbly/logging"
	"github.com/mattn/go-isatty"
)

var TextWriter = logging.MustTextWriter(os.Stdout, DefaultTextFormat)
var Writer = logging.MustFilter(TextWriter)

func init() {
	tty := isatty.IsTerminal(os.Stdout.Fd())
	Writer.SetLevel(logging.INFO, "*")
	TextWriter.NoColor = !tty
	logging.DefaultBackend.Target = Writer
}

const DefaultTextFormat = "%{color}%{time:15:04:05.000} %{level:-8s} [%{module}|%{shortfile:%s:%d}]%{/color} %{leftmargin}%{message}\n"


// ForceColor causes coloring to be enabled for the writer owned by
// this package.  Useful for implementing a `--color` flag
func ForceColor() {
	TextWriter.NoColor = false
}

// ForceColor causes coloring to be disabled for the writer owned by
// this package.  Useful for implementing a `--no-color` flag
func ForceNoColor() {
	TextWriter.NoColor = true
}

func Debug() {
	Writer.SetLevel(logging.DEBUG, "*")
}

