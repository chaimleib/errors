package errors

import (
	"runtime"
)

type funcInfo struct {
	file, funcName string
	line           int
}

// File gives the absolute path to the go source file.
func (fi *funcInfo) File() string {
	return fi.file
}

// FuncName gives the fully-qualified identifier of the function. This includes
// the enclosing package name, with child packaged separated with slashes, and
// the function name appended after a dot.
func (fi *funcInfo) FuncName() string {
	return fi.funcName
}

// Line gives the line number.
func (fi *funcInfo) Line() int {
	return fi.line
}

// FuncInfo identifies a line of code, as provided by go runtime package
// functions.
type FuncInfo interface {
	File() string
	Line() int
	FuncName() string
}

// NewFuncInfo ascends the number of stack frames indicated by calldepth and
// returns a FuncInfo describing the location of that line of code.
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
