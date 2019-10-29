package errors

import (
	"fmt"
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

type signatured struct {
	message string
	fi      FuncInfo
	params  interface{ String() string }
}

func (s *signatured) Error() string {
	return s.message
}

func (s *signatured) Location() FuncInfo {
	return s.fi
}

func (s *signatured) Parameters() interface{ String() string } {
	return s.params
}

type stringStringer string

func (ss stringStringer) String() string {
	return string(ss)
}

type formatStringer struct {
	fmt    string
	params []interface{}
}

func (fs formatStringer) String() string {
	return fmt.Sprintf(fs.fmt, fs.params...)
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

// Errorf is the same as fmt.Errorf, except that the error message gets a
// prefix with line number info and, if provided to NewBuilder, descriptions
// of the arguments passed to this function call.
func (ab *argsBuilder) Errorf(msg string, args ...interface{}) error {
	return &signatured{
		message: fmt.Sprintf(msg, args...),
		fi:      NewFuncInfo(1),
		params:  stringStringer(ab.argStr),
	}
}

// Wrap is the same as errors.Wrap, except that the error message gets a prefix
// with line number info and, if provided to NewBuilder, descriptions of the
// arguments passed to this function call.
func (ab *argsBuilder) Wrap(
	err error,
	msg string,
	args ...interface{},
) Wrapped {
	return WrapWith(err, &signatured{
		message: fmt.Sprintf(msg, args...),
		fi:      NewFuncInfo(1),
		params:  stringStringer(ab.argStr),
	})
}

type lazyArgsBuilder struct {
	argFmt string
	args   []interface{}
}

// NewLazyBuilder SHOULD NOT be used unless it is known that NewBuilder
// won't work. Frequent undisciplined usage of NewLazyBuilder can lead to
// poor code maintainability. It is similar to NewBuilder, except that the
// prefix is formatted later, when the error actually occurs. This can be
// important for performance in functions called hundreds of times per second,
// but misleading debug messages can result if the arguments have changed since
// the function was first called. As a debug warning, any args are labeled
// "<lazy>" to indicate that the argument values were formatted after the
// error appeared, not when the Builder was created.
func NewLazyBuilder(argFmt string, args ...interface{}) Builder {
	lab := new(lazyArgsBuilder)
	lab.argFmt = argFmt
	lab.args = args
	return lab
}

// Errorf is the same as fmt.Errorf, except that the error message gets a
// prefix with line number info and, if provided to NewLazyBuilder,
// descriptions of the arguments passed to this function call.
func (lab *lazyArgsBuilder) Errorf(msg string, args ...interface{}) error {
	argFmt := lab.argFmt
	if len(argFmt) != 0 { // if no args, nothing to warn about
		argFmt = "<lazy> " + argFmt
	}
	return &signatured{
		message: fmt.Sprintf(msg, args...),
		fi:      NewFuncInfo(1),
		params:  formatStringer{argFmt, lab.args},
	}
}

// Wrap is the same as errors.Wrap, except that the error message gets a prefix
// with line number info and, if provided to NewLazyBuilder, descriptions of
// the arguments passed to this function call.
func (lab *lazyArgsBuilder) Wrap(
	err error,
	msg string,
	args ...interface{},
) Wrapped {
	argFmt := lab.argFmt
	if len(argFmt) != 0 { // if no args, nothing to warn about
		argFmt = "<lazy> " + argFmt
	}
	return WrapWith(err, &signatured{
		message: fmt.Sprintf(msg, args...),
		fi:      NewFuncInfo(1),
		params:  formatStringer{argFmt, lab.args},
	})
}
