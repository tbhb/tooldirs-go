package toolpaths_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tbhb/toolpaths-go"
)

// Tests for Linux/XDG path resolution logic.
// These use explicit Platform: PlatformLinux and mock HOME,
// so they can run on any platform.

// setTestHome sets up a test home directory and returns a cleanup function.
func setTestHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	// Clear XDG env vars that might override home-based paths
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("XDG_DATA_HOME", "")
	t.Setenv("XDG_CACHE_HOME", "")
	t.Setenv("XDG_STATE_HOME", "")
	t.Setenv("XDG_RUNTIME_DIR", "")
	toolpaths.SetHomeDirFunc(func() string { return home })
	t.Cleanup(func() { toolpaths.SetHomeDirFunc(nil) })
	return home
}

func TestLinuxPlatformUserConfigDir(t *testing.T) {
	home := setTestHome(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, ".config", "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestLinuxPlatformUserDataDir(t *testing.T) {
	home := setTestHome(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, ".local", "share", "testapp")
	assert.Equal(t, expected, dirs.UserDataDir())
}

func TestLinuxPlatformUserCacheDir(t *testing.T) {
	home := setTestHome(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, ".cache", "testapp")
	assert.Equal(t, expected, dirs.UserCacheDir())
}

func TestLinuxPlatformUserStateDir(t *testing.T) {
	home := setTestHome(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, ".local", "state", "testapp")
	assert.Equal(t, expected, dirs.UserStateDir())
}

func TestLinuxPlatformUserLogDir(t *testing.T) {
	home := setTestHome(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	// Log is a subdirectory of state on XDG platforms
	expected := filepath.Join(home, ".local", "state", "testapp", "log")
	assert.Equal(t, expected, dirs.UserLogDir())
}

func TestLinuxPlatformSystemConfigDir(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	expected := filepath.Join("/etc", "xdg", "testapp")
	assert.Equal(t, expected, dirs.SystemConfigDir())
}

func TestLinuxPlatformSystemDataDir(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	// First dir is /usr/local/share
	expected := filepath.Join("/usr", "local", "share", "testapp")
	assert.Equal(t, expected, dirs.SystemDataDir())
}

func TestLinuxPlatformSystemDataDirs(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	sysDirs := dirs.SystemDataDirs()
	require.Len(t, sysDirs, 2)
	assert.Equal(t, filepath.Join("/usr", "local", "share", "testapp"), sysDirs[0])
	assert.Equal(t, filepath.Join("/usr", "share", "testapp"), sysDirs[1])
}

func TestLinuxPlatformSystemCacheDir(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	expected := filepath.Join("/var", "cache", "testapp")
	assert.Equal(t, expected, dirs.SystemCacheDir())
}

func TestLinuxPlatformSystemStateDir(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	expected := filepath.Join("/var", "lib", "testapp")
	assert.Equal(t, expected, dirs.SystemStateDir())
}

func TestLinuxPlatformSystemLogDir(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	expected := filepath.Join("/var", "log", "testapp")
	assert.Equal(t, expected, dirs.SystemLogDir())
}

func TestLinuxPlatformSystemRuntimeDir(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	expected := filepath.Join("/run", "testapp")
	assert.Equal(t, expected, dirs.SystemRuntimeDir())
}

func TestLinuxPlatformXDGConfigDirsEnv(t *testing.T) {
	testDir := t.TempDir()
	t.Setenv("XDG_CONFIG_DIRS", testDir+":/opt/config")

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	sysDirs := dirs.SystemConfigDirs()
	require.Len(t, sysDirs, 2)
	assert.Equal(t, filepath.Join(testDir, "testapp"), sysDirs[0])
	assert.Equal(t, filepath.Join("/opt", "config", "testapp"), sysDirs[1])
}

func TestLinuxPlatformXDGDataDirsEnv(t *testing.T) {
	testDir := t.TempDir()
	t.Setenv("XDG_DATA_DIRS", testDir+":/opt/data")

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	sysDirs := dirs.SystemDataDirs()
	require.Len(t, sysDirs, 2)
	assert.Equal(t, filepath.Join(testDir, "testapp"), sysDirs[0])
	assert.Equal(t, filepath.Join("/opt", "data", "testapp"), sysDirs[1])
}

func TestLinuxPlatformWithVersion(t *testing.T) {
	home := setTestHome(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Version:  "3.0",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, ".config", "testapp", "3.0")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestLinuxPlatformUserConfigDirsNoFallback(t *testing.T) {
	home := setTestHome(t)

	// On Linux (XDG platform), there are no fallbacks since XDG is native
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	configDirs := dirs.UserConfigDirs()
	require.Len(t, configDirs, 1)

	expected := filepath.Join(home, ".config", "testapp")
	assert.Equal(t, expected, configDirs[0])
}

func TestLinuxPlatformXDGRuntimeDir(t *testing.T) {
	testDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", testDir)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	path, err := dirs.UserRuntimeDir()
	require.NoError(t, err)

	expected := filepath.Join(testDir, "testapp")
	assert.Equal(t, expected, path)
}

func TestLinuxPlatformRuntimeDirFallback(t *testing.T) {
	// When XDG_RUNTIME_DIR is not set, falls back to temp dir
	t.Setenv("XDG_RUNTIME_DIR", "")

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformLinux,
	})
	require.NoError(t, err)

	path, err := dirs.UserRuntimeDir()
	require.NoError(t, err)

	// Should contain app name and uid
	assert.Contains(t, path, "testapp")
}
