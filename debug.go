package logging

import (
	"fmt"
	"runtime"
	"path"
)

// SelfDebug is an output used for debugging the logging system itself
type SelfDebug struct{}

func (std SelfDebug) Write(rec *Record, skip int) {
	fmt.Printf("LOG Time %s Level %s Module %q\n",
		rec.Timestamp.Format("2006-01-02 15:04:05"),
		rec.Level,
		rec.Module)
	fmt.Printf("    Format: %q\n", rec.Format)
	fmt.Printf("    Args: %#v\n", rec.Args)
	fmt.Printf("    Formatted: %q\n", fmt.Sprintf(rec.Format, rec.Args...))
	fmt.Printf("    %d Annotations:\n", len(rec.Annotations))
	for k, v := range rec.Annotations {
		fmt.Printf("      %q := %#v\n", k, v)
	}
	for i := -skip+2; i < 3; i++ {
		_, file, line, ok := runtime.Caller(skip+i)
		if ok {
			fmt.Printf("   [%d+%d] %s:%d\n",
				skip,
				i,
				path.Base(file),
				line)
		}
	}
				

}
