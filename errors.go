package errors

import (
	"errors"
	"fmt"
	"strings"
)

var (
	As     = errors.As
	Is     = errors.Is
	New    = errors.New
	Unwrap = errors.Unwrap
)

type wrapped struct {
	error
	wrapped error
}

type Wrapped interface {
	error
	Unwrap() error
}

func Wrap(cause error, msg string, vars ...interface{}) Wrapped {
	return WrapWith(cause, fmt.Errorf(msg, vars...))
}

func WrapWith(cause, err error) Wrapped {
	w := new(wrapped)
	w.wrapped = cause
	w.error = err
	return w
}

func (w *wrapped) Unwrap() error {
	return w.wrapped
}

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
	// reverse the slice
	// for a, b := 0, len(errSlice)-1; a < b; a, b = a+1, b-1 {
	// 	errSlice[a], errSlice[b] = errSlice[b], errSlice[a]
	// }
}

func StackString(err error) string {
	stack := Stack(err)
	messages := make([]string, 0, len(stack))
	for _, err := range stack {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "\n")
}

type Group []error

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
