//go:build darwin

package tooldirs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tbhb/tooldirs-go"
)

// Tests that verify auto-detection works correctly on macOS.
// These tests do NOT specify Platform explicitly - they rely on
// PlatformAuto detecting darwin and using macOS paths.

func TestAutoMacOSUserConfigDir(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, "Library", "Application Support", "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestAutoMacOSUserDataDir(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, "Library", "Application Support", "testapp")
	assert.Equal(t, expected, dirs.UserDataDir())
}

func TestAutoMacOSUserCacheDir(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, "Library", "Caches", "testapp")
	assert.Equal(t, expected, dirs.UserCacheDir())
}

func TestAutoMacOSUserLogDir(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, "Library", "Logs", "testapp")
	assert.Equal(t, expected, dirs.UserLogDir())
}

func TestAutoMacOSUserStateDir(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, "Library", "Application Support", "testapp")
	assert.Equal(t, expected, dirs.UserStateDir())
}

func TestAutoMacOSSystemConfigDir(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	expected := filepath.Join("/Library", "Application Support", "testapp")
	assert.Equal(t, expected, dirs.SystemConfigDir())
}

func TestAutoMacOSSystemCacheDir(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	expected := filepath.Join("/Library", "Caches", "testapp")
	assert.Equal(t, expected, dirs.SystemCacheDir())
}

func TestAutoMacOSSystemLogDir(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	expected := filepath.Join("/Library", "Logs", "testapp")
	assert.Equal(t, expected, dirs.SystemLogDir())
}

func TestAutoMacOSUserRuntimeDir(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	path, err := dirs.UserRuntimeDir()
	require.NoError(t, err)

	// On macOS, runtime uses TMPDIR which is per-user
	assert.Contains(t, path, "testapp")
}

func TestAutoMacOSWithVersion(t *testing.T) {
	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName: "testapp",
		Version: "2.0",
	})
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, "Library", "Application Support", "testapp", "2.0")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestAutoMacOSUserConfigDirs(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	configDirs := dirs.UserConfigDirs()
	require.Len(t, configDirs, 2)

	nativePath := filepath.Join(home, "Library", "Application Support", "testapp")
	xdgPath := filepath.Join(home, ".config", "testapp")

	assert.Equal(t, nativePath, configDirs[0])
	assert.Equal(t, xdgPath, configDirs[1])
}
