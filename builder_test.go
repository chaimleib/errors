package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuiltinBuilder(t *testing.T) {
	b := BuiltinBuilder.Errorf("a")
	assert.Error(t, b)
	assert.Equal(t, "a", b.Error())
}

func TestNewBuilder(t *testing.T) {
	changingMap := map[string]string{"first": "orig"}
	b := NewBuilder("t, %q", changingMap)
	changingMap["second"] = "new"
	err := b.Errorf("a")
	assert.Equal(t, "a", err.Error())
	assert.Regexp(
		t,
		`^[^ ]+/errors\.TestNewBuilder\(t, map\["first":"orig"\]\) builder_test\.go:[0-9]+ a$`,
		StackString(err),
	)
}

func TestNewLazyBuilder(t *testing.T) {
	changingMap := map[string]string{"first": "orig"}
	b := NewLazyBuilder("t, %q", changingMap)
	changingMap["second"] = "new"
	err := b.Errorf("a")
	assert.Equal(t, "a", err.Error())
	assert.Regexp(
		t,
		`^[^ ]+/errors\.TestNewLazyBuilder\(<lazy> t, map\[("[^"]+":"[^"]+" ?){2}\]\) builder_test\.go:[0-9]+ a$`,
		StackString(err),
	)
}
