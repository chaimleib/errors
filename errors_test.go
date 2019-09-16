package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	inner := fmt.Errorf("a")
	outer := Wrap(inner, "b")
	assert.NotNil(t, outer)
	assert.Equal(t, inner, outer.Unwrap())
	assert.Equal(t, inner, Unwrap(outer))
	assert.Equal(t, "b", outer.Error())
	assert.Equal(t, []error{outer, inner}, Stack(outer))
}

type testError struct{}

func (te *testError) Error() string { return "test" }
func (te *testError) Unwrap() error { return fmt.Errorf("cause") }
func Test3rdPartyWrappedError(t *testing.T) {
	assert.Equal(t, "test\ncause", StackString(new(testError)))
}

func TestStackString(t *testing.T) {
	inner := fmt.Errorf("a")
	outer := Wrap(inner, "b")
	assert.Equal(t, "b\na", StackString(outer))
}

func TestGroup(t *testing.T) {
	assert.Equal(t, "", Group(nil).Error())

	a := fmt.Errorf("a")
	assert.Equal(t, "a", Group([]error{a}).Error())

	bc := Wrap(fmt.Errorf("c"), "b")
	assert.Equal(t, "[\n\ta\n\t,\n\tb\n\tc\n]", Group([]error{a, bc}).Error())
}
