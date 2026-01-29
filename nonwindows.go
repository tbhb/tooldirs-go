//go:build !windows

package toolpaths

import (
	"path/filepath"
)

// These functions provide sensible defaults for Windows path conventions
// on non-Windows platforms. This enables:
// 1. Testing Windows path resolution logic on any platform
// 2. Explicit Platform: PlatformWindows usage on non-Windows systems

func windowsRoamingAppData() string {
	return filepath.Join(userHomeDir(), "AppData", "Roaming")
}

func windowsLocalAppData() string {
	return filepath.Join(userHomeDir(), "AppData", "Local")
}

func windowsProgramData() string {
	return filepath.Join(string(filepath.Separator), "ProgramData")
}
