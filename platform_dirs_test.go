package tooldirs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tbhb/tooldirs-go"
)

func TestNew(t *testing.T) {
	dirs, err := tooldirs.New("myapp")
	require.NoError(t, err)
	require.NotNil(t, dirs)

	// Should return non-empty config path
	configDir := dirs.UserConfigDir()
	assert.NotEmpty(t, configDir)
	assert.Contains(t, configDir, "myapp")
}

func TestNewWithConfig(t *testing.T) {
	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName: "testapp",
		Version: "1.0",
	})
	require.NoError(t, err)
	require.NotNil(t, dirs)

	// Should include version in path
	configDir := dirs.UserConfigDir()
	assert.Contains(t, configDir, "testapp")
	assert.Contains(t, configDir, "1.0")
}

func TestNewAppNameRequired(t *testing.T) {
	tests := []struct {
		name    string
		appName string
	}{
		{"empty string", ""},
		{"whitespace only", "   "},
		{"tabs only", "\t\t"},
		{"newline only", "\n"},
		{"mixed whitespace", " \t\n "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirs, err := tooldirs.New(tt.appName)
			require.ErrorIs(t, err, tooldirs.ErrAppNameRequired)
			assert.Nil(t, dirs)
		})
	}
}

func TestNewWithConfigAppNameRequired(t *testing.T) {
	tests := []struct {
		name    string
		appName string
	}{
		{"empty string", ""},
		{"whitespace only", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
				AppName: tt.appName,
			})
			require.ErrorIs(t, err, tooldirs.ErrAppNameRequired)
			assert.Nil(t, dirs)
		})
	}
}

func TestUserConfigPath(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	path := dirs.UserConfigPath("config.yaml")
	assert.True(t, filepath.IsAbs(path))
	assert.Contains(t, path, "testapp")
	assert.Contains(t, path, "config.yaml")
}

func TestUserDataPath(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	path := dirs.UserDataPath("data.db")
	assert.True(t, filepath.IsAbs(path))
	assert.Contains(t, path, "testapp")
	assert.Contains(t, path, "data.db")
}

func TestUserCachePath(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	path := dirs.UserCachePath("cache.bin")
	assert.True(t, filepath.IsAbs(path))
	assert.Contains(t, path, "testapp")
	assert.Contains(t, path, "cache.bin")
}

func TestUserStatePath(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	path := dirs.UserStatePath("state.json")
	assert.True(t, filepath.IsAbs(path))
	assert.Contains(t, path, "testapp")
	assert.Contains(t, path, "state.json")
}

func TestUserLogPath(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	path := dirs.UserLogPath("app.log")
	assert.True(t, filepath.IsAbs(path))
	assert.Contains(t, path, "testapp")
	assert.Contains(t, path, "app.log")
}

func TestUserRuntimeDir(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	path, err := dirs.UserRuntimeDir()
	require.NoError(t, err)
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "testapp")
}

func TestUserRuntimePath(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	path, err := dirs.UserRuntimePath("socket")
	require.NoError(t, err)
	assert.Contains(t, path, "testapp")
	assert.Contains(t, path, "socket")
}

func TestSystemConfigDir(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	dir := dirs.SystemConfigDir()
	assert.NotEmpty(t, dir)
	assert.Contains(t, dir, "testapp")
}

func TestSystemConfigDirs(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	sysDirs := dirs.SystemConfigDirs()
	assert.NotEmpty(t, sysDirs)
	// First element should be the primary system config dir
	assert.Equal(t, dirs.SystemConfigDir(), sysDirs[0])
}

func TestSystemDataDir(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	dir := dirs.SystemDataDir()
	assert.NotEmpty(t, dir)
	assert.Contains(t, dir, "testapp")
}

func TestSystemDataDirs(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	sysDirs := dirs.SystemDataDirs()
	assert.NotEmpty(t, sysDirs)
	// First element should be the primary system data dir
	assert.Equal(t, dirs.SystemDataDir(), sysDirs[0])
}

func TestSystemCacheDir(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	dir := dirs.SystemCacheDir()
	assert.NotEmpty(t, dir)
	assert.Contains(t, dir, "testapp")
}

func TestSystemStateDir(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	dir := dirs.SystemStateDir()
	assert.NotEmpty(t, dir)
	assert.Contains(t, dir, "testapp")
}

func TestSystemLogDir(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	dir := dirs.SystemLogDir()
	assert.NotEmpty(t, dir)
	assert.Contains(t, dir, "testapp")
}

func TestEnvOverrides(t *testing.T) {
	t.Run("user config override", func(t *testing.T) {
		envVar := "TEST_MYAPP_CONFIG"
		testDir := t.TempDir()
		t.Setenv(envVar, testDir)

		dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
			AppName: "myapp",
			EnvOverrides: &tooldirs.EnvOverrides{
				AppendAppName: false,
				UserConfig:    envVar,
			},
		})
		require.NoError(t, err)
		assert.Equal(t, testDir, dirs.UserConfigDir())
	})

	t.Run("user config override with app name", func(t *testing.T) {
		envVar := "TEST_MYAPP_CONFIG"
		testDir := t.TempDir()
		t.Setenv(envVar, testDir)

		dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
			AppName: "myapp",
			EnvOverrides: &tooldirs.EnvOverrides{
				AppendAppName: true,
				UserConfig:    envVar,
			},
		})
		require.NoError(t, err)

		expected := filepath.Join(testDir, "myapp")
		assert.Equal(t, expected, dirs.UserConfigDir())
	})

	t.Run("user data override", func(t *testing.T) {
		envVar := "TEST_MYAPP_DATA"
		testDir := t.TempDir()
		t.Setenv(envVar, testDir)

		dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
			AppName: "myapp",
			EnvOverrides: &tooldirs.EnvOverrides{
				AppendAppName: false,
				UserData:      envVar,
			},
		})
		require.NoError(t, err)
		assert.Equal(t, testDir, dirs.UserDataDir())
	})
}

func TestXDGEnvOverrides(t *testing.T) {
	tests := []struct {
		name   string
		envVar string
		getDir func(*tooldirs.PlatformDirs) string
	}{
		{
			name:   "XDG_CONFIG_HOME override",
			envVar: "XDG_CONFIG_HOME",
			getDir: (*tooldirs.PlatformDirs).UserConfigDir,
		},
		{
			name:   "XDG_DATA_HOME override",
			envVar: "XDG_DATA_HOME",
			getDir: (*tooldirs.PlatformDirs).UserDataDir,
		},
		{
			name:   "XDG_CACHE_HOME override",
			envVar: "XDG_CACHE_HOME",
			getDir: (*tooldirs.PlatformDirs).UserCacheDir,
		},
		{
			name:   "XDG_STATE_HOME override",
			envVar: "XDG_STATE_HOME",
			getDir: (*tooldirs.PlatformDirs).UserStateDir,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := t.TempDir()
			t.Setenv(tt.envVar, testDir)

			dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
				AppName:           "myapp",
				XDGOnAllPlatforms: true,
			})
			require.NoError(t, err)

			expected := filepath.Join(testDir, "myapp")
			assert.Equal(t, expected, tt.getDir(dirs))
		})
	}

	t.Run("XDG_RUNTIME_DIR override", func(t *testing.T) {
		testDir := t.TempDir()
		t.Setenv("XDG_RUNTIME_DIR", testDir)

		dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
			AppName:           "myapp",
			XDGOnAllPlatforms: true,
		})
		require.NoError(t, err)

		path, err := dirs.UserRuntimeDir()
		require.NoError(t, err)

		expected := filepath.Join(testDir, "myapp")
		assert.Equal(t, expected, path)
	})
}

func TestUserConfigDirs(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	configDirs := dirs.UserConfigDirs()
	require.NotEmpty(t, configDirs)
	// First element should always be the primary config dir
	assert.Equal(t, dirs.UserConfigDir(), configDirs[0])
}

func TestUserDataDirs(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	dataDirs := dirs.UserDataDirs()
	require.NotEmpty(t, dataDirs)
	assert.Equal(t, dirs.UserDataDir(), dataDirs[0])
}

func TestUserCacheDirs(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	cacheDirs := dirs.UserCacheDirs()
	require.NotEmpty(t, cacheDirs)
	assert.Equal(t, dirs.UserCacheDir(), cacheDirs[0])
}

func TestUserStateDirs(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	stateDirs := dirs.UserStateDirs()
	require.NotEmpty(t, stateDirs)
	assert.Equal(t, dirs.UserStateDir(), stateDirs[0])
}

func TestUserLogDirs(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	logDirs := dirs.UserLogDirs()
	require.NotEmpty(t, logDirs)
	assert.Equal(t, dirs.UserLogDir(), logDirs[0])
}

func TestFindConfigFile(t *testing.T) {
	// Create a temp dir and file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(configFile, []byte("test"), 0o644)
	require.NoError(t, err)

	// Use env override to set config dir
	t.Setenv("TEST_CONFIG", tmpDir)

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName: "testapp",
		EnvOverrides: &tooldirs.EnvOverrides{
			AppendAppName: false,
			UserConfig:    "TEST_CONFIG",
		},
	})
	require.NoError(t, err)

	path, found := dirs.FindConfigFile("config.yaml")
	assert.True(t, found)
	assert.Equal(t, configFile, path)
}

func TestFindConfigFileNotFound(t *testing.T) {
	dirs, err := tooldirs.New("testapp-notexist-12345")
	require.NoError(t, err)

	path, found := dirs.FindConfigFile("nonexistent-file-xyz.yaml")
	assert.False(t, found)
	assert.Empty(t, path)
}

func TestAllConfigPaths(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	paths := dirs.AllConfigPaths("config.yaml")
	require.NotEmpty(t, paths)
	// First path should be user config
	assert.Equal(t, dirs.UserConfigPath("config.yaml"), paths[0])
}

func TestExistingConfigFiles(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(configFile, []byte("test"), 0o644)
	require.NoError(t, err)

	t.Setenv("TEST_CONFIG", tmpDir)

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName: "testapp",
		EnvOverrides: &tooldirs.EnvOverrides{
			AppendAppName: false,
			UserConfig:    "TEST_CONFIG",
		},
	})
	require.NoError(t, err)

	existing := dirs.ExistingConfigFiles("config.yaml")
	assert.Contains(t, existing, configFile)
}

func TestEnsureUserConfigDir(t *testing.T) {
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "newdir")
	t.Setenv("TEST_CONFIG", testDir)

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName: "testapp",
		EnvOverrides: &tooldirs.EnvOverrides{
			AppendAppName: false,
			UserConfig:    "TEST_CONFIG",
		},
	})
	require.NoError(t, err)

	path, err := dirs.EnsureUserConfigDir()
	require.NoError(t, err)
	assert.Equal(t, testDir, path)

	// Directory should exist
	info, err := os.Stat(testDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestEnsureUserDataDir(t *testing.T) {
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "newdir")
	t.Setenv("TEST_DATA", testDir)

	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName: "testapp",
		EnvOverrides: &tooldirs.EnvOverrides{
			AppendAppName: false,
			UserData:      "TEST_DATA",
		},
	})
	require.NoError(t, err)

	path, err := dirs.EnsureUserDataDir()
	require.NoError(t, err)
	assert.Equal(t, testDir, path)

	info, err := os.Stat(testDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestPlatformDirsString(t *testing.T) {
	dirs, err := tooldirs.New("testapp")
	require.NoError(t, err)

	str := dirs.String()

	assert.Contains(t, str, "testapp")
	assert.Contains(t, str, "User directories")
	assert.Contains(t, str, "System directories")
	assert.Contains(t, str, "Config:")
	assert.Contains(t, str, "Data:")
	assert.Contains(t, str, "Cache:")
}

func TestPlatformDirsStringWithVersion(t *testing.T) {
	dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
		AppName: "testapp",
		Version: "2.0",
	})
	require.NoError(t, err)

	str := dirs.String()

	assert.Contains(t, str, "testapp")
	assert.Contains(t, str, "2.0")
}

func TestPlatformString(t *testing.T) {
	tests := []struct {
		platform tooldirs.Platform
		expected string
	}{
		{tooldirs.PlatformAuto, "auto"},
		{tooldirs.PlatformLinux, "linux"},
		{tooldirs.PlatformMacOS, "macos"},
		{tooldirs.PlatformWindows, "windows"},
		{tooldirs.PlatformFreeBSD, "freebsd"},
		{tooldirs.PlatformOpenBSD, "openbsd"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.platform.String())
		})
	}
}

func TestPlatformOverride(t *testing.T) {
	// This test verifies that explicitly setting Platform to a specific value
	// (including PlatformLinux) is respected, rather than triggering auto-detection.

	t.Run("explicit PlatformLinux is respected", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)

		dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
			AppName:  "testapp",
			Platform: tooldirs.PlatformLinux,
		})
		require.NoError(t, err)

		// Linux uses XDG paths by default (~/.config/appname)
		configDir := dirs.UserConfigDir()
		expected := filepath.Join(home, ".config", "testapp")
		assert.Equal(t, expected, configDir)
	})

	t.Run("explicit PlatformMacOS is respected", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)

		dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
			AppName:  "testapp",
			Platform: tooldirs.PlatformMacOS,
		})
		require.NoError(t, err)

		// macOS uses Library paths (~/Library/Application Support/appname)
		configDir := dirs.UserConfigDir()
		expected := filepath.Join(home, "Library", "Application Support", "testapp")
		assert.Equal(t, expected, configDir)
	})

	t.Run("PlatformAuto triggers detection", func(t *testing.T) {
		// PlatformAuto (zero value) should trigger auto-detection
		dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
			AppName:  "testapp",
			Platform: tooldirs.PlatformAuto,
		})
		require.NoError(t, err)

		// Just verify it doesn't panic and returns a non-empty path
		configDir := dirs.UserConfigDir()
		assert.NotEmpty(t, configDir)
	})

	t.Run("zero value Config uses auto-detection", func(t *testing.T) {
		// Omitting Platform field should default to PlatformAuto
		dirs, err := tooldirs.NewWithConfig(tooldirs.Config{
			AppName: "testapp",
		})
		require.NoError(t, err)

		configDir := dirs.UserConfigDir()
		assert.NotEmpty(t, configDir)
	})
}
