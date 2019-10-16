package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncInfo(t *testing.T) {
	name, file, line := FuncInfo(1)
	assert.Contains(t, name, "errors.TestFuncInfo")
	assert.Contains(t, file, "errors/funcinfo_test.go")
	assert.LessOrEqual(t, 10, line)
}
