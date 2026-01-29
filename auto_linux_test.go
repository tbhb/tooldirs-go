//go:build linux

package toolpaths_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tbhb/toolpaths-go"
)

// Tests that verify auto-detection works correctly on Linux.
// These tests do NOT specify Platform explicitly - they rely on
// PlatformAuto detecting linux and using XDG paths.

func TestAutoLinuxUserConfigDir(t *testing.T) {
	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, ".config", "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestAutoLinuxUserDataDir(t *testing.T) {
	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, ".local", "share", "testapp")
	assert.Equal(t, expected, dirs.UserDataDir())
}

func TestAutoLinuxUserCacheDir(t *testing.T) {
	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, ".cache", "testapp")
	assert.Equal(t, expected, dirs.UserCacheDir())
}

func TestAutoLinuxUserStateDir(t *testing.T) {
	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, ".local", "state", "testapp")
	assert.Equal(t, expected, dirs.UserStateDir())
}

func TestAutoLinuxUserLogDir(t *testing.T) {
	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	// Log is a subdirectory of state on XDG platforms
	expected := filepath.Join(home, ".local", "state", "testapp", "log")
	assert.Equal(t, expected, dirs.UserLogDir())
}

func TestAutoLinuxSystemConfigDir(t *testing.T) {
	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	expected := filepath.Join("/etc", "xdg", "testapp")
	assert.Equal(t, expected, dirs.SystemConfigDir())
}

func TestAutoLinuxSystemDataDir(t *testing.T) {
	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	// First dir is /usr/local/share
	expected := filepath.Join("/usr", "local", "share", "testapp")
	assert.Equal(t, expected, dirs.SystemDataDir())
}

func TestAutoLinuxSystemDataDirs(t *testing.T) {
	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	sysDirs := dirs.SystemDataDirs()
	require.Len(t, sysDirs, 2)
	assert.Equal(t, filepath.Join("/usr", "local", "share", "testapp"), sysDirs[0])
	assert.Equal(t, filepath.Join("/usr", "share", "testapp"), sysDirs[1])
}

func TestAutoLinuxSystemCacheDir(t *testing.T) {
	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	expected := filepath.Join("/var", "cache", "testapp")
	assert.Equal(t, expected, dirs.SystemCacheDir())
}

func TestAutoLinuxSystemStateDir(t *testing.T) {
	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	expected := filepath.Join("/var", "lib", "testapp")
	assert.Equal(t, expected, dirs.SystemStateDir())
}

func TestAutoLinuxSystemLogDir(t *testing.T) {
	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	expected := filepath.Join("/var", "log", "testapp")
	assert.Equal(t, expected, dirs.SystemLogDir())
}

func TestAutoLinuxSystemRuntimeDir(t *testing.T) {
	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	expected := filepath.Join("/run", "testapp")
	assert.Equal(t, expected, dirs.SystemRuntimeDir())
}

func TestAutoLinuxWithVersion(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName: "testapp",
		Version: "3.0",
	})
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, ".config", "testapp", "3.0")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestAutoLinuxUserConfigDirsNoFallback(t *testing.T) {
	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	// On Linux (XDG platform), there are no fallbacks since XDG is native
	configDirs := dirs.UserConfigDirs()
	require.Len(t, configDirs, 1)

	expected := filepath.Join(home, ".config", "testapp")
	assert.Equal(t, expected, configDirs[0])
}

func TestAutoLinuxXDGRuntimeDir(t *testing.T) {
	testDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", testDir)

	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	path, err := dirs.UserRuntimeDir()
	require.NoError(t, err)

	expected := filepath.Join(testDir, "testapp")
	assert.Equal(t, expected, path)
}

func TestAutoLinuxRuntimeDirFallback(t *testing.T) {
	// When XDG_RUNTIME_DIR is not set, falls back to temp dir
	t.Setenv("XDG_RUNTIME_DIR", "")

	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	path, err := dirs.UserRuntimeDir()
	require.NoError(t, err)

	// Should contain app name and uid
	assert.Contains(t, path, "testapp")
}

func TestAutoLinuxXDGConfigDirsEnv(t *testing.T) {
	testDir := t.TempDir()
	t.Setenv("XDG_CONFIG_DIRS", testDir+":/opt/config")

	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	sysDirs := dirs.SystemConfigDirs()
	require.Len(t, sysDirs, 2)
	assert.Equal(t, filepath.Join(testDir, "testapp"), sysDirs[0])
	assert.Equal(t, filepath.Join("/opt", "config", "testapp"), sysDirs[1])
}

func TestAutoLinuxXDGDataDirsEnv(t *testing.T) {
	testDir := t.TempDir()
	t.Setenv("XDG_DATA_DIRS", testDir+":/opt/data")

	dirs, err := toolpaths.New("testapp")
	require.NoError(t, err)

	sysDirs := dirs.SystemDataDirs()
	require.Len(t, sysDirs, 2)
	assert.Equal(t, filepath.Join(testDir, "testapp"), sysDirs[0])
	assert.Equal(t, filepath.Join("/opt", "data", "testapp"), sysDirs[1])
}
