package logging

import (
	"path"
	"sync"
)

type levelRule struct {
	pattern string
	level   Level
}

type LevelFilter struct {
	rules  []levelRule
	target Writer
	// cache:
	lock      sync.Mutex
	threshold map[string]Level
}

func MustFilter(w Writer) *LevelFilter {
	return &LevelFilter{
		target: w,
	}
}

func (f *LevelFilter) Write(rec *Record, skip int) {
	if rec.Level <= f.GetLevel(rec.Module) {
		f.target.Write(rec, skip+1)
	}
}

// GetLevel returns the log level for the given module.
func (f *LevelFilter) GetLevel(module string) Level {
	f.lock.Lock()
	defer f.lock.Unlock()
	if level, ok := f.threshold[module]; ok {
		return level
	}
	level := DEBUG // default value in case of no match
	for _, r := range f.rules {
		match, err := path.Match(r.pattern, module)
		if err != nil {
			// fall back to exact matching
			match = r.pattern == module
		}
		if match {
			level = r.level
			break
		}
	}
	// cache the result for later
	f.threshold[module] = level
	return level
}

// SetLevel sets the log level for the given module.  If the
// module contains one of the special characters '*' or '?',
// then it is interpreted as a glob (unless `path.Match` rejects
// the pattern as being malformed, in which case an error is logged
// and the module is considered an exact match)
func (f *LevelFilter) SetLevel(level Level, module string) {
	// add the rule
	f.lock.Lock()
	defer f.lock.Unlock()

	// clear the cache
	f.threshold = make(map[string]Level)

	// see if there's an exact match, in which case overwrite the levelRule
	for i, rule := range f.rules {
		if rule.pattern == module {
			f.rules[i].level = level
			return
		}
	}
	f.rules = append(f.rules, levelRule{module, level})

}
