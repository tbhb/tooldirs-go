package toolpaths_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tbhb/toolpaths-go"
)

// Tests for Windows path resolution logic.
// These use explicit Platform: PlatformWindows and mock HOME,
// so they can run on any platform.

// setTestHomeWindows sets up a test home directory and returns a cleanup function.
func setTestHomeWindows(t *testing.T) string {
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

func TestWindowsPlatformUserConfigDir(t *testing.T) {
	home := setTestHomeWindows(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformWindows,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "AppData", "Local", "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestWindowsPlatformUserConfigDirRoaming(t *testing.T) {
	home := setTestHomeWindows(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformWindows,
		Roaming:  true,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "AppData", "Roaming", "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestWindowsPlatformUserDataDir(t *testing.T) {
	home := setTestHomeWindows(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformWindows,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "AppData", "Local", "testapp")
	assert.Equal(t, expected, dirs.UserDataDir())
}

func TestWindowsPlatformUserCacheDir(t *testing.T) {
	home := setTestHomeWindows(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformWindows,
	})
	require.NoError(t, err)

	// Cache is always under LocalAppData
	expected := filepath.Join(home, "AppData", "Local", "testapp", "cache")
	assert.Equal(t, expected, dirs.UserCacheDir())
}

func TestWindowsPlatformUserLogDir(t *testing.T) {
	home := setTestHomeWindows(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformWindows,
	})
	require.NoError(t, err)

	// Log is always under LocalAppData
	expected := filepath.Join(home, "AppData", "Local", "testapp", "log")
	assert.Equal(t, expected, dirs.UserLogDir())
}

func TestWindowsPlatformUserStateDir(t *testing.T) {
	home := setTestHomeWindows(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformWindows,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "AppData", "Local", "testapp")
	assert.Equal(t, expected, dirs.UserStateDir())
}

func TestWindowsPlatformWithAppAuthor(t *testing.T) {
	home := setTestHomeWindows(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:   "testapp",
		AppAuthor: "MyCompany",
		Platform:  toolpaths.PlatformWindows,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "AppData", "Local", "MyCompany", "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestWindowsPlatformWithVersion(t *testing.T) {
	home := setTestHomeWindows(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Version:  "2.0",
		Platform: toolpaths.PlatformWindows,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "AppData", "Local", "testapp", "2.0")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestWindowsPlatformWithAppAuthorAndVersion(t *testing.T) {
	home := setTestHomeWindows(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:   "testapp",
		AppAuthor: "MyCompany",
		Version:   "2.0",
		Platform:  toolpaths.PlatformWindows,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "AppData", "Local", "MyCompany", "testapp", "2.0")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestWindowsPlatformSystemConfigDir(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformWindows,
	})
	require.NoError(t, err)

	// On actual Windows, path includes drive letter (C:\ProgramData\testapp)
	// On other platforms, we get \ProgramData\testapp
	result := dirs.SystemConfigDir()
	expectedSuffix := filepath.Join("ProgramData", "testapp")
	assert.True(t, strings.HasSuffix(result, expectedSuffix),
		"expected path to end with %q, got %q", expectedSuffix, result)
}

func TestWindowsPlatformSystemDataDir(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformWindows,
	})
	require.NoError(t, err)

	// On actual Windows, path includes drive letter
	result := dirs.SystemDataDir()
	expectedSuffix := filepath.Join("ProgramData", "testapp")
	assert.True(t, strings.HasSuffix(result, expectedSuffix),
		"expected path to end with %q, got %q", expectedSuffix, result)
}

func TestWindowsPlatformSystemCacheDir(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformWindows,
	})
	require.NoError(t, err)

	// On actual Windows, path includes drive letter
	result := dirs.SystemCacheDir()
	expectedSuffix := filepath.Join("ProgramData", "testapp", "cache")
	assert.True(t, strings.HasSuffix(result, expectedSuffix),
		"expected path to end with %q, got %q", expectedSuffix, result)
}

func TestWindowsPlatformSystemLogDir(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformWindows,
	})
	require.NoError(t, err)

	// On actual Windows, path includes drive letter
	result := dirs.SystemLogDir()
	expectedSuffix := filepath.Join("ProgramData", "testapp", "log")
	assert.True(t, strings.HasSuffix(result, expectedSuffix),
		"expected path to end with %q, got %q", expectedSuffix, result)
}

func TestWindowsPlatformSystemRuntimeDir(t *testing.T) {
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformWindows,
	})
	require.NoError(t, err)

	// Windows has no system runtime dir
	assert.Empty(t, dirs.SystemRuntimeDir())
}

func TestWindowsPlatformXDGOnAllPlatforms(t *testing.T) {
	home := setTestHomeWindows(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:           "testapp",
		Platform:          toolpaths.PlatformWindows,
		XDGOnAllPlatforms: true,
	})
	require.NoError(t, err)

	// With XDGOnAllPlatforms, should use XDG paths
	expected := filepath.Join(home, ".config", "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestWindowsPlatformIncludeXDGFallbacks(t *testing.T) {
	home := setTestHomeWindows(t)

	// Default: IncludeXDGFallbacks is true
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformWindows,
	})
	require.NoError(t, err)

	configDirs := dirs.UserConfigDirs()
	require.Len(t, configDirs, 2)

	nativePath := filepath.Join(home, "AppData", "Local", "testapp")
	xdgPath := filepath.Join(home, ".config", "testapp")

	assert.Equal(t, nativePath, configDirs[0])
	assert.Equal(t, xdgPath, configDirs[1])
}

func TestWindowsPlatformIncludeXDGFallbacksDisabled(t *testing.T) {
	home := setTestHomeWindows(t)

	falseVal := false
	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:             "testapp",
		Platform:            toolpaths.PlatformWindows,
		IncludeXDGFallbacks: &falseVal,
	})
	require.NoError(t, err)

	configDirs := dirs.UserConfigDirs()
	require.Len(t, configDirs, 1)

	nativePath := filepath.Join(home, "AppData", "Local", "testapp")
	assert.Equal(t, nativePath, configDirs[0])
}

func TestWindowsPlatformUserRuntimeDir(t *testing.T) {
	home := setTestHomeWindows(t)

	dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
		AppName:  "testapp",
		Platform: toolpaths.PlatformWindows,
	})
	require.NoError(t, err)

	path, err := dirs.UserRuntimeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, "AppData", "Local", "testapp", "runtime")
	assert.Equal(t, expected, path)
}
