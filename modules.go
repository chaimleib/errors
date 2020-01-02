package errors

import (
	"runtime/debug"
	"strings"
)

// MainModule returns the path of the currently-running binary's main module.
// If not built with module support, returns "".
func MainModule() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return buildInfo.Main.Path
}

// RelativeModule replaces occurrences of home inside modName with a tilde `~`.
// If there is no common root in the first two elements (e.g.
// `github.com/username`), just return modName without changing it.
func RelativeModule(modName, home string) string {
	if !strings.Contains(home, "/") || !strings.HasPrefix(modName, home) {
		return modName
	}
	postHome := modName[len(home):]
	if postHome == "" || postHome[0] == '/' || postHome[0] == '.' {
		return "~" + postHome
	}
	return modName
}
