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

func TestWrapIs(t *testing.T) {
	tokenErr := fmt.Errorf("token")
	a := Wrap(tokenErr, "a")
	b := Wrap(tokenErr, "b")
	unrelated := fmt.Errorf("unrelated")

	assert.True(t, Is(tokenErr, tokenErr))
	assert.True(t, Is(a, tokenErr))
	assert.True(t, Is(b, tokenErr))
	assert.False(t, Is(a, b))
	assert.False(t, Is(b, a))
	assert.False(t, Is(unrelated, tokenErr))
	assert.False(t, Is(a, unrelated))
	assert.False(t, Is(b, unrelated))
}

func TestWrapWithIs(t *testing.T) {
	tokenErr := fmt.Errorf("token")
	a := fmt.Errorf("a")
	a1 := WrapWith(tokenErr, a)
	a2 := WrapWith(tokenErr, a)
	a3 := Wrap(tokenErr, "a")
	unrelated := fmt.Errorf("unrelated")

	assert.True(t, Is(a1, a))
	assert.True(t, Is(a2, a))
	assert.False(t, Is(a3, a))
	assert.True(t, Is(a1, a2))
	assert.True(t, Is(a2, a1))
	assert.False(t, Is(a1, a3))
	assert.False(t, Is(a3, a1))
	assert.False(t, Is(a2, a3))
	assert.False(t, Is(a3, a2))
	assert.False(t, Is(a1, unrelated))
	assert.False(t, Is(a2, unrelated))
	assert.False(t, Is(a3, unrelated))
	assert.True(t, Is(a1, tokenErr))
	assert.True(t, Is(a2, tokenErr))
	assert.True(t, Is(a3, tokenErr))
	assert.False(t, Is(tokenErr, a1))
	assert.False(t, Is(tokenErr, a2))
	assert.False(t, Is(tokenErr, a3))
}
