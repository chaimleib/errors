package errors

import (
	"fmt"
	"path"
)

type ErrorMaker interface {
	Errorf(msg string, args ...interface{}) error
	Wrap(err error, msg string, args ...interface{}) Wrapped
}

type defaultErrorMaker struct{}

func (dem *defaultErrorMaker) Errorf(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

func (dem *defaultErrorMaker) Wrap(
	err error,
	msg string,
	args ...interface{},
) Wrapped {
	return Wrap(err, msg, args...)
}

var SimpleErrorMaker ErrorMaker = (*defaultErrorMaker)(nil)

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
	fi := NewFuncInfo(3)
	return fmt.Sprintf(
		"%s:%d %s(%s): ",
		path.Base(fi.File()),
		fi.Line(),
		fi.FuncName(),
		aem.argStr,
	)
}

func (aem *argsErrorMaker) Errorf(msg string, args ...interface{}) error {
	return fmt.Errorf("%s%s", aem.prefix(), fmt.Sprintf(msg, args...))
}

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
// important for performance for functions called hundreds of times per second,
// but misleading debug messages can result in the error if the arguments have
// changed since the function was first called. As a debug warning, errors are
// labeled "<lazyMsg>" to indicate that the argument values were formatted
// after the error appeared, not when the ErrorMaker was created.
func NewLazyErrorMaker(argStr string, args ...interface{}) ErrorMaker {
	laem := new(lazyArgsErrorMaker)
	laem.argStr = argStr
	laem.args = args
	return laem
}

func (laem *lazyArgsErrorMaker) prefix() string {
	fi := NewFuncInfo(3)
	argStr := fmt.Sprintf(laem.argStr, laem.args...)
	return fmt.Sprintf(
		"<lazyMsg> %s:%d %s(%s): ",
		path.Base(fi.File()),
		fi.Line(),
		fi.FuncName(),
		argStr,
	)
}

func (laem *lazyArgsErrorMaker) Errorf(msg string, args ...interface{}) error {
	return fmt.Errorf("%s%s", laem.prefix(), fmt.Sprintf(msg, args...))
}

func (laem *lazyArgsErrorMaker) Wrap(
	err error,
	msg string,
	args ...interface{},
) Wrapped {
	return Wrap(err, "%s%s", laem.prefix(), fmt.Sprintf(msg, args...))
}
