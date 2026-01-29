//go:build windows

package tooldirs

import (
	"os"

	"golang.org/x/sys/windows"
)

func windowsRoamingAppData() string {
	path, err := windows.KnownFolderPath(windows.FOLDERID_RoamingAppData, 0)
	if err != nil {
		// Fallback to environment variable
		return os.Getenv("APPDATA")
	}
	return path
}

func windowsLocalAppData() string {
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
