package errors

import (
	"fmt"
	"path"
)

// Builder implementors can make and wrap errors.
type Builder interface {
	Errorf(msg string, args ...interface{}) error
	Wrap(err error, msg string, args ...interface{}) Wrapped
}

type builtinBuilder struct{}

// SimpleBuilder has no frills. It is a proxy to built-in go packages.
var BuiltinBuilder Builder = (*builtinBuilder)(nil)

// Errorf is the same as fmt.Errorf
func (bb *builtinBuilder) Errorf(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

// Wrap is the same as errors.Wrap
func (bb *builtinBuilder) Wrap(
	err error,
	msg string,
	args ...interface{},
) Wrapped {
	return Wrap(err, msg, args...)
}

type argsBuilder struct {
	argStr string
}

// NewBuilder prefixes info about which function was called, and the args
// provided.
func NewBuilder(argStr string, args ...interface{}) Builder {
	ab := new(argsBuilder)
	ab.argStr = fmt.Sprintf(argStr, args...)
	return ab
}

func (ab *argsBuilder) prefix() string {
	fi := NewFuncInfo(2)
	return fmt.Sprintf(
		"%s:%d %s(%s): ",
		path.Base(fi.File()),
		fi.Line(),
		fi.FuncName(),
		ab.argStr,
	)
}

// Errorf is the same as fmt.Errorf, except that the error message gets a
// prefix with line number info and, if provided to NewBuilder, descriptions
// of the arguments passed to this function call.
func (ab *argsBuilder) Errorf(msg string, args ...interface{}) error {
	return fmt.Errorf("%s%s", ab.prefix(), fmt.Sprintf(msg, args...))
}

// Wrap is the same as errors.Wrap, except that the error message gets a prefix
// with line number info and, if provided to NewBuilder, descriptions of the
// arguments passed to this function call.
func (ab *argsBuilder) Wrap(
	err error,
	msg string,
	args ...interface{},
) Wrapped {
	return Wrap(err, "%s%s", ab.prefix(), fmt.Sprintf(msg, args...))
}

type lazyArgsBuilder struct {
	argStr string
	args   []interface{}
}

// NewLazyBuilder SHOULD NOT be used unless it is known that NewBuilder
// won't work. Frequent undisciplined usage of NewLazyBuilder can lead to
// poor code maintainability. It is similar to NewBuilder, except that the
// prefix is formatted later, when the error actually occurs. This can be
// important for performance in functions called hundreds of times per second,
// but misleading debug messages can result if the arguments have changed since
// the function was first called. As a debug warning, errors are labeled
// "<lazyMsg>" to indicate that the argument values were formatted after the
// error appeared, not when the Builder was created.
func NewLazyBuilder(argStr string, args ...interface{}) Builder {
	lab := new(lazyArgsBuilder)
	lab.argStr = argStr
	lab.args = args
	return lab
}

func (lab *lazyArgsBuilder) prefix() string {
	fi := NewFuncInfo(2)
	argStr := fmt.Sprintf(lab.argStr, lab.args...)
	return fmt.Sprintf(
		"<lazyMsg> %s:%d %s(%s): ",
		path.Base(fi.File()),
		fi.Line(),
		fi.FuncName(),
		argStr,
	)
}

// Errorf is the same as fmt.Errorf, except that the error message gets a
// prefix with line number info and, if provided to NewLazyBuilder,
// descriptions of the arguments passed to this function call.
func (lab *lazyArgsBuilder) Errorf(msg string, args ...interface{}) error {
	return fmt.Errorf("%s%s", lab.prefix(), fmt.Sprintf(msg, args...))
}

// Wrap is the same as errors.Wrap, except that the error message gets a prefix
// with line number info and, if provided to NewLazyBuilder, descriptions of
// the arguments passed to this function call.
func (lab *lazyArgsBuilder) Wrap(
	err error,
	msg string,
	args ...interface{},
) Wrapped {
	return Wrap(err, "%s%s", lab.prefix(), fmt.Sprintf(msg, args...))
}
