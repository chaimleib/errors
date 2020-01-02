package errors

import (
	"fmt"
	"path"
	"strings"
)

type wrapped struct {
	error
	wrapped error
}

// Wrapped is an error whose cause can be retrieved with Unwrap()
type Wrapped interface {
	error
	Unwrap() error
}

// Wrap is like fmt.Errorf, except that calling Unwrap() on the result yields
// the provided cause.
func Wrap(cause error, msg string, vars ...interface{}) Wrapped {
	return WrapWith(cause, fmt.Errorf(msg, vars...))
}

// WrapWith takes two errors and wraps the first provided error with the
// second. This is particularly useful when the wrapping error itself is a
// special type with custom fields and methods.
func WrapWith(cause, err error) Wrapped {
	return wrapped{error: err, wrapped: cause}
}

// Unwrap returns the cause of the sender.
func (w wrapped) Unwrap() error {
	return w.wrapped
}

// Wrapper returns the error value without the error it wraps. This allows
// access to custom fields of the error, without interference from the Unwrap
// chain like in Is or As.
func (w wrapped) Wrapper() error {
	return w.error
}

// Is returns whether the wrapping error or a member of its cause chain matches
// the target.
func (w wrapped) Is(target error) bool {
	switch target := target.(type) {
	case wrapped:
		if w == target {
			return true
		}
	}
	return Is(w.Wrapper(), target) || Is(w.Unwrap(), target)
}

// As assigns the first compatible error it finds in the cause chain to the
// target, and returns true if successful. If successful, Is(w, target) will be
// true.
func (w wrapped) As(target interface{}) bool {
	switch target := target.(type) {
	case wrapped:
		target.error = w.error
		target.wrapped = w.wrapped
		return true
	}
	return As(w.Wrapper(), target) || As(w.Unwrap(), target)
}

// Stack returns a slice of all the errors found by recursively calling
// Unwrap() on the provided error. Errors causing other errors appear later.
func Stack(err error) []error {
	var errSlice []error
	for err != nil {
		errSlice = append(errSlice, err)
		if w, ok := err.(Wrapped); ok {
			err = w.Unwrap()
			continue
		}
		break
	}
	return errSlice
}

// StackString recursively calls Unwrap() on the given error and stringifies
// all the errors in the chain.
//
// The stringification adds location prefixes to errors that additionally
// implement `FuncInfo() FuncInfo` and optionally `ArgStringer() interface{
// String() string }`.
//
// The stringification can be overridden if the error implements `StackString()
// string`.
func StackString(err error) string {
	stack := Stack(err)
	messages := make([]string, 0, len(stack))
	for _, err := range stack {
		messages = append(messages, stackStringAt(err))
	}
	return strings.Join(messages, "\n")
}

// stackStringAt returns one level of StackString.
func stackStringAt(err error) string {
	switch err := err.(type) {
	case interface{ StackString() string }:
		return err.StackString()
	case interface {
		error
		FuncInfo() FuncInfo
		ArgStringer() interface{ String() string }
	}:
		return fmt.Sprintf(
			"%s(%s) %s:%d %v",
			RelativeModule(err.FuncInfo().FuncName(), MainModule()),
			err.ArgStringer().String(),
			path.Base(err.FuncInfo().File()),
			err.FuncInfo().Line(),
			err,
		)
	case interface {
		error
		FuncInfo() FuncInfo
	}:
		return fmt.Sprintf(
			"%s %s:%d %v",
			RelativeModule(err.FuncInfo().FuncName(), MainModule()),
			path.Base(err.FuncInfo().File()),
			err.FuncInfo().Line(),
			err,
		)
	case interface {
		error
		Wrapper() error
	}:
		return stackStringAt(err.Wrapper())
	}
	return err.Error()
}

// Group allows treating a slice of errors as an error. This is useful when
// many errors together lead to one error downstream. For example, a network
// fetch might fail only if all the mirrors fail to respond.
type Group []error

// StackString formats a []error as a recursive list of StackString-ed errors
// if len() is 2 or more, otherwise as a single StackString-ed error if len()
// is 1, otherwise as nothing if len() is 0.
func (g Group) StackString() string {
	l := []error(g)
	switch len(l) {
	case 0:
		return ""
	case 1:
		return StackString(l[0])
	}
	messages := make([]string, 0, len(l))
	for _, err := range l {
		messages = append(messages, StackString(err))
	}
	return fmt.Sprintf("[\n%s\n]", indent(strings.Join(messages, "\n,\n")))
}

// Error formats a []error as a list of errors if len() is 2 or more, otherwise
// as a single error if len() is 1, otherwise as nothing if len() is 0.
func (g Group) Error() string {
	l := []error(g)
	switch len(l) {
	case 0:
		return ""
	case 1:
		return l[0].Error()
	}
	messages := make([]string, 0, len(l))
	for _, err := range l {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("[\n%s\n]", indent(strings.Join(messages, "\n,\n")))
}

func indent(s string) string {
	return "\t" + strings.ReplaceAll(s, "\n", "\n\t")
}
