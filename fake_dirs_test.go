package toolpaths_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tbhb/toolpaths-go"
)

// testBase returns a platform-appropriate base path for tests.
// On Unix: /base, on Windows: \base.
func testBase() string {
	return filepath.Join(string(filepath.Separator), "base")
}

// p joins path segments using the platform's path separator.
// This is a helper to make cross-platform test assertions cleaner.
func p(parts ...string) string {
	return filepath.Join(parts...)
}

func TestFakeDirsImplementsInterface(_ *testing.T) {
	var _ toolpaths.Dirs = (*toolpaths.FakeDirs)(nil)
}

func TestNewFakeDirs(t *testing.T) {
	base := testBase()
	fake := toolpaths.NewFakeDirs(base)

	assert.Equal(t, p(base, "config"), fake.UserConfigDir())
	assert.Equal(t, p(base, "data"), fake.UserDataDir())
	assert.Equal(t, p(base, "cache"), fake.UserCacheDir())
	assert.Equal(t, p(base, "state"), fake.UserStateDir())
	assert.Equal(t, p(base, "log"), fake.UserLogDir())

	runtime, err := fake.UserRuntimeDir()
	require.NoError(t, err)
	assert.Equal(t, p(base, "runtime"), runtime)

	assert.Equal(t, p(base, "system", "config"), fake.SystemConfigDir())
	assert.Equal(t, p(base, "system", "data"), fake.SystemDataDir())
	assert.Equal(t, p(base, "system", "cache"), fake.SystemCacheDir())
	assert.Equal(t, p(base, "system", "state"), fake.SystemStateDir())
	assert.Equal(t, p(base, "system", "log"), fake.SystemLogDir())
	assert.Equal(t, p(base, "system", "runtime"), fake.SystemRuntimeDir())
}

func TestNewFakeDirsWithTempDir(t *testing.T) {
	fake, cleanup := toolpaths.NewFakeDirsWithTempDir("test")
	defer cleanup()

	// Should have valid paths
	configDir := fake.UserConfigDir()
	assert.NotEmpty(t, configDir)
	assert.True(t, filepath.IsAbs(configDir))

	// cleanup should be safe to call multiple times
	cleanup()
}

func TestFakeDirsUserConfigDirs(t *testing.T) {
	base := testBase()
	fake := toolpaths.NewFakeDirs(base)

	// Default: returns slice with home only
	dirs := fake.UserConfigDirs()
	assert.Equal(t, []string{p(base, "config")}, dirs)

	// With custom dirs
	fake.UserConfigDirsVal = []string{"/custom1", "/custom2"}
	dirs = fake.UserConfigDirs()
	assert.Equal(t, []string{"/custom1", "/custom2"}, dirs)
}

func TestFakeDirsUserDataDirs(t *testing.T) {
	base := testBase()
	fake := toolpaths.NewFakeDirs(base)

	dirs := fake.UserDataDirs()
	assert.Equal(t, []string{p(base, "data")}, dirs)

	fake.UserDataDirsVal = []string{"/custom"}
	dirs = fake.UserDataDirs()
	assert.Equal(t, []string{"/custom"}, dirs)
}

func TestFakeDirsUserRuntimeDirError(t *testing.T) {
	base := testBase()
	fake := toolpaths.NewFakeDirs(base)
	fake.UserRuntimeDirErr = errors.New("runtime dir not available")

	_, err := fake.UserRuntimeDir()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "runtime dir not available")
}

func TestFakeDirsPaths(t *testing.T) {
	base := testBase()
	fake := toolpaths.NewFakeDirs(base)

	assert.Equal(t, p(base, "config", "myfile.yaml"), fake.UserConfigPath("myfile.yaml"))
	assert.Equal(t, p(base, "data", "db.sqlite"), fake.UserDataPath("db.sqlite"))
	assert.Equal(t, p(base, "cache", "tmp.bin"), fake.UserCachePath("tmp.bin"))
	assert.Equal(t, p(base, "state", "state.json"), fake.UserStatePath("state.json"))
	assert.Equal(t, p(base, "log", "app.log"), fake.UserLogPath("app.log"))

	path, err := fake.UserRuntimePath("socket")
	require.NoError(t, err)
	assert.Equal(t, p(base, "runtime", "socket"), path)
}

func TestFakeDirsSystemDirs(t *testing.T) {
	base := testBase()
	fake := toolpaths.NewFakeDirs(base)

	configDirs := fake.SystemConfigDirs()
	assert.Equal(t, []string{p(base, "system", "config")}, configDirs)

	dataDirs := fake.SystemDataDirs()
	assert.Equal(t, []string{p(base, "system", "data")}, dataDirs)
}

func TestFakeDirsSetExisting(t *testing.T) {
	base := testBase()
	fake := toolpaths.NewFakeDirs(base)

	// Initially nothing exists
	_, found := fake.FindConfigFile("config.yaml")
	assert.False(t, found)

	// Mark file as existing
	fake.SetExisting(p(base, "config", "config.yaml"))

	path, found := fake.FindConfigFile("config.yaml")
	assert.True(t, found)
	assert.Equal(t, p(base, "config", "config.yaml"), path)
}

func TestFakeDirsSetNotExisting(t *testing.T) {
	base := testBase()
	fake := toolpaths.NewFakeDirs(base)
	fake.SetExisting(p(base, "config", "config.yaml"))

	// File exists
	_, found := fake.FindConfigFile("config.yaml")
	assert.True(t, found)

	// Mark as not existing
	fake.SetNotExisting(p(base, "config", "config.yaml"))

	_, found = fake.FindConfigFile("config.yaml")
	assert.False(t, found)
}

func TestFakeDirsAllConfigPaths(t *testing.T) {
	base := testBase()
	fake := toolpaths.NewFakeDirs(base)

	paths := fake.AllConfigPaths("config.yaml")
	assert.Contains(t, paths, p(base, "config", "config.yaml"))
	assert.Contains(t, paths, p(base, "system", "config", "config.yaml"))
}

func TestFakeDirsExistingConfigFiles(t *testing.T) {
	base := testBase()
	fake := toolpaths.NewFakeDirs(base)
	fake.SetExisting(p(base, "config", "config.yaml"))
	fake.SetExisting(p(base, "system", "config", "config.yaml"))

	existing := fake.ExistingConfigFiles("config.yaml")
	assert.Len(t, existing, 2)
	assert.Contains(t, existing, p(base, "config", "config.yaml"))
	assert.Contains(t, existing, p(base, "system", "config", "config.yaml"))
}

func TestFakeDirsFindDataFile(t *testing.T) {
	base := testBase()
	fake := toolpaths.NewFakeDirs(base)
	fake.SetExisting(p(base, "data", "db.sqlite"))

	path, found := fake.FindDataFile("db.sqlite")
	assert.True(t, found)
	assert.Equal(t, p(base, "data", "db.sqlite"), path)
}

func TestFakeDirsFindCacheFile(t *testing.T) {
	base := testBase()
	fake := toolpaths.NewFakeDirs(base)
	fake.SetExisting(p(base, "cache", "data.bin"))

	path, found := fake.FindCacheFile("data.bin")
	assert.True(t, found)
	assert.Equal(t, p(base, "cache", "data.bin"), path)
}

func TestFakeDirsFindStateFile(t *testing.T) {
	base := testBase()
	fake := toolpaths.NewFakeDirs(base)
	fake.SetExisting(p(base, "state", "state.json"))

	path, found := fake.FindStateFile("state.json")
	assert.True(t, found)
	assert.Equal(t, p(base, "state", "state.json"), path)
}

func TestFakeDirsFindLogFile(t *testing.T) {
	base := testBase()
	fake := toolpaths.NewFakeDirs(base)
	fake.SetExisting(p(base, "log", "app.log"))

	path, found := fake.FindLogFile("app.log")
	assert.True(t, found)
	assert.Equal(t, p(base, "log", "app.log"), path)
}

func TestFakeDirsFindRuntimeFile(t *testing.T) {
	base := testBase()
	fake := toolpaths.NewFakeDirs(base)
	fake.SetExisting(p(base, "runtime", "socket"))

	path, found := fake.FindRuntimeFile("socket")
	assert.True(t, found)
	assert.Equal(t, p(base, "runtime", "socket"), path)
}

func TestFakeDirsEnsureUserConfigDir(t *testing.T) {
	fake, cleanup := toolpaths.NewFakeDirsWithTempDir("test")
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
	base := testBase()
	fake := toolpaths.NewFakeDirs(base)
	fake.EnsureErrors["config"] = errors.New("permission denied")

	_, err := fake.EnsureUserConfigDir()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
}

func TestFakeDirsEnsureUserDataDir(t *testing.T) {
	fake, cleanup := toolpaths.NewFakeDirsWithTempDir("test")
	defer cleanup()
	fake.CreateDirs = true

	path, err := fake.EnsureUserDataDir()
	require.NoError(t, err)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestFakeDirsEnsureUserCacheDir(t *testing.T) {
	fake, cleanup := toolpaths.NewFakeDirsWithTempDir("test")
	defer cleanup()
	fake.CreateDirs = true

	path, err := fake.EnsureUserCacheDir()
	require.NoError(t, err)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestFakeDirsEnsureUserStateDir(t *testing.T) {
	fake, cleanup := toolpaths.NewFakeDirsWithTempDir("test")
	defer cleanup()
	fake.CreateDirs = true

	path, err := fake.EnsureUserStateDir()
	require.NoError(t, err)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestFakeDirsEnsureUserLogDir(t *testing.T) {
	fake, cleanup := toolpaths.NewFakeDirsWithTempDir("test")
	defer cleanup()
	fake.CreateDirs = true

	path, err := fake.EnsureUserLogDir()
	require.NoError(t, err)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestFakeDirsEmptyValues(t *testing.T) {
	fake := &toolpaths.FakeDirs{}

	// Should not panic with empty values
	assert.Empty(t, fake.UserConfigDir())
	assert.Nil(t, fake.UserConfigDirs())
	assert.Empty(t, fake.SystemConfigDir())
	assert.Nil(t, fake.SystemConfigDirs())
}

func TestFakeDirsSystemRuntimePathEmpty(t *testing.T) {
	fake := &toolpaths.FakeDirs{
		SystemRuntimeDirVal: "",
	}

	// Should return empty when system runtime dir is empty
	assert.Empty(t, fake.SystemRuntimePath("socket"))
}
