package errors

import (
	"runtime"
)

type funcInfo struct {
	file, funcName string
	line           int
}

func (fi *funcInfo) File() string {
	return fi.file
}

func (fi *funcInfo) FuncName() string {
	return fi.funcName
}

func (fi *funcInfo) Line() int {
	return fi.line
}

type FuncInfo interface {
	File() string
	Line() int
	FuncName() string
}

func NewFuncInfo(calldepth int) FuncInfo {
	var pc uintptr
	var ok bool
	fi := new(funcInfo)
	pc, fi.file, fi.line, ok = runtime.Caller(calldepth)
	if fi.file == "" {
		fi.file = "?file?"
	}
	if !ok {
		return nil
	}
	fi.funcName = runtime.FuncForPC(pc).Name()
	return fi
}
