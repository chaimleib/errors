package errors

import (
	"errors"
)

var (
	// provided by Go 1.13
	As     = errors.As
	Is     = errors.Is
	New    = errors.New
	Unwrap = errors.Unwrap
)
