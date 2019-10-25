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
	assert.Regexp(
		t,
		`^errormaker_test\.go:[0-9]+ [^ ]+/errors\.TestNewBuilder\(t, map\["first":"orig"\]\): a$`,
		err.Error(),
	)
}

func TestNewLazyBuilder(t *testing.T) {
	changingMap := map[string]string{"first": "orig"}
	b := NewLazyBuilder("t, %q", changingMap)
	changingMap["second"] = "new"
	err := b.Errorf("a")
	assert.Regexp(
		t,
		`^<lazyMsg> errormaker_test\.go:[0-9]+ [^ ]+/errors\.TestNewLazyBuilder\(t, map\[("[^"]+":"[^"]+" ?){2}\]\): a$`,
		err.Error(),
	)
}
