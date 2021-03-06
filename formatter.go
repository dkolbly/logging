package logging

import (
	"bytes"
	"fmt"
	"path"
	"regexp"
	"runtime"
	"strings"
)

type Formatter interface {
	Format(*Record, bool, int) []byte
}

type outputContext struct {
	dst        bytes.Buffer
	src        *Record
	stackSkip  int
	leftMargin int
	column     int
	bol        bool
}

func (ctx *outputContext) WriteString(s string) (int, error) {
	return ctx.Write([]byte(s))
}

func (ctx *outputContext) Write(data []byte) (int, error) {
	for _, b := range data {
		if b == '\n' {
			ctx.dst.WriteByte(b)
			ctx.bol = true
		} else {
			if ctx.bol {
				for j := 0; j < ctx.leftMargin; j++ {
					ctx.dst.WriteByte(' ')
				}
				ctx.column = ctx.leftMargin
				ctx.bol = false
			}
			ctx.dst.WriteByte(b)
			ctx.column++
		}
	}
	return len(data), nil
}

type fragmentFormatter func(*outputContext)

// TODO we can unfold the Write() loop by pre-scanning the literal data
// and splitting it into lines (and special case the 1-line case, and maybe
// even the 1 byte case)
func literals(lit []byte) fragmentFormatter {
	return func(ctx *outputContext) {
		ctx.Write(lit)
	}
}

func literal(lit string) fragmentFormatter {
	return func(ctx *outputContext) {
		ctx.Write([]byte(lit))
	}
}

type LegacyPatternFormatter struct {
	fragments        []fragmentFormatter
	nocolorFragments []fragmentFormatter
	requireAnnots    []string
}

func (pf *LegacyPatternFormatter) Match(rec *Record) bool {
	// check to see if all the annotations we want are here
	for _, a := range pf.requireAnnots {
		if _, ok := rec.Annotations[a]; !ok {
			// oops, missing one!
			return false
		}
	}
	return true
}

func (pf *LegacyPatternFormatter) Format(r *Record, nocolor bool, skip int) []byte {
	ctx := &outputContext{
		src:       r,
		stackSkip: skip + 1,
	}
	frags := pf.fragments
	if nocolor {
		frags = pf.nocolorFragments
	}
	for _, frag := range frags {
		frag(ctx)
	}
	return ctx.dst.Bytes()
}

func MustPatternFormatter(pat string) *LegacyPatternFormatter {
	f, err := PatternFormatter(pat)
	if err != nil {
		panic(err)
	}
	return f
}

func PatternFormatter(pat string) (*LegacyPatternFormatter, error) {
	frags, ncfrags, use, err := compilePattern(pat)
	if err != nil {
		return nil, err
	}
	return &LegacyPatternFormatter{
		fragments:        frags,
		nocolorFragments: ncfrags,
		requireAnnots:    use,
	}, nil
}

var formatRe *regexp.Regexp = regexp.MustCompile(`%{([a-z/]+)(?::(.*?[^\\]))?}`)

func compilePattern(pat string) ([]fragmentFormatter, []fragmentFormatter, []string, error) {
	// Find all the %{...} pieces
	matches := formatRe.FindAllStringSubmatchIndex(pat, -1)
	if matches == nil {
		return nil, nil, nil, fmt.Errorf("logger: invalid log format: %q", pat)
	}

	var frags []fragmentFormatter
	var nocolorFrags []fragmentFormatter

	push := func(ff fragmentFormatter, iscolor bool) {
		frags = append(frags, ff)
		if !iscolor {
			nocolorFrags = append(nocolorFrags, ff)
		}
	}

	prev := 0
	require := []string{}
	for _, m := range matches {
		start, end := m[0], m[1]
		if start > prev {
			push(literal(pat[prev:start]), false)
		}
		verb := pat[m[2]:m[3]]
		layout := ""
		if m[4] != -1 {
			layout = pat[m[4]:m[5]]
		}
		var frag fragmentFormatter
		var err error

		if strings.HasPrefix(verb, "annot/") {
			annot := verb[6:]
			frag = makeAnnotFrag(annot, layout)
			require = append(require, annot)
		} else {
			fragMaker, ok := verbTable[verb]
			if !ok {
				return nil, nil, nil, fmt.Errorf("logger: unknown verb %q in %q", verb, pat)
			}
			frag, err = fragMaker(layout)
			if err != nil {
				return nil, nil, nil, err
			}
		}
		push(frag, colorTable[verb])
		prev = end
	}
	if prev < len(pat) {
		push(literal(pat[prev:]), false)
	}
	return frags, nocolorFrags, require, nil
}

type fragMaker func(string) (fragmentFormatter, error)

var colorTable = map[string]bool{
	"color":  true,
	"/color": true,
}

var verbTable = map[string]fragMaker{
	"time":       makeTimeFrag,
	"message":    makeMessageFrag,
	"color":      makeColorFrag,
	"/color":     makeColorResetFrag,
	"module":     makeModuleFrag,
	"shortfile":  makeShortFileFrag,
	"level":      makeLevelFrag,
	"id":         makeIdFrag,
	"leftmargin": makeLeftMarginFrag,
}

func stringopt(options string) string {
	if options == "" {
		return "%s"
	} else {
		return "%" + options
	}
}

func makeLeftMarginFrag(_ string) (fragmentFormatter, error) {
	return func(ctx *outputContext) {
		ctx.leftMargin = ctx.column
	}, nil
}

func makeModuleFrag(options string) (fragmentFormatter, error) {
	options = stringopt(options)
	return func(ctx *outputContext) {
		fmt.Fprintf(ctx, options, ctx.src.Module)
	}, nil
}

func makeIdFrag(options string) (fragmentFormatter, error) {
	if options == "" {
		options = "%d"
	} else {
		options = "%" + options
	}

	return func(ctx *outputContext) {
		fmt.Fprintf(ctx, options, ctx.src.ID)
	}, nil
}

type Sourcer interface {
	File() string
	Line() int
}

func makeShortFileFrag(options string) (fragmentFormatter, error) {
	if options == "" {
		options = "%[1]s:%[2]d"
	}
	return func(ctx *outputContext) {
		if src, ok := ctx.src.Annotations["source"]; ok {
			if sourcer, ok := src.(Sourcer); ok {
				fmt.Fprintf(ctx, options, path.Base(sourcer.File()), sourcer.Line())
				return
			}
		}

		_, file, line, ok := runtime.Caller(ctx.stackSkip)
		if !ok {
			file = "???"
			line = 0
		} else {
			file = path.Base(file)
		}
		fmt.Fprintf(ctx, options, file, line)
	}, nil
}

func makeLevelFrag(options string) (fragmentFormatter, error) {
	options = stringopt(options)
	return func(ctx *outputContext) {
		fmt.Fprintf(ctx, options, levelNames[ctx.src.Level])
	}, nil
}

/*
func makeFrag(verb, extra string) (fragmentFormatter, error) {
	switch verb {
	case "time":
		return makeTimeFrag(extra)
	case "level":
		return makeLevelFrag(extra)
	case "id":
		return makeIdFrag(extra)
	case "pid":
	case "program":
	case "module":
	case "message":
	case "longfile":
	case "shortfile":
	case "color":
*/

const rfc3339Milli = "2006-01-02T15:04:05.999Z07:00"

func makeTimeFrag(options string) (fragmentFormatter, error) {
	if options == "" {
		options = rfc3339Milli
	}
	return func(ctx *outputContext) {
		ctx.WriteString(ctx.src.Timestamp.Format(options))
	}, nil

}

func makeMessageFrag(options string) (fragmentFormatter, error) {
	if options == "" {
		options = "%s"
	} else {
		options = "%" + options
	}
	return func(ctx *outputContext) {
		str := fmt.Sprintf(ctx.src.Format, ctx.src.Args...)
		fmt.Fprintf(ctx, options, str)
	}, nil
}

func makeAnnotFrag(annot, options string) fragmentFormatter {
	// allow the output to be completely suppressed with %{annot/foo:-},
	// in which case the tag is only used for format message selection.
	// it's basically equivalent to using ".0s" but more efficient and
	// easier to understand
	if options == "-" {
		return func(ctx *outputContext) {}
	}

	if options == "" {
		options = "%s"
	} else {
		options = "%" + options
	}

	return func(ctx *outputContext) {
		if value, ok := ctx.src.Annotations[annot]; ok {
			fmt.Fprintf(ctx, options, value)
		}
	}
}
