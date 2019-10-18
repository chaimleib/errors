package errors

import (
	"fmt"
	"path"
)

// ErrorMaker implementors can make and wrap errors.
type ErrorMaker interface {
	Errorf(msg string, args ...interface{}) error
	Wrap(err error, msg string, args ...interface{}) Wrapped
}

type defaultErrorMaker struct{}

// SimpleErrorMaker has no frills. It is a proxy to built-in go packages.
var SimpleErrorMaker ErrorMaker = (*defaultErrorMaker)(nil)

// Errorf is the same as fmt.Errorf
func (dem *defaultErrorMaker) Errorf(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

// Wrap is the same as errors.Wrap
func (dem *defaultErrorMaker) Wrap(
	err error,
	msg string,
	args ...interface{},
) Wrapped {
	return Wrap(err, msg, args...)
}

type argsErrorMaker struct {
	argStr string
}

// NewErrorMaker prefixes info about which function was called, and the args
// provided.
func NewErrorMaker(argStr string, args ...interface{}) ErrorMaker {
	aem := new(argsErrorMaker)
	aem.argStr = fmt.Sprintf(argStr, args...)
	return aem
}

func (aem *argsErrorMaker) prefix() string {
	fi := NewFuncInfo(2)
	return fmt.Sprintf(
		"%s:%d %s(%s): ",
		path.Base(fi.File()),
		fi.Line(),
		fi.FuncName(),
		aem.argStr,
	)
}

// Errorf is the same as fmt.Errorf, except that the error message gets a
// prefix with line number info and, if provided to NewErrorMaker, descriptions
// of the arguments passed to this function call.
func (aem *argsErrorMaker) Errorf(msg string, args ...interface{}) error {
	return fmt.Errorf("%s%s", aem.prefix(), fmt.Sprintf(msg, args...))
}

// Wrap is the same as errors.Wrap, except that the error message gets a prefix
// with line number info and, if provided to NewErrorMaker, descriptions of the
// arguments passed to this function call.
func (aem *argsErrorMaker) Wrap(
	err error,
	msg string,
	args ...interface{},
) Wrapped {
	return Wrap(err, "%s%s", aem.prefix(), fmt.Sprintf(msg, args...))
}

type lazyArgsErrorMaker struct {
	argStr string
	args   []interface{}
}

// NewLazyErrorMaker SHOULD NOT be used unless it is known that NewErrorMaker
// won't work. Frequent undisciplined usage of NewLazyErrorMaker can lead to
// poor code maintainability. It is similar to NewErrorMaker, except that the
// prefix is formatted later, when the error actually occurs. This can be
// important for performance in functions called hundreds of times per second,
// but misleading debug messages can result if the arguments have changed since
// the function was first called. As a debug warning, errors are labeled
// "<lazyMsg>" to indicate that the argument values were formatted after the
// error appeared, not when the ErrorMaker was created.
func NewLazyErrorMaker(argStr string, args ...interface{}) ErrorMaker {
	laem := new(lazyArgsErrorMaker)
	laem.argStr = argStr
	laem.args = args
	return laem
}

func (laem *lazyArgsErrorMaker) prefix() string {
	fi := NewFuncInfo(2)
	argStr := fmt.Sprintf(laem.argStr, laem.args...)
	return fmt.Sprintf(
		"<lazyMsg> %s:%d %s(%s): ",
		path.Base(fi.File()),
		fi.Line(),
		fi.FuncName(),
		argStr,
	)
}

// Errorf is the same as fmt.Errorf, except that the error message gets a
// prefix with line number info and, if provided to NewLazyErrorMaker,
// descriptions of the arguments passed to this function call.
func (laem *lazyArgsErrorMaker) Errorf(msg string, args ...interface{}) error {
	return fmt.Errorf("%s%s", laem.prefix(), fmt.Sprintf(msg, args...))
}

// Wrap is the same as errors.Wrap, except that the error message gets a prefix
// with line number info and, if provided to NewLazyErrorMaker, descriptions of
// the arguments passed to this function call.
func (laem *lazyArgsErrorMaker) Wrap(
	err error,
	msg string,
	args ...interface{},
) Wrapped {
	return Wrap(err, "%s%s", laem.prefix(), fmt.Sprintf(msg, args...))
}
