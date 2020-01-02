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

// FuncInfo returns the location of the error.
func (s *signatured) FuncInfo() FuncInfo {
	return s.fi
}

// ArgStringer returns a value with a String() method, which describes the
// arguments the erroring function was called with. This can be printed within
// parenthesis after the function name for debugging, as in StackString.
func (s *signatured) ArgStringer() interface{ String() string } {
	return s.argStringer
}

// stringStringer has a String method that returns the underlying string value.
type stringStringer string

func (ss stringStringer) String() string {
	return string(ss)
}

// formatStringer has a String method that calls fmt.Sprintf with predefined
// arguments.
type formatStringer struct {
	fmt    string
	params []interface{}
}

func (fs formatStringer) String() string {
	return fmt.Sprintf(fs.fmt, fs.params...)
}

// argsBuilder is the underlying type for NewBuilder.
type argsBuilder struct {
	argString string
}

// NewBuilder returns an error builder that attaches info about the function
// where the error happened, and the args with which the function was called.
func NewBuilder(argFmt string, args ...interface{}) Builder {
	ab := new(argsBuilder)
	ab.argString = fmt.Sprintf(argFmt, args...)
	return ab
}

// Errorf is the same as fmt.Errorf, except that the error message gets
// FuncInfo() and ArgStringer() methods, describing the context of the error.
func (ab *argsBuilder) Errorf(msg string, args ...interface{}) error {
	return &signatured{
		message:     fmt.Sprintf(msg, args...),
		fi:          NewFuncInfo(1),
		argStringer: stringStringer(ab.argString),
	}
}

// Wrap replaces errors.Wrap, except that the error additionally implements the
// Wrapper() method, whose return value implements the FuncInfo() and
// ArgStringer() methods to describe the context of the error.
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
// embedded ArgStringer()s do their formatting work lazily when the error is
// formatted for display, rather than up-front when the Builder is initialized.
// This can be important for performance in functions called thousands of times
// per second, but misleading debug messages can result if the arguments have
// changed since the function was first called. As a debug warning, any args
// are labeled "<lazy>" by the ArgStringer().
func NewLazyBuilder(argFmt string, args ...interface{}) Builder {
	lab := new(lazyArgsBuilder)
	lab.argFmt = argFmt
	lab.args = args
	return lab
}

// Errorf is the same as fmt.Errorf, except that the error message gets
// FuncInfo() and ArgStringer() methods, describing the context of the error.
// On lazy builders, ArgStringer() does its formatting computations when its
// String() method gets called.
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

// Wrap replaces errors.Wrap, except that the error additionally implements the
// Wrapper() method, whose return value implements the FuncInfo() and
// ArgStringer() methods to describe the context of the error.
// On lazy builders, ArgStringer() does its formatting computations when its
// String() method gets called.
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
