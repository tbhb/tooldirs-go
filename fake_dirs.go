package tooldirs

import (
	"os"
	"path/filepath"
)

// FakeDirs is a test double for Dirs that returns configurable paths.
// All fields are exported for direct manipulation in tests.
//
// Example usage:
//
//	fake := &tooldirs.FakeDirs{
//	    UserConfigHomeVal: "/tmp/test/config",
//	    UserDataHomeVal:   "/tmp/test/data",
//	}
//	// Use fake wherever PlatformDirs interface is expected
//	app := NewApp(fake)
type FakeDirs struct {
	// User directory homes (write targets)
	UserConfigHomeVal string
	UserDataHomeVal   string
	UserCacheHomeVal  string
	UserStateHomeVal  string
	UserLogHomeVal    string
	UserRuntimeDirVal string
	UserRuntimeDirErr error

	// User directory search paths (if nil, defaults to []string{*HomeVal})
	UserConfigDirsVal []string
	UserDataDirsVal   []string
	UserCacheDirsVal  []string
	UserStateDirsVal  []string
	UserLogDirsVal    []string

	// System directories
	SystemConfigDirsVal []string
	SystemDataDirsVal   []string
	SystemCacheDirVal   string
	SystemStateDirVal   string
	SystemLogDirVal     string
	SystemRuntimeDirVal string

	// ExistingFiles maps paths to existence. Used by Find* and Existing* methods.
	// If nil, file existence checks use the real filesystem.
	// If non-nil, only paths in this map with true values are considered to exist.
	ExistingFiles map[string]bool

	// EnsureErrors maps directory types to errors returned by Ensure* methods.
	// Keys are: "config", "data", "cache", "state", "log"
	EnsureErrors map[string]error

	// CreateDirs controls whether Ensure* methods actually create directories.
	// If false (default), Ensure* methods just return the path (and any configured error).
	// If true, Ensure* methods call os.MkdirAll.
	CreateDirs bool
}

// Compile-time check that FakeDirs implements Dirs.
var _ Dirs = (*FakeDirs)(nil)

// NewFakeDirs creates a FakeDirs with all paths set to subdirectories of the given base.
// This is convenient for tests that need a complete, consistent fake.
//
// Example:
//
//	fake := tooldirs.NewFakeDirs("/tmp/test-app")
//	// fake.UserConfigHomeVal == "/tmp/test-app/config"
//	// fake.UserDataHomeVal == "/tmp/test-app/data"
//	// etc.
func NewFakeDirs(base string) *FakeDirs {
	return &FakeDirs{
		UserConfigHomeVal:   filepath.Join(base, "config"),
		UserDataHomeVal:     filepath.Join(base, "data"),
		UserCacheHomeVal:    filepath.Join(base, "cache"),
		UserStateHomeVal:    filepath.Join(base, "state"),
		UserLogHomeVal:      filepath.Join(base, "log"),
		UserRuntimeDirVal:   filepath.Join(base, "runtime"),
		SystemConfigDirsVal: []string{filepath.Join(base, "system", "config")},
		SystemDataDirsVal:   []string{filepath.Join(base, "system", "data")},
		SystemCacheDirVal:   filepath.Join(base, "system", "cache"),
		SystemStateDirVal:   filepath.Join(base, "system", "state"),
		SystemLogDirVal:     filepath.Join(base, "system", "log"),
		SystemRuntimeDirVal: filepath.Join(base, "system", "runtime"),
		ExistingFiles:       make(map[string]bool),
		EnsureErrors:        make(map[string]error),
	}
}

// NewFakeDirsWithTempDir creates a FakeDirs rooted in a new temporary directory.
// Returns the FakeDirs and a cleanup function that removes the temp directory.
// The cleanup function is safe to call multiple times.
//
// Example:
//
//	fake, cleanup := tooldirs.NewFakeDirsWithTempDir(t.Name())
//	defer cleanup()
//	fake.CreateDirs = true  // Actually create directories
func NewFakeDirsWithTempDir(prefix string) (*FakeDirs, func()) {
	base, err := os.MkdirTemp("", prefix)
	if err != nil {
		// Fall back to a path that won't exist
		base = filepath.Join(os.TempDir(), "tooldirs-fake-"+prefix)
	}
	fake := NewFakeDirs(base)
	cleanup := func() {
		_ = os.RemoveAll(base)
	}
	return fake, cleanup
}

// SetExisting marks a path as existing for Find* and Existing* methods.
func (f *FakeDirs) SetExisting(path string) {
	if f.ExistingFiles == nil {
		f.ExistingFiles = make(map[string]bool)
	}
	f.ExistingFiles[path] = true
}

// SetNotExisting marks a path as not existing.
func (f *FakeDirs) SetNotExisting(path string) {
	if f.ExistingFiles == nil {
		f.ExistingFiles = make(map[string]bool)
	}
	f.ExistingFiles[path] = false
}

// fileExists checks if a path exists, using ExistingFiles map if set.
func (f *FakeDirs) fileExists(path string) bool {
	if f.ExistingFiles != nil {
		return f.ExistingFiles[path]
	}
	_, err := os.Stat(path)
	return err == nil
}

// --- User config ---

func (f *FakeDirs) UserConfigDir() string {
	return f.UserConfigHomeVal
}

func (f *FakeDirs) UserConfigDirs() []string {
	if f.UserConfigDirsVal != nil {
		return f.UserConfigDirsVal
	}
	if f.UserConfigHomeVal != "" {
		return []string{f.UserConfigHomeVal}
	}
	return nil
}

func (f *FakeDirs) UserConfigPath(elem ...string) string {
	return path(f.UserConfigHomeVal, elem...)
}

// --- User data ---

func (f *FakeDirs) UserDataDir() string {
	return f.UserDataHomeVal
}

func (f *FakeDirs) UserDataDirs() []string {
	if f.UserDataDirsVal != nil {
		return f.UserDataDirsVal
	}
	if f.UserDataHomeVal != "" {
		return []string{f.UserDataHomeVal}
	}
	return nil
}

func (f *FakeDirs) UserDataPath(elem ...string) string {
	return path(f.UserDataHomeVal, elem...)
}

// --- User cache ---

func (f *FakeDirs) UserCacheDir() string {
	return f.UserCacheHomeVal
}

func (f *FakeDirs) UserCacheDirs() []string {
	if f.UserCacheDirsVal != nil {
		return f.UserCacheDirsVal
	}
	if f.UserCacheHomeVal != "" {
		return []string{f.UserCacheHomeVal}
	}
	return nil
}

func (f *FakeDirs) UserCachePath(elem ...string) string {
	return path(f.UserCacheHomeVal, elem...)
}

// --- User state ---

func (f *FakeDirs) UserStateDir() string {
	return f.UserStateHomeVal
}

func (f *FakeDirs) UserStateDirs() []string {
	if f.UserStateDirsVal != nil {
		return f.UserStateDirsVal
	}
	if f.UserStateHomeVal != "" {
		return []string{f.UserStateHomeVal}
	}
	return nil
}

func (f *FakeDirs) UserStatePath(elem ...string) string {
	return path(f.UserStateHomeVal, elem...)
}

// --- User log ---

func (f *FakeDirs) UserLogDir() string {
	return f.UserLogHomeVal
}

func (f *FakeDirs) UserLogDirs() []string {
	if f.UserLogDirsVal != nil {
		return f.UserLogDirsVal
	}
	if f.UserLogHomeVal != "" {
		return []string{f.UserLogHomeVal}
	}
	return nil
}

func (f *FakeDirs) UserLogPath(elem ...string) string {
	return path(f.UserLogHomeVal, elem...)
}

// --- User runtime ---

func (f *FakeDirs) UserRuntimeDir() (string, error) {
	if f.UserRuntimeDirErr != nil {
		return "", f.UserRuntimeDirErr
	}
	return f.UserRuntimeDirVal, nil
}

func (f *FakeDirs) UserRuntimePath(elem ...string) (string, error) {
	base, err := f.UserRuntimeDir()
	if err != nil {
		return "", err
	}
	return path(base, elem...), nil
}

// --- System config ---

func (f *FakeDirs) SystemConfigDirs() []string {
	return f.SystemConfigDirsVal
}

func (f *FakeDirs) SystemConfigDir() string {
	if len(f.SystemConfigDirsVal) == 0 {
		return ""
	}
	return f.SystemConfigDirsVal[0]
}

func (f *FakeDirs) SystemConfigPath(elem ...string) string {
	return path(f.SystemConfigDir(), elem...)
}

// --- System data ---

func (f *FakeDirs) SystemDataDirs() []string {
	return f.SystemDataDirsVal
}

func (f *FakeDirs) SystemDataDir() string {
	if len(f.SystemDataDirsVal) == 0 {
		return ""
	}
	return f.SystemDataDirsVal[0]
}

func (f *FakeDirs) SystemDataPath(elem ...string) string {
	return path(f.SystemDataDir(), elem...)
}

// --- System cache ---

func (f *FakeDirs) SystemCacheDir() string {
	return f.SystemCacheDirVal
}

func (f *FakeDirs) SystemCachePath(elem ...string) string {
	return path(f.SystemCacheDirVal, elem...)
}

// --- System state ---

func (f *FakeDirs) SystemStateDir() string {
	return f.SystemStateDirVal
}

func (f *FakeDirs) SystemStatePath(elem ...string) string {
	return path(f.SystemStateDirVal, elem...)
}

// --- System log ---

func (f *FakeDirs) SystemLogDir() string {
	return f.SystemLogDirVal
}

func (f *FakeDirs) SystemLogPath(elem ...string) string {
	return path(f.SystemLogDirVal, elem...)
}

// --- System runtime ---

func (f *FakeDirs) SystemRuntimeDir() string {
	return f.SystemRuntimeDirVal
}

func (f *FakeDirs) SystemRuntimePath(elem ...string) string {
	if f.SystemRuntimeDirVal == "" {
		return ""
	}
	return path(f.SystemRuntimeDirVal, elem...)
}

// --- Find utilities ---

func (f *FakeDirs) FindConfigFile(filename string) (string, bool) {
	for _, p := range f.AllConfigPaths(filename) {
		if f.fileExists(p) {
			return p, true
		}
	}
	return "", false
}

func (f *FakeDirs) AllConfigPaths(filename string) []string {
	var paths []string
	for _, dir := range f.UserConfigDirs() {
		paths = append(paths, filepath.Join(dir, filename))
	}
	for _, dir := range f.SystemConfigDirs() {
		paths = append(paths, filepath.Join(dir, filename))
	}
	return paths
}

func (f *FakeDirs) ExistingConfigFiles(filename string) []string {
	var existing []string
	for _, p := range f.AllConfigPaths(filename) {
		if f.fileExists(p) {
			existing = append(existing, p)
		}
	}
	return existing
}

func (f *FakeDirs) FindDataFile(filename string) (string, bool) {
	for _, p := range f.AllDataPaths(filename) {
		if f.fileExists(p) {
			return p, true
		}
	}
	return "", false
}

func (f *FakeDirs) AllDataPaths(filename string) []string {
	var paths []string
	for _, dir := range f.UserDataDirs() {
		paths = append(paths, filepath.Join(dir, filename))
	}
	for _, dir := range f.SystemDataDirs() {
		paths = append(paths, filepath.Join(dir, filename))
	}
	return paths
}

func (f *FakeDirs) ExistingDataFiles(filename string) []string {
	var existing []string
	for _, p := range f.AllDataPaths(filename) {
		if f.fileExists(p) {
			existing = append(existing, p)
		}
	}
	return existing
}

func (f *FakeDirs) FindCacheFile(filename string) (string, bool) {
	for _, p := range f.AllCachePaths(filename) {
		if f.fileExists(p) {
			return p, true
		}
	}
	return "", false
}

func (f *FakeDirs) AllCachePaths(filename string) []string {
	var paths []string
	for _, dir := range f.UserCacheDirs() {
		paths = append(paths, filepath.Join(dir, filename))
	}
	if f.SystemCacheDirVal != "" {
		paths = append(paths, filepath.Join(f.SystemCacheDirVal, filename))
	}
	return paths
}

func (f *FakeDirs) ExistingCacheFiles(filename string) []string {
	var existing []string
	for _, p := range f.AllCachePaths(filename) {
		if f.fileExists(p) {
			existing = append(existing, p)
		}
	}
	return existing
}

func (f *FakeDirs) FindStateFile(filename string) (string, bool) {
	for _, p := range f.AllStatePaths(filename) {
		if f.fileExists(p) {
			return p, true
		}
	}
	return "", false
}

func (f *FakeDirs) AllStatePaths(filename string) []string {
	var paths []string
	for _, dir := range f.UserStateDirs() {
		paths = append(paths, filepath.Join(dir, filename))
	}
	if f.SystemStateDirVal != "" {
		paths = append(paths, filepath.Join(f.SystemStateDirVal, filename))
	}
	return paths
}

func (f *FakeDirs) ExistingStateFiles(filename string) []string {
	var existing []string
	for _, p := range f.AllStatePaths(filename) {
		if f.fileExists(p) {
			existing = append(existing, p)
		}
	}
	return existing
}

func (f *FakeDirs) FindLogFile(filename string) (string, bool) {
	for _, p := range f.AllLogPaths(filename) {
		if f.fileExists(p) {
			return p, true
		}
	}
	return "", false
}

func (f *FakeDirs) AllLogPaths(filename string) []string {
	var paths []string
	for _, dir := range f.UserLogDirs() {
		paths = append(paths, filepath.Join(dir, filename))
	}
	if f.SystemLogDirVal != "" {
		paths = append(paths, filepath.Join(f.SystemLogDirVal, filename))
	}
	return paths
}

func (f *FakeDirs) ExistingLogFiles(filename string) []string {
	var existing []string
	for _, p := range f.AllLogPaths(filename) {
		if f.fileExists(p) {
			existing = append(existing, p)
		}
	}
	return existing
}

func (f *FakeDirs) FindRuntimeFile(filename string) (string, bool) {
	for _, p := range f.AllRuntimePaths(filename) {
		if p != "" && f.fileExists(p) {
			return p, true
		}
	}
	return "", false
}

func (f *FakeDirs) AllRuntimePaths(filename string) []string {
	var paths []string
	if f.UserRuntimeDirVal != "" && f.UserRuntimeDirErr == nil {
		paths = append(paths, filepath.Join(f.UserRuntimeDirVal, filename))
	}
	if f.SystemRuntimeDirVal != "" {
		paths = append(paths, filepath.Join(f.SystemRuntimeDirVal, filename))
	}
	return paths
}

func (f *FakeDirs) ExistingRuntimeFiles(filename string) []string {
	var existing []string
	for _, p := range f.AllRuntimePaths(filename) {
		if p != "" && f.fileExists(p) {
			existing = append(existing, p)
		}
	}
	return existing
}

// --- Ensure utilities ---

func (f *FakeDirs) EnsureUserConfigDir() (string, error) {
	if err := f.EnsureErrors["config"]; err != nil {
		return "", err
	}
	if f.CreateDirs {
		if err := os.MkdirAll(f.UserConfigHomeVal, 0o700); err != nil {
			return "", err
		}
	}
	return f.UserConfigHomeVal, nil
}

func (f *FakeDirs) EnsureUserDataDir() (string, error) {
	if err := f.EnsureErrors["data"]; err != nil {
		return "", err
	}
	if f.CreateDirs {
		if err := os.MkdirAll(f.UserDataHomeVal, 0o700); err != nil {
			return "", err
		}
	}
	return f.UserDataHomeVal, nil
}

func (f *FakeDirs) EnsureUserCacheDir() (string, error) {
	if err := f.EnsureErrors["cache"]; err != nil {
		return "", err
	}
	if f.CreateDirs {
		if err := os.MkdirAll(f.UserCacheHomeVal, 0o700); err != nil {
			return "", err
		}
	}
	return f.UserCacheHomeVal, nil
}

func (f *FakeDirs) EnsureUserStateDir() (string, error) {
	if err := f.EnsureErrors["state"]; err != nil {
		return "", err
	}
	if f.CreateDirs {
		if err := os.MkdirAll(f.UserStateHomeVal, 0o700); err != nil {
			return "", err
		}
	}
	return f.UserStateHomeVal, nil
}

func (f *FakeDirs) EnsureUserLogDir() (string, error) {
	if err := f.EnsureErrors["log"]; err != nil {
		return "", err
	}
	if f.CreateDirs {
		if err := os.MkdirAll(f.UserLogHomeVal, 0o700); err != nil {
			return "", err
		}
	}
	return f.UserLogHomeVal, nil
}

// --- Project discovery methods ---

// FindUp walks up from start, returning the first directory containing any marker.
func (f *FakeDirs) FindUp(start string, markers ...string) (string, string, bool) {
	matches := f.walkUp(start, markers, nil, nil, false)
	if len(matches) == 0 {
		return "", "", false
	}
	return matches[0].Dir, matches[0].Marker, true
}

// FindUpFunc adds a predicate to FindUp.
func (f *FakeDirs) FindUpFunc(
	start string,
	markers []string,
	match func(markerPath string) bool,
) (string, string, bool) {
	matches := f.walkUp(start, markers, nil, match, false)
	if len(matches) == 0 {
		return "", "", false
	}
	return matches[0].Dir, matches[0].Marker, true
}

// FindUpUntil stops traversal when encountering a stop marker.
func (f *FakeDirs) FindUpUntil(
	start string,
	markers, stopAt []string,
) (string, string, bool) {
	matches := f.walkUp(start, markers, stopAt, nil, false)
	if len(matches) == 0 {
		return "", "", false
	}
	return matches[0].Dir, matches[0].Marker, true
}

// FindUpUntilFunc combines predicate and stop behavior.
func (f *FakeDirs) FindUpUntilFunc(
	start string,
	markers, stopAt []string,
	match func(markerPath string) bool,
) (string, string, bool) {
	matches := f.walkUp(start, markers, stopAt, match, false)
	if len(matches) == 0 {
		return "", "", false
	}
	return matches[0].Dir, matches[0].Marker, true
}

// FindAllUp returns all directories containing any marker.
func (f *FakeDirs) FindAllUp(start string, markers ...string) []Match {
	return f.walkUp(start, markers, nil, nil, true)
}

// FindAllUpFunc filters matches through a predicate.
func (f *FakeDirs) FindAllUpFunc(
	start string,
	markers []string,
	match func(markerPath string) bool,
) []Match {
	return f.walkUp(start, markers, nil, match, true)
}

// FindAllUpUntil collects matches until encountering a stop marker.
func (f *FakeDirs) FindAllUpUntil(start string, markers, stopAt []string) []Match {
	return f.walkUp(start, markers, stopAt, nil, true)
}

// FindAllUpUntilFunc combines collection, predicate, and stop behavior.
func (f *FakeDirs) FindAllUpUntilFunc(
	start string,
	markers, stopAt []string,
	match func(markerPath string) bool,
) []Match {
	return f.walkUp(start, markers, stopAt, match, true)
}

// walkUp is the internal traversal function for FakeDirs.
func (f *FakeDirs) walkUp(
	start string,
	markers, stopAt []string,
	matchFn func(string) bool,
	collectAll bool,
) []Match {
	if len(markers) == 0 {
		return nil
	}

	dir := cleanAbsPath(start)
	var results []Match

	for {
		if match, found := f.checkMarkers(dir, markers, matchFn); found {
			results = append(results, match)
			if !collectAll {
				return results
			}
		}

		if f.shouldStop(dir, stopAt) {
			return results
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return results
}

// checkMarkers checks if any marker exists in the directory.
func (f *FakeDirs) checkMarkers(
	dir string,
	markers []string,
	matchFn func(string) bool,
) (Match, bool) {
	for _, m := range markers {
		markerPath := filepath.Join(dir, m)
		if f.fileExists(markerPath) {
			if matchFn == nil || matchFn(markerPath) {
				return Match{Dir: dir, Marker: m}, true
			}
		}
	}
	return Match{}, false
}

// shouldStop checks if any stop marker exists in the directory.
func (f *FakeDirs) shouldStop(dir string, stopAt []string) bool {
	for _, s := range stopAt {
		if f.fileExists(filepath.Join(dir, s)) {
			return true
		}
	}
	return false
}
