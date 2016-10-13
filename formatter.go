package logging

import (
	"fmt"
	"regexp"
	"bytes"
)

type Formatter interface {
	Format(*Record) []byte
}

type fragmentFormatter func(*bytes.Buffer, *Record)

func literal(lit string) fragmentFormatter {
	return func(dst *bytes.Buffer, src *Record) {
		dst.WriteString(lit)
	}
}

type LegacyPatternFormatter struct {
	fragments []fragmentFormatter
}

func (pf *LegacyPatternFormatter) Format(r *Record) []byte {
	buf := &bytes.Buffer{}
	for _, frag := range pf.fragments {
		frag(buf, r)
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
	frags, err := compilePattern(pat)
	if err != nil {
		return nil, err
	}
	return &LegacyPatternFormatter{
		fragments: frags,
	}, nil
}

var formatRe *regexp.Regexp = regexp.MustCompile(`%{([a-z]+)(?::(.*?[^\\]))?}`)

func compilePattern(pat string) ([]fragmentFormatter, error) {
	// Find all the %{...} pieces
	matches := formatRe.FindAllStringSubmatchIndex(pat, -1)
	if matches == nil {
		return nil, fmt.Errorf("logger: invalid log format: %q", pat)
	}

	var frags []fragmentFormatter

	prev := 0
	for _, m := range matches {
		start, end := m[0], m[1]
		if start > prev {
			frags = append(frags, literal(pat[prev:start]))
		}
		verb := pat[m[2]:m[3]]
		layout := ""
		if m[4] != -1 {
			layout = pat[m[4]:m[5]]
		}
		fragMaker, ok := verbTable[verb]
		if !ok {
			return nil, fmt.Errorf("logger: unknown verb %q in %q", verb, pat)
		}
		frag, err := fragMaker(layout)
		if err != nil {
			return nil, err
		}
		frags = append(frags, frag)
		prev = end
	}
	if prev < len(pat) {
		frags = append(frags, literal(pat[prev:]))
	}
	return frags, nil
}

const rfc3339Milli = "2006-01-02T15:04:05.999Z07:00"

type fragMaker func(string) (fragmentFormatter, error)

var verbTable = map[string]fragMaker{
	"time": makeTimeFrag,
	"message": makeMessageFrag,
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
func makeTimeFrag(options string) (fragmentFormatter, error) {
	return func(dest *bytes.Buffer, r *Record) {
	}, nil
	
}

func makeMessageFrag(options string) (fragmentFormatter, error) {
	if options == "" {
		options = "%s"
	} else {
		options = "%" + options
	}
	return func(dest *bytes.Buffer, r *Record) {
		str := fmt.Sprintf(r.Format, r.Args...)
		fmt.Fprintf(dest, options, str)
	}, nil
}
