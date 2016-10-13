package main

import (
	"os"
	"context"
	"github.com/dkolbly/logging"
)

var log = logging.New("yo")

func main() {
	log.Info("Hello, world (%d)", 10)

	bg := context.Background()

	foo(log.Re(proc("foo")).In(bg))
	bar(log.To(logging.SelfDebug{}).Re(proc("bob")).In(bg))
}

type proc string

func (p proc) Annotate(rec *logging.Record) {
	rec.Annotate("proc", string(p))
}

func foo(ctx context.Context) {
	log := logging.In(ctx)
	log.Info("Hi")
}

func bar(ctx context.Context) {
	log := logging.In(ctx)
	log.Info("Bar (%s)", "cat")

	wr, err := logging.NewTextWriter(os.Stdout, "Yo (%{message:-10s})\n")
	if err != nil {
		panic(err)
	}
	log.Tee(wr).Info("---")
}
