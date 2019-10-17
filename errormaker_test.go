package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleErrorMaker(t *testing.T) {
	e := SimpleErrorMaker.Errorf("a")
	assert.Error(t, e)
	assert.Equal(t, "a", e.Error())
}

func TestNewErrorMaker(t *testing.T) {
	changingMap := map[string]string{"first": "orig"}
	em := NewErrorMaker("t, %q", changingMap)
	changingMap["second"] = "new"
	e := em.Errorf("a")
	assert.Regexp(
		t,
		`^errormaker_test\.go:[0-9]+ [^ ]+/errors\.TestNewErrorMaker\(t, map\["first":"orig"\]\): a$`,
		e.Error(),
	)
}

func TestNewLazyErrorMaker(t *testing.T) {
	changingMap := map[string]string{"first": "orig"}
	em := NewLazyErrorMaker("t, %q", changingMap)
	changingMap["second"] = "new"
	e := em.Errorf("a")
	assert.Regexp(
		t,
		`^<lazyMsg> errormaker_test\.go:[0-9]+ [^ ]+/errors\.TestNewLazyErrorMaker\(t, map\[("[^"]+":"[^"]+" ?){2}\]\): a$`,
		e.Error(),
	)
}
