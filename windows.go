//go:build windows

package toolpaths

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func windowsRoamingAppData() string {
	// When userHomeDirFunc is overridden (for testing), use home-based paths
	if homeDirFuncOverridden {
		return filepath.Join(userHomeDir(), "AppData", "Roaming")
	}
	path, err := windows.KnownFolderPath(windows.FOLDERID_RoamingAppData, 0)
	if err != nil {
		// Fallback to environment variable
		return os.Getenv("APPDATA")
	}
	return path
}

func windowsLocalAppData() string {
	// When userHomeDirFunc is overridden (for testing), use home-based paths
	if homeDirFuncOverridden {
		return filepath.Join(userHomeDir(), "AppData", "Local")
	}
	path, err := windows.KnownFolderPath(windows.FOLDERID_LocalAppData, 0)
	if err != nil {
		// Fallback to environment variable
		return os.Getenv("LOCALAPPDATA")
	}
	return path
}

func windowsProgramData() string {
	path, err := windows.KnownFolderPath(windows.FOLDERID_ProgramData, 0)
	if err != nil {
		// Fallback to environment variable
		return os.Getenv("ProgramData")
	}
	return path
}
