package tooldirs_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tbhb/tooldirs-go"
)

func TestFakeDirsImplementsInterface(_ *testing.T) {
	var _ tooldirs.Dirs = (*tooldirs.FakeDirs)(nil)
}

func TestNewFakeDirs(t *testing.T) {
	fake := tooldirs.NewFakeDirs("/base")

	assert.Equal(t, "/base/config", fake.UserConfigDir())
	assert.Equal(t, "/base/data", fake.UserDataDir())
	assert.Equal(t, "/base/cache", fake.UserCacheDir())
	assert.Equal(t, "/base/state", fake.UserStateDir())
	assert.Equal(t, "/base/log", fake.UserLogDir())

	runtime, err := fake.UserRuntimeDir()
	require.NoError(t, err)
	assert.Equal(t, "/base/runtime", runtime)

	assert.Equal(t, "/base/system/config", fake.SystemConfigDir())
	assert.Equal(t, "/base/system/data", fake.SystemDataDir())
	assert.Equal(t, "/base/system/cache", fake.SystemCacheDir())
	assert.Equal(t, "/base/system/state", fake.SystemStateDir())
	assert.Equal(t, "/base/system/log", fake.SystemLogDir())
	assert.Equal(t, "/base/system/runtime", fake.SystemRuntimeDir())
}

func TestNewFakeDirsWithTempDir(t *testing.T) {
	fake, cleanup := tooldirs.NewFakeDirsWithTempDir("test")
	defer cleanup()

	// Should have valid paths
	configDir := fake.UserConfigDir()
	assert.NotEmpty(t, configDir)
	assert.True(t, filepath.IsAbs(configDir))

	// cleanup should be safe to call multiple times
	cleanup()
}

func TestFakeDirsUserConfigDirs(t *testing.T) {
	fake := tooldirs.NewFakeDirs("/base")

	// Default: returns slice with home only
	dirs := fake.UserConfigDirs()
	assert.Equal(t, []string{"/base/config"}, dirs)

	// With custom dirs
	fake.UserConfigDirsVal = []string{"/custom1", "/custom2"}
	dirs = fake.UserConfigDirs()
	assert.Equal(t, []string{"/custom1", "/custom2"}, dirs)
}

func TestFakeDirsUserDataDirs(t *testing.T) {
	fake := tooldirs.NewFakeDirs("/base")

	dirs := fake.UserDataDirs()
	assert.Equal(t, []string{"/base/data"}, dirs)

	fake.UserDataDirsVal = []string{"/custom"}
	dirs = fake.UserDataDirs()
	assert.Equal(t, []string{"/custom"}, dirs)
}

func TestFakeDirsUserRuntimeDirError(t *testing.T) {
	fake := tooldirs.NewFakeDirs("/base")
	fake.UserRuntimeDirErr = errors.New("runtime dir not available")

	_, err := fake.UserRuntimeDir()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "runtime dir not available")
}

func TestFakeDirsPaths(t *testing.T) {
	fake := tooldirs.NewFakeDirs("/base")

	assert.Equal(t, "/base/config/myfile.yaml", fake.UserConfigPath("myfile.yaml"))
	assert.Equal(t, "/base/data/db.sqlite", fake.UserDataPath("db.sqlite"))
	assert.Equal(t, "/base/cache/tmp.bin", fake.UserCachePath("tmp.bin"))
	assert.Equal(t, "/base/state/state.json", fake.UserStatePath("state.json"))
	assert.Equal(t, "/base/log/app.log", fake.UserLogPath("app.log"))

	path, err := fake.UserRuntimePath("socket")
	require.NoError(t, err)
	assert.Equal(t, "/base/runtime/socket", path)
}

func TestFakeDirsSystemDirs(t *testing.T) {
	fake := tooldirs.NewFakeDirs("/base")

	configDirs := fake.SystemConfigDirs()
	assert.Equal(t, []string{"/base/system/config"}, configDirs)

	dataDirs := fake.SystemDataDirs()
	assert.Equal(t, []string{"/base/system/data"}, dataDirs)
}

func TestFakeDirsSetExisting(t *testing.T) {
	fake := tooldirs.NewFakeDirs("/base")

	// Initially nothing exists
	_, found := fake.FindConfigFile("config.yaml")
	assert.False(t, found)

	// Mark file as existing
	fake.SetExisting("/base/config/config.yaml")

	path, found := fake.FindConfigFile("config.yaml")
	assert.True(t, found)
	assert.Equal(t, "/base/config/config.yaml", path)
}

func TestFakeDirsSetNotExisting(t *testing.T) {
	fake := tooldirs.NewFakeDirs("/base")
	fake.SetExisting("/base/config/config.yaml")

	// File exists
	_, found := fake.FindConfigFile("config.yaml")
	assert.True(t, found)

	// Mark as not existing
	fake.SetNotExisting("/base/config/config.yaml")

	_, found = fake.FindConfigFile("config.yaml")
	assert.False(t, found)
}

func TestFakeDirsAllConfigPaths(t *testing.T) {
	fake := tooldirs.NewFakeDirs("/base")

	paths := fake.AllConfigPaths("config.yaml")
	assert.Contains(t, paths, "/base/config/config.yaml")
	assert.Contains(t, paths, "/base/system/config/config.yaml")
}

func TestFakeDirsExistingConfigFiles(t *testing.T) {
	fake := tooldirs.NewFakeDirs("/base")
	fake.SetExisting("/base/config/config.yaml")
	fake.SetExisting("/base/system/config/config.yaml")

	existing := fake.ExistingConfigFiles("config.yaml")
	assert.Len(t, existing, 2)
	assert.Contains(t, existing, "/base/config/config.yaml")
	assert.Contains(t, existing, "/base/system/config/config.yaml")
}

func TestFakeDirsFindDataFile(t *testing.T) {
	fake := tooldirs.NewFakeDirs("/base")
	fake.SetExisting("/base/data/db.sqlite")

	path, found := fake.FindDataFile("db.sqlite")
	assert.True(t, found)
	assert.Equal(t, "/base/data/db.sqlite", path)
}

func TestFakeDirsFindCacheFile(t *testing.T) {
	fake := tooldirs.NewFakeDirs("/base")
	fake.SetExisting("/base/cache/data.bin")

	path, found := fake.FindCacheFile("data.bin")
	assert.True(t, found)
	assert.Equal(t, "/base/cache/data.bin", path)
}

func TestFakeDirsFindStateFile(t *testing.T) {
	fake := tooldirs.NewFakeDirs("/base")
	fake.SetExisting("/base/state/state.json")

	path, found := fake.FindStateFile("state.json")
	assert.True(t, found)
	assert.Equal(t, "/base/state/state.json", path)
}

func TestFakeDirsFindLogFile(t *testing.T) {
	fake := tooldirs.NewFakeDirs("/base")
	fake.SetExisting("/base/log/app.log")

	path, found := fake.FindLogFile("app.log")
	assert.True(t, found)
	assert.Equal(t, "/base/log/app.log", path)
}

func TestFakeDirsFindRuntimeFile(t *testing.T) {
	fake := tooldirs.NewFakeDirs("/base")
	fake.SetExisting("/base/runtime/socket")

	path, found := fake.FindRuntimeFile("socket")
	assert.True(t, found)
	assert.Equal(t, "/base/runtime/socket", path)
}

func TestFakeDirsEnsureUserConfigDir(t *testing.T) {
	fake, cleanup := tooldirs.NewFakeDirsWithTempDir("test")
	defer cleanup()
	fake.CreateDirs = true

	path, err := fake.EnsureUserConfigDir()
	require.NoError(t, err)

	// Directory should exist
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestFakeDirsEnsureUserConfigDirError(t *testing.T) {
	fake := tooldirs.NewFakeDirs("/base")
	fake.EnsureErrors["config"] = errors.New("permission denied")

	_, err := fake.EnsureUserConfigDir()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
}

func TestFakeDirsEnsureUserDataDir(t *testing.T) {
	fake, cleanup := tooldirs.NewFakeDirsWithTempDir("test")
	defer cleanup()
	fake.CreateDirs = true

	path, err := fake.EnsureUserDataDir()
	require.NoError(t, err)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestFakeDirsEnsureUserCacheDir(t *testing.T) {
	fake, cleanup := tooldirs.NewFakeDirsWithTempDir("test")
	defer cleanup()
	fake.CreateDirs = true

	path, err := fake.EnsureUserCacheDir()
	require.NoError(t, err)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestFakeDirsEnsureUserStateDir(t *testing.T) {
	fake, cleanup := tooldirs.NewFakeDirsWithTempDir("test")
	defer cleanup()
	fake.CreateDirs = true

	path, err := fake.EnsureUserStateDir()
	require.NoError(t, err)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestFakeDirsEnsureUserLogDir(t *testing.T) {
	fake, cleanup := tooldirs.NewFakeDirsWithTempDir("test")
	defer cleanup()
	fake.CreateDirs = true

	path, err := fake.EnsureUserLogDir()
	require.NoError(t, err)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestFakeDirsEmptyValues(t *testing.T) {
	fake := &tooldirs.FakeDirs{}

	// Should not panic with empty values
	assert.Empty(t, fake.UserConfigDir())
	assert.Nil(t, fake.UserConfigDirs())
	assert.Empty(t, fake.SystemConfigDir())
	assert.Nil(t, fake.SystemConfigDirs())
}

func TestFakeDirsSystemRuntimePathEmpty(t *testing.T) {
	fake := &tooldirs.FakeDirs{
		SystemRuntimeDirVal: "",
	}

	// Should return empty when system runtime dir is empty
	assert.Empty(t, fake.SystemRuntimePath("socket"))
}
