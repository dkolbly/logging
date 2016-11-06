package main

import (
	stdlog "log"

	"github.com/dkolbly/logging"
	_ "github.com/dkolbly/logging/pretty"
)

var log = logging.New("samp")

func main() {
	log.Info("Version 1")
	log.Info("Here is some info you'd like:\n1. Red\n2. Blue")
	log.Info("Here is some more info you'd like:\n1. Cats")

	log.Debug("debug")
	log.Info("info")
	log.Notice("notice")
	log.Warning("warning")
	log.Error("error")
	log.Critical("critical")
	log.Alert("alert")
	log.Emergency("emergency")

	var x stdlog.Logger
	x = log
	x.Fatal("fatal")

	log.Warning("FATAL DID NOT EXIT?")
}
