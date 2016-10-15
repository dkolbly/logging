package logging

type MultiFormat struct {
	options  []MultiFormatCase
	fallback Formatter
}

func NewMultiFormat(base Formatter, f ...MultiFormatCase) Formatter {
	return &MultiFormat{
		options: f,
		fallback: base,
	}
}

type MultiFormatCase interface {
	Formatter
	Match(*Record) bool
}

func (mf *MultiFormat) Format(rec *Record, noColor bool, skip int) []byte {
	for _, opt := range mf.options {
		if opt.Match(rec) {
			return opt.Format(rec, noColor, skip+1)
		}
	}
	return mf.fallback.Format(rec, noColor, skip+1)
}
