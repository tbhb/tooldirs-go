//go:build windows

package tooldirs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tbhb/tooldirs-go"
)

// Tests that verify auto-detection works correctly on Windows.
// These tests do NOT specify Platform explicitly - they rely on
// PlatformAuto detecting windows and using Windows paths.

func TestAutoWindowsUserConfigDir(t *testing.T) {
	localAppData := os.Getenv("LOCALAPPDATA")
	require.NotEmpty(t, localAppData, "LOCALAPPDATA should be set on Windows")

	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	expected := filepath.Join(localAppData, "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestAutoWindowsUserConfigDirRoaming(t *testing.T) {
	appData := os.Getenv("APPDATA")
	require.NotEmpty(t, appData, "APPDATA should be set on Windows")

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName: "testapp",
		Roaming: true,
	})
	require.NoError(t, err)

	expected := filepath.Join(appData, "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestAutoWindowsUserDataDir(t *testing.T) {
	localAppData := os.Getenv("LOCALAPPDATA")
	require.NotEmpty(t, localAppData, "LOCALAPPDATA should be set on Windows")

	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	expected := filepath.Join(localAppData, "testapp")
	assert.Equal(t, expected, dirs.UserDataDir())
}

func TestAutoWindowsUserCacheDir(t *testing.T) {
	localAppData := os.Getenv("LOCALAPPDATA")
	require.NotEmpty(t, localAppData, "LOCALAPPDATA should be set on Windows")

	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	// Cache is always under LOCALAPPDATA
	expected := filepath.Join(localAppData, "testapp", "cache")
	assert.Equal(t, expected, dirs.UserCacheDir())
}

func TestAutoWindowsUserLogDir(t *testing.T) {
	localAppData := os.Getenv("LOCALAPPDATA")
	require.NotEmpty(t, localAppData, "LOCALAPPDATA should be set on Windows")

	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	// Log is always under LOCALAPPDATA
	expected := filepath.Join(localAppData, "testapp", "log")
	assert.Equal(t, expected, dirs.UserLogDir())
}

func TestAutoWindowsUserStateDir(t *testing.T) {
	localAppData := os.Getenv("LOCALAPPDATA")
	require.NotEmpty(t, localAppData, "LOCALAPPDATA should be set on Windows")

	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	expected := filepath.Join(localAppData, "testapp")
	assert.Equal(t, expected, dirs.UserStateDir())
}

func TestAutoWindowsWithAppAuthor(t *testing.T) {
	localAppData := os.Getenv("LOCALAPPDATA")
	require.NotEmpty(t, localAppData, "LOCALAPPDATA should be set on Windows")

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:   "testapp",
		AppAuthor: "MyCompany",
	})
	require.NoError(t, err)

	expected := filepath.Join(localAppData, "MyCompany", "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestAutoWindowsWithVersion(t *testing.T) {
	localAppData := os.Getenv("LOCALAPPDATA")
	require.NotEmpty(t, localAppData, "LOCALAPPDATA should be set on Windows")

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName: "testapp",
		Version: "2.0",
	})
	require.NoError(t, err)

	expected := filepath.Join(localAppData, "testapp", "2.0")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestAutoWindowsWithAppAuthorAndVersion(t *testing.T) {
	localAppData := os.Getenv("LOCALAPPDATA")
	require.NotEmpty(t, localAppData, "LOCALAPPDATA should be set on Windows")

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:   "testapp",
		AppAuthor: "MyCompany",
		Version:   "2.0",
	})
	require.NoError(t, err)

	expected := filepath.Join(localAppData, "MyCompany", "testapp", "2.0")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestAutoWindowsSystemConfigDir(t *testing.T) {
	programData := os.Getenv("ProgramData")
	require.NotEmpty(t, programData, "ProgramData should be set on Windows")

	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	expected := filepath.Join(programData, "testapp")
	assert.Equal(t, expected, dirs.SystemConfigDir())
}

func TestAutoWindowsSystemRuntimeDir(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	// Windows has no system runtime dir
	assert.Empty(t, dirs.SystemRuntimeDir())
}

func TestAutoWindowsUserConfigDirs(t *testing.T) {
	localAppData := os.Getenv("LOCALAPPDATA")
	require.NotEmpty(t, localAppData, "LOCALAPPDATA should be set on Windows")

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	// Default: IncludeXDGFallbacks is true
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	configDirs := dirs.UserConfigDirs()
	require.Len(t, configDirs, 2)

	nativePath := filepath.Join(localAppData, "testapp")
	xdgPath := filepath.Join(home, ".config", "testapp")

	assert.Equal(t, nativePath, configDirs[0])
	assert.Equal(t, xdgPath, configDirs[1])
}

func TestAutoWindowsXDGOnAllPlatforms(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:           "testapp",
		XDGOnAllPlatforms: true,
	})
	require.NoError(t, err)

	// With XDGOnAllPlatforms, should use XDG paths
	expected := filepath.Join(home, ".config", "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}
