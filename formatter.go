package logging

import (
	"bytes"
	"fmt"
	"path"
	"regexp"
	"runtime"
)

type Formatter interface {
	Format(*Record, bool, int) []byte
}

type fragmentFormatter func(*bytes.Buffer, *Record, int)

func literals(lit []byte) fragmentFormatter {
	return func(dst *bytes.Buffer, src *Record, _ int) {
		dst.Write(lit)
	}
}

func literal(lit string) fragmentFormatter {
	return func(dst *bytes.Buffer, src *Record, _ int) {
		dst.WriteString(lit)
	}
}

type LegacyPatternFormatter struct {
	fragments        []fragmentFormatter
	nocolorFragments []fragmentFormatter
}

func (pf *LegacyPatternFormatter) Format(r *Record, nocolor bool, skip int) []byte {
	buf := &bytes.Buffer{}
	frags := pf.fragments
	if nocolor {
		frags = pf.nocolorFragments
	}
	for _, frag := range frags {
		frag(buf, r, skip+1)
	}
	return buf.Bytes()
}

func MustPatternFormatter(pat string) Formatter {
	f, err := PatternFormatter(pat)
	if err != nil {
		panic(err)
	}
	return f
}

func PatternFormatter(pat string) (Formatter, error) {
	frags, ncfrags, err := compilePattern(pat)
	if err != nil {
		return nil, err
	}
	return &LegacyPatternFormatter{
		fragments:        frags,
		nocolorFragments: ncfrags,
	}, nil
}

var formatRe *regexp.Regexp = regexp.MustCompile(`%{([a-z/]+)(?::(.*?[^\\]))?}`)

func compilePattern(pat string) ([]fragmentFormatter, []fragmentFormatter, error) {
	// Find all the %{...} pieces
	matches := formatRe.FindAllStringSubmatchIndex(pat, -1)
	if matches == nil {
		return nil, nil, fmt.Errorf("logger: invalid log format: %q", pat)
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
		fragMaker, ok := verbTable[verb]
		if !ok {
			return nil, nil, fmt.Errorf("logger: unknown verb %q in %q", verb, pat)
		}
		frag, err := fragMaker(layout)
		if err != nil {
			return nil, nil, err
		}
		push(frag, colorTable[verb])
		prev = end
	}
	if prev < len(pat) {
		push(literal(pat[prev:]), false)
	}
	return frags, nocolorFrags, nil
}

type fragMaker func(string) (fragmentFormatter, error)

var colorTable = map[string]bool{
	"color":     true,
	"/color":    true,
}

var verbTable = map[string]fragMaker{
	"time":      makeTimeFrag,
	"message":   makeMessageFrag,
	"color":     makeColorFrag,
	"/color":    makeColorResetFrag,
	"module":    makeModuleFrag,
	"shortfile": makeShortFileFrag,
	"level":     makeLevelFrag,
	"id": makeIdFrag,
}

func stringopt(options string) string {
	if options == "" {
		return "%s"
	} else {
		return "%" + options
	}
}

func makeModuleFrag(options string) (fragmentFormatter, error) {
	options = stringopt(options)
	return func(dst *bytes.Buffer, src *Record, _ int) {
		fmt.Fprintf(dst, options, src.Module)
	}, nil
}

func makeIdFrag(options string) (fragmentFormatter, error) {
	if options == "" {
		options = "%d"
	} else {
		options = "%" + options
	}
	
	return func(dst *bytes.Buffer, src *Record, _ int) {
		fmt.Fprintf(dst, options, src.ID)
	}, nil
}

func makeShortFileFrag(options string) (fragmentFormatter, error) {
	if options == "" {
		options = "%[1]s:%[2]d"
	}
	return func(dst *bytes.Buffer, src *Record, depth int) {
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			file = "???"
			line = 0
		} else {
			file = path.Base(file)
		}
		fmt.Fprintf(dst, options, file, line)
	}, nil
}

func makeLevelFrag(options string) (fragmentFormatter, error) {
	options = stringopt(options)
	return func(dst *bytes.Buffer, src *Record, _ int) {
		fmt.Fprintf(dst, options, levelNames[src.Level])
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
	return func(dest *bytes.Buffer, r *Record, _ int) {
		dest.WriteString(r.Timestamp.Format(options))
	}, nil

}

func makeMessageFrag(options string) (fragmentFormatter, error) {
	if options == "" {
		options = "%s"
	} else {
		options = "%" + options
	}
	return func(dest *bytes.Buffer, r *Record, _ int) {
		str := fmt.Sprintf(r.Format, r.Args...)
		fmt.Fprintf(dest, options, str)
	}, nil
}
