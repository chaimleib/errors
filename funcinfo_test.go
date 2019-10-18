package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncInfo(t *testing.T) {
	fi := NewFuncInfo(1)
	assert.NotNil(t, fi)
	assert.Contains(t, fi.FuncName(), "errors.TestFuncInfo")
	assert.Contains(t, fi.File(), "/funcinfo_test.go")
	assert.LessOrEqual(t, 10, fi.Line())
}
