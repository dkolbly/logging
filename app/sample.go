package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dkolbly/logging"
)

var log = logging.New("yo")

func main() {
	log.Info("Hello, world (%d)", 10)

	bg := context.Background()

	foo(log.Re(proc("foo")).In(bg))
	bar(log.To(logging.SelfDebug{}).Re(proc("bob")).In(bg))
	baz(log.To(&JSONWriter{}).Re(proc("alice")).In(bg), false)
	baz(log.To(&JSONWriter{}).Re(proc("alice")).In(bg), true)
	blech(log.Re(proc("carol")).In(bg), false)
	blech(log.Re(proc("carol")).In(bg), true)
	blech(log.Re(proc("carol")).In(bg), false)

	f1 := logging.MustPatternFormatter("aaa--- %{message} NO TAGS\n")
	f2 := logging.MustPatternFormatter("bbb--- %{message} TAGS: PROC=%{annot/proc}\n")
	f3 := logging.MustPatternFormatter("ccc--- %{leftmargin}%{message}\nTAGS: %{annot/slub:-}%{color:=#r}SLUB%{/color}\n")
		
	mflog := log.To(logging.NewTextWriterUsing(
		os.Stdout,
		logging.NewMultiFormat(f1, f2, f3)))
	mflog.Info("Hello")
	mflog.Re(proc("alice")).Info("Hi!")
	mflog.Re(proc("alice")).Re(slub("xxx")).Info("Hi!")
	mflog.Re(slub("Xxx")).Re(proc("Alice")).Info("Hiyo!")
	mflog.Re(slub("Xxx")).Info("slub only...")
}

type proc string

func (p proc) Annotate(rec *logging.Record) {
	rec.Annotate("proc", string(p))
}

type slub string

func (p slub) Annotate(rec *logging.Record) {
	rec.Annotate("slub", string(p))
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

func baz(ctx context.Context, nocolor bool) {
	log := logging.In(ctx)

	wr, err := logging.NewTextWriter(os.Stdout, "%{color}%{time:15:04:05.000} %{level:-8s} [%{module}|%{shortfile:%.5s:%06d}]%{/color} %{message}\n")
	if err != nil {
		panic(err)
	}
	if nocolor {
		wr.NoColor = true
	}
	log.Tee(wr).Info("This is my message: %d", 5)

	log.Info("A%d B=%q", 37, "fun")
}

func blech(ctx context.Context, nocolor bool) {
	log := logging.In(ctx)

	wr, err := logging.NewTextWriter(os.Stdout, "%{color:=b}[%{id:05d}]%{/color} %{color}%{time:15:04:05.000} %{level:-8s} [%{module}|%{shortfile:%.5s:%06d}]%{/color} %{message}\n")
	if err != nil {
		panic(err)
	}
	wr.NoColor = nocolor
	blech1(log.To(wr).In(ctx))
}

func blech1(ctx context.Context) {
	log := logging.In(ctx)
	for i := 0; i<3; i++ {
		log.Info("Hello (%d)", i)
	}
}


type JSONWriter struct{}

func (j *JSONWriter) Write(rec *logging.Record, skip int) {
	buf, err := rec.JSON(skip+1)
	if err == nil {
		fmt.Printf("%s\n", buf)
	}
}
