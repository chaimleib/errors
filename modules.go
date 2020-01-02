package errors

import (
	"runtime/debug"
	"strings"
)

// MainModule returns the path of the currently-running binary's main module,
// as defined in the go.mod file. If not built with module support, returns
// "".
func MainModule() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return buildInfo.Main.Path
}

// RelativeModule replaces occurrences of `home` inside `modName` with a tilde
// "~".  If `home` is the empty string or if `modName` is not a child of home,
// just return `modName` without changing it.
func RelativeModule(modName, home string) string {
	if home == "" || !strings.HasPrefix(modName, home) {
		return modName
	}
	postHome := modName[len(home):]
	if postHome == "" || postHome[0] == '/' || postHome[0] == '.' {
		return "~" + postHome
	}
	return modName
}
