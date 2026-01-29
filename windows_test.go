//go:build windows

package tooldirs_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tbhb/tooldirs-go"
)

// Tests specific to Windows that verify real Windows API behavior.
// These complement auto_windows_test.go by testing Windows-specific
// edge cases and API integration.

func TestWindowsKnownFolderPathsUsed(t *testing.T) {
	// Verify that the Windows API (KnownFolderPath) returns paths
	// that match the expected environment variables
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	localAppData := os.Getenv("LOCALAPPDATA")
	require.NotEmpty(t, localAppData)

	// The actual path should use the KnownFolderPath result,
	// which should match LOCALAPPDATA
	configDir := dirs.UserConfigDir()
	assert.True(t, strings.HasPrefix(configDir, localAppData),
		"UserConfigDir should be under LOCALAPPDATA")
}

func TestWindowsRoamingUsesAppData(t *testing.T) {
	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName: "testapp",
		Roaming: true,
	})
	require.NoError(t, err)

	appData := os.Getenv("APPDATA")
	require.NotEmpty(t, appData)

	configDir := dirs.UserConfigDir()
	assert.True(t, strings.HasPrefix(configDir, appData),
		"Roaming UserConfigDir should be under APPDATA")
}

func TestWindowsSystemUsesProgramData(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	programData := os.Getenv("ProgramData")
	require.NotEmpty(t, programData)

	sysConfigDir := dirs.SystemConfigDir()
	assert.True(t, strings.HasPrefix(sysConfigDir, programData),
		"SystemConfigDir should be under ProgramData")
}
