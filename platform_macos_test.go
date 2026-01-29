package toolpaths_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tbhb/toolpaths-go"
)

// Tests for macOS path resolution logic.
// These use explicit Platform: PlatformMacOS and mock HOME,
// so they can run on any platform.

// setTestHomeMacOS sets up a test home directory and returns a cleanup function.
func setTestHomeMacOS(t *testing.T) string {
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

func TestMacOSPlatformUserConfigDir(t *testing.T) {
	home := setTestHomeMacOS(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformMacOS,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "Library", "Application Support", "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestMacOSPlatformUserDataDir(t *testing.T) {
	home := setTestHomeMacOS(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformMacOS,
	})
	require.NoError(t, err)

	// On macOS, data dir is same as config dir
	expected := filepath.Join(home, "Library", "Application Support", "testapp")
	assert.Equal(t, expected, dirs.UserDataDir())
}

func TestMacOSPlatformUserCacheDir(t *testing.T) {
	home := setTestHomeMacOS(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformMacOS,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "Library", "Caches", "testapp")
	assert.Equal(t, expected, dirs.UserCacheDir())
}

func TestMacOSPlatformUserLogDir(t *testing.T) {
	home := setTestHomeMacOS(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformMacOS,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "Library", "Logs", "testapp")
	assert.Equal(t, expected, dirs.UserLogDir())
}

func TestMacOSPlatformUserStateDir(t *testing.T) {
	home := setTestHomeMacOS(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformMacOS,
	})
	require.NoError(t, err)

	// On macOS, state dir is same as config dir
	expected := filepath.Join(home, "Library", "Application Support", "testapp")
	assert.Equal(t, expected, dirs.UserStateDir())
}

func TestMacOSPlatformSystemConfigDir(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformMacOS,
	})
	require.NoError(t, err)

	expected := filepath.Join("/Library", "Application Support", "testapp")
	assert.Equal(t, expected, dirs.SystemConfigDir())
}

func TestMacOSPlatformSystemDataDir(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformMacOS,
	})
	require.NoError(t, err)

	expected := filepath.Join("/Library", "Application Support", "testapp")
	assert.Equal(t, expected, dirs.SystemDataDir())
}

func TestMacOSPlatformSystemCacheDir(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformMacOS,
	})
	require.NoError(t, err)

	expected := filepath.Join("/Library", "Caches", "testapp")
	assert.Equal(t, expected, dirs.SystemCacheDir())
}

func TestMacOSPlatformSystemLogDir(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformMacOS,
	})
	require.NoError(t, err)

	expected := filepath.Join("/Library", "Logs", "testapp")
	assert.Equal(t, expected, dirs.SystemLogDir())
}

func TestMacOSPlatformSystemRuntimeDir(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformMacOS,
	})
	require.NoError(t, err)

	// macOS has no system runtime dir
	assert.Empty(t, dirs.SystemRuntimeDir())
}

func TestMacOSPlatformXDGOnAllPlatforms(t *testing.T) {
	home := setTestHomeMacOS(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:           "testapp",
		Platform:          toolpaths.PlatformMacOS,
		XDGOnAllPlatforms: true,
	})
	require.NoError(t, err)

	// With XDGOnAllPlatforms, should use XDG paths
	expected := filepath.Join(home, ".config", "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestMacOSPlatformIncludeXDGFallbacks(t *testing.T) {
	home := setTestHomeMacOS(t)

	// Default: IncludeXDGFallbacks is true
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformMacOS,
	})
	require.NoError(t, err)

	configDirs := dirs.UserConfigDirs()
	require.Len(t, configDirs, 2)

	nativePath := filepath.Join(home, "Library", "Application Support", "testapp")
	xdgPath := filepath.Join(home, ".config", "testapp")

	assert.Equal(t, nativePath, configDirs[0])
	assert.Equal(t, xdgPath, configDirs[1])
}

func TestMacOSPlatformIncludeXDGFallbacksDisabled(t *testing.T) {
	home := setTestHomeMacOS(t)

	falseVal := false
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:             "testapp",
		Platform:            toolpaths.PlatformMacOS,
		IncludeXDGFallbacks: &falseVal,
	})
	require.NoError(t, err)

	configDirs := dirs.UserConfigDirs()
	require.Len(t, configDirs, 1)

	nativePath := filepath.Join(home, "Library", "Application Support", "testapp")
	assert.Equal(t, nativePath, configDirs[0])
}

func TestMacOSPlatformWithVersion(t *testing.T) {
	home := setTestHomeMacOS(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Version:  "2.0",
		Platform: toolpaths.PlatformMacOS,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "Library", "Application Support", "testapp", "2.0")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestMacOSPlatformXDGEnvRespected(t *testing.T) {
	home := setTestHomeMacOS(t)

	testDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", testDir)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformMacOS,
	})
	require.NoError(t, err)

	// XDG env var should be respected even on macOS
	expected := filepath.Join(testDir, "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())

	_ = home // prevent unused variable warning
}
