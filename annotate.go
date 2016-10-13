package logging

type Annotater interface {
	Annotate(*Record)
}

const secretLevel = Level(99)

func grow(arg map[string]interface{}) map[string]interface{} {
	n := len(arg)
	tbl := make(map[string]interface{}, n+1)
	for k, v := range arg {
		tbl[k] = v
	}
	return tbl
}

func (l *Logger) Re(a Annotater) *Logger {
	n := &Logger{
		module: l.module,
		annot: grow(l.annot),
		outputs: l.outputs,
	}
	// now add whatever new annotations we have in mind
	var tmp Record
	tmp.Module = n.module
	tmp.Annotations = n.annot
	tmp.Level = secretLevel
	a.Annotate(&tmp)
	return n
}

func (a *Record) Annotate(key string, val interface{}) {
	// secret level used to indicate that this record is
	// a prototype use by Re()
	if a.Level != secretLevel {
		// need to make a copy of the annotations so we
		// don't side-effect the logger's
		a.Annotations = grow(a.Annotations)
	}
	a.Annotations[key] = val
}


	
	

