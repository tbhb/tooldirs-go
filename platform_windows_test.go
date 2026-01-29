package tooldirs_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tbhb/tooldirs-go"
)

// Tests for Windows path resolution logic.
// These use explicit Platform: PlatformWindows and mock HOME,
// so they can run on any platform.

func TestWindowsPlatformUserConfigDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:  "testapp",
		Platform: tooldirs.PlatformWindows,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "AppData", "Local", "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestWindowsPlatformUserConfigDirRoaming(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:  "testapp",
		Platform: tooldirs.PlatformWindows,
		Roaming:  true,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "AppData", "Roaming", "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestWindowsPlatformUserDataDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:  "testapp",
		Platform: tooldirs.PlatformWindows,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "AppData", "Local", "testapp")
	assert.Equal(t, expected, dirs.UserDataDir())
}

func TestWindowsPlatformUserCacheDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:  "testapp",
		Platform: tooldirs.PlatformWindows,
	})
	require.NoError(t, err)

	// Cache is always under LocalAppData
	expected := filepath.Join(home, "AppData", "Local", "testapp", "cache")
	assert.Equal(t, expected, dirs.UserCacheDir())
}

func TestWindowsPlatformUserLogDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:  "testapp",
		Platform: tooldirs.PlatformWindows,
	})
	require.NoError(t, err)

	// Log is always under LocalAppData
	expected := filepath.Join(home, "AppData", "Local", "testapp", "log")
	assert.Equal(t, expected, dirs.UserLogDir())
}

func TestWindowsPlatformUserStateDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:  "testapp",
		Platform: tooldirs.PlatformWindows,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "AppData", "Local", "testapp")
	assert.Equal(t, expected, dirs.UserStateDir())
}

func TestWindowsPlatformWithAppAuthor(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:   "testapp",
		AppAuthor: "MyCompany",
		Platform:  tooldirs.PlatformWindows,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "AppData", "Local", "MyCompany", "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestWindowsPlatformWithVersion(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:  "testapp",
		Version:  "2.0",
		Platform: tooldirs.PlatformWindows,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "AppData", "Local", "testapp", "2.0")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestWindowsPlatformWithAppAuthorAndVersion(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:   "testapp",
		AppAuthor: "MyCompany",
		Version:   "2.0",
		Platform:  tooldirs.PlatformWindows,
	})
	require.NoError(t, err)

	expected := filepath.Join(home, "AppData", "Local", "MyCompany", "testapp", "2.0")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestWindowsPlatformSystemConfigDir(t *testing.T) {
	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:  "testapp",
		Platform: tooldirs.PlatformWindows,
	})
	require.NoError(t, err)

	expected := filepath.Join(string(filepath.Separator), "ProgramData", "testapp")
	assert.Equal(t, expected, dirs.SystemConfigDir())
}

func TestWindowsPlatformSystemDataDir(t *testing.T) {
	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:  "testapp",
		Platform: tooldirs.PlatformWindows,
	})
	require.NoError(t, err)

	expected := filepath.Join(string(filepath.Separator), "ProgramData", "testapp")
	assert.Equal(t, expected, dirs.SystemDataDir())
}

func TestWindowsPlatformSystemCacheDir(t *testing.T) {
	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:  "testapp",
		Platform: tooldirs.PlatformWindows,
	})
	require.NoError(t, err)

	expected := filepath.Join(string(filepath.Separator), "ProgramData", "testapp", "cache")
	assert.Equal(t, expected, dirs.SystemCacheDir())
}

func TestWindowsPlatformSystemLogDir(t *testing.T) {
	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:  "testapp",
		Platform: tooldirs.PlatformWindows,
	})
	require.NoError(t, err)

	expected := filepath.Join(string(filepath.Separator), "ProgramData", "testapp", "log")
	assert.Equal(t, expected, dirs.SystemLogDir())
}

func TestWindowsPlatformSystemRuntimeDir(t *testing.T) {
	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:  "testapp",
		Platform: tooldirs.PlatformWindows,
	})
	require.NoError(t, err)

	// Windows has no system runtime dir
	assert.Empty(t, dirs.SystemRuntimeDir())
}

func TestWindowsPlatformXDGOnAllPlatforms(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:           "testapp",
		Platform:          tooldirs.PlatformWindows,
		XDGOnAllPlatforms: true,
	})
	require.NoError(t, err)

	// With XDGOnAllPlatforms, should use XDG paths
	expected := filepath.Join(home, ".config", "testapp")
	assert.Equal(t, expected, dirs.UserConfigDir())
}

func TestWindowsPlatformIncludeXDGFallbacks(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	// Default: IncludeXDGFallbacks is true
	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:  "testapp",
		Platform: tooldirs.PlatformWindows,
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
	home := t.TempDir()
	t.Setenv("HOME", home)

	falseVal := false
	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:             "testapp",
		Platform:            tooldirs.PlatformWindows,
		IncludeXDGFallbacks: &falseVal,
	})
	require.NoError(t, err)

	configDirs := dirs.UserConfigDirs()
	require.Len(t, configDirs, 1)

	nativePath := filepath.Join(home, "AppData", "Local", "testapp")
	assert.Equal(t, nativePath, configDirs[0])
}

func TestWindowsPlatformUserRuntimeDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName:  "testapp",
		Platform: tooldirs.PlatformWindows,
	})
	require.NoError(t, err)

	path, err := dirs.UserRuntimeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, "AppData", "Local", "testapp", "runtime")
	assert.Equal(t, expected, path)
}
