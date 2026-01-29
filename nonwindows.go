//go:build !windows

package tooldirs

import (
	"os"
	"path/filepath"
)

// These functions provide sensible defaults for Windows path conventions
// on non-Windows platforms. This enables:
// 1. Testing Windows path resolution logic on any platform
// 2. Explicit Platform: PlatformWindows usage on non-Windows systems

func windowsRoamingAppData() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "AppData", "Roaming")
}

func windowsLocalAppData() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "AppData", "Local")
}

func windowsProgramData() string {
	return filepath.Join(string(filepath.Separator), "ProgramData")
}
