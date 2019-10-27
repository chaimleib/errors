package errors

import (
	"fmt"
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
	w := new(wrapped)
	w.wrapped = cause
	w.error = err
	return w
}

// Unwrap returns the cause of the sender.
func (w *wrapped) Unwrap() error {
	return w.wrapped
}

// Is returns whether the wrapping error or a member of its cause chain matches
// the target.
func (w *wrapped) Is(target error) bool {
	switch target := target.(type) {
	case *wrapped:
		if w == target {
			return true
		}
	}
	return Is(w.error, target) || Is(w.wrapped, target)
}

// As assigns the first compatible error it finds in the cause chain to the
// target, and returns true if successful. If successful, Is(w, target) will be
// true.
func (w *wrapped) As(target interface{}) bool {
	switch target := target.(type) {
	case *wrapped:
		target.error = w.error
		target.wrapped = w.wrapped
		return true
	}
	return As(w.error, target) || As(w.wrapped, target)
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
func StackString(err error) string {
	stack := Stack(err)
	messages := make([]string, 0, len(stack))
	for _, err := range stack {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "\n")
}

// Group allows treating a slice of errors as an error. This is useful when
// many errors together lead to one error downstream. For example, a network
// fetch might fail only if all the mirrors fail to respond.
type Group []error

// Error formats a []error in an easily-readable way for human consumption.
func (g Group) Error() string {
	l := []error(g)
	switch len(l) {
	case 0:
		return ""
	case 1:
		return StackString(l[0])
	}
	messages := make([]string, 0, len(l))
	for _, err := range l {
		messages = append(messages, indent(StackString(err)))
	}
	return fmt.Sprintf("[\n%s\n]", strings.Join(messages, "\n\t,\n"))
}

func indent(s string) string {
	return "\t" + strings.ReplaceAll(s, "\n", "\n\t")
}
