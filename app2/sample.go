package main

import (
	"github.com/dkolbly/logging"
	_ "github.com/dkolbly/logging/pretty"
)

var log = logging.New("samp")

func main() {
	log.Info("Version 1")
	log.Info("Here is some info you'd like:\n1. Red\n2. Blue")
	log.Info("Here is some more info you'd like:\n1. Cats")
}
