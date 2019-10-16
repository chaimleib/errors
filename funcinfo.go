package errors

import (
	"runtime"
)

func FuncInfo(calldepth int) (name, file string, line int) {
	var pc uintptr
	var ok bool
	pc, file, line, ok = runtime.Caller(calldepth)
	if file == "" {
		file = "?file?"
	}
	if !ok {
		return
	}
	name = runtime.FuncForPC(pc).Name()
	return
}
