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

// BuiltinBuilder has no frills. It is a proxy to built-in go packages.
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

// signatured is an error that also has info about the function where it
// happened. It behaves like an error created with the builtin errors.New,
// except when processed with a function that is aware of its extra methods,
// like StackString.
type signatured struct {
	// message describes the error
	message string

	// fi stores the location of the error
	fi FuncInfo

	// argStringer has a String() method that describes the arguments provided to
	// the erroring function. The StackString function will add this between
	// parenthesis after the function name.
	argStringer interface{ String() string }
}

func (s *signatured) Error() string {
	return s.message
}

func (s *signatured) FuncInfo() FuncInfo {
	return s.fi
}

func (s *signatured) ArgStringer() interface{ String() string } {
	return s.argStringer
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
	argString string
}

// NewBuilder prefixes info about which function was called, and the args
// provided.
func NewBuilder(argFmt string, args ...interface{}) Builder {
	ab := new(argsBuilder)
	ab.argString = fmt.Sprintf(argFmt, args...)
	return ab
}

// Errorf is the same as fmt.Errorf, except that the error message gets a
// prefix with line number info and, if provided to NewBuilder, descriptions
// of the arguments passed to this function call.
func (ab *argsBuilder) Errorf(msg string, args ...interface{}) error {
	return &signatured{
		message:     fmt.Sprintf(msg, args...),
		fi:          NewFuncInfo(1),
		argStringer: stringStringer(ab.argString),
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
		message:     fmt.Sprintf(msg, args...),
		fi:          NewFuncInfo(1),
		argStringer: stringStringer(ab.argString),
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
		message:     fmt.Sprintf(msg, args...),
		fi:          NewFuncInfo(1),
		argStringer: formatStringer{argFmt, lab.args},
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
		message:     fmt.Sprintf(msg, args...),
		fi:          NewFuncInfo(1),
		argStringer: formatStringer{argFmt, lab.args},
	})
}
