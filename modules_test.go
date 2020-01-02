package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainModule(t *testing.T) {
	mm := MainModule()
	if mm != "" {
		assert.Equal(t, "github.com/chaimleib/errors", mm)
	}
}

func TestRelativeModule(t *testing.T) {
	type testCase struct {
		mod, exp string
	}
	home := "github.com/chaimleib/errors"
	cases := []testCase{
		{"runtime", "runtime"},
		{"runtime/debug", "runtime/debug"},
		{"github.com/chaimleib/other", "github.com/chaimleib/other"},
		{home, "~"},
		{home + "/sub", "~/sub"},
	}
	for i, c := range cases {
		msg := fmt.Sprintf("case %d %+v", i, c)
		assert.Equal(t, c.exp, RelativeModule(c.mod, home), msg)
	}
	assert.Equal(t, "runtime/debug", RelativeModule("runtime/debug", "github.com/chaimleib/errors"))
}
