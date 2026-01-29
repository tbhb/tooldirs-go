package tooldirs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ErrAppNameRequired is returned when Config.AppName is empty or whitespace-only.
var ErrAppNameRequired = errors.New("tooldirs: AppName is required")

// PlatformDirs provides access to platform-appropriate directories.
type PlatformDirs struct {
	cfg      Config
	platform Platform
}

// New creates a PlatformDirs instance with default configuration.
// Returns ErrAppNameRequired if appName is empty or whitespace-only.
func New(appName string) (*PlatformDirs, error) {
	return NewWithConfig(Config{AppName: appName})
}

// NewWithConfig creates a PlatformDirs instance with custom configuration.
// Returns ErrAppNameRequired if Config.AppName is empty or whitespace-only.
func NewWithConfig(cfg Config) (*PlatformDirs, error) {
	if strings.TrimSpace(cfg.AppName) == "" {
		return nil, ErrAppNameRequired
	}
	platform := cfg.Platform
	if platform == PlatformAuto {
		platform = detectPlatform()
	}
	return &PlatformDirs{cfg: cfg, platform: platform}, nil
}

func detectPlatform() Platform {
	switch runtime.GOOS {
	case "darwin":
		return PlatformMacOS
	case "windows":
		return PlatformWindows
	case "freebsd":
		return PlatformFreeBSD
	case "openbsd":
		return PlatformOpenBSD
	default:
		return PlatformLinux
	}
}

// String returns a human-readable summary of all resolved directory paths.
func (d *PlatformDirs) String() string {
	var b strings.Builder

	b.WriteString("tooldirs.PlatformDirs{\n")
	fmt.Fprintf(&b, "  AppName:  %q\n", d.cfg.AppName)
	if d.cfg.Version != "" {
		fmt.Fprintf(&b, "  Version:  %q\n", d.cfg.Version)
	}
	fmt.Fprintf(&b, "  Platform: %s\n", d.platform)
	b.WriteString("\n")

	b.WriteString("  User directories:\n")
	fmt.Fprintf(&b, "    Config:  %s\n", d.UserConfigDir())
	fmt.Fprintf(&b, "    Data:    %s\n", d.UserDataDir())
	fmt.Fprintf(&b, "    Cache:   %s\n", d.UserCacheDir())
	fmt.Fprintf(&b, "    State:   %s\n", d.UserStateDir())
	fmt.Fprintf(&b, "    Log:     %s\n", d.UserLogDir())
	if userRuntime, err := d.UserRuntimeDir(); err == nil {
		fmt.Fprintf(&b, "    Runtime: %s\n", userRuntime)
	} else {
		fmt.Fprintf(&b, "    Runtime: <error: %v>\n", err)
	}
	b.WriteString("\n")

	b.WriteString("  System directories:\n")
	fmt.Fprintf(&b, "    Config:  %s\n", d.SystemConfigDir())
	fmt.Fprintf(&b, "    Data:    %s\n", d.SystemDataDir())
	fmt.Fprintf(&b, "    Cache:   %s\n", d.SystemCacheDir())
	fmt.Fprintf(&b, "    State:   %s\n", d.SystemStateDir())
	fmt.Fprintf(&b, "    Log:     %s\n", d.SystemLogDir())
	if sysRuntime := d.SystemRuntimeDir(); sysRuntime != "" {
		fmt.Fprintf(&b, "    Runtime: %s\n", sysRuntime)
	} else {
		b.WriteString("    Runtime: <not available>\n")
	}
	b.WriteString("}")

	return b.String()
}

// ---------------------------------------------------------------------
// Path resolution helpers
// ---------------------------------------------------------------------

// path joins path elements onto a base directory.
func path(base string, elem ...string) string {
	return filepath.Join(append([]string{base}, elem...)...)
}

// ---------------------------------------------------------------------
// User directories
// ---------------------------------------------------------------------

// UserConfigDir returns the user-specific configuration directory.
// This is the primary location for writing config files.
func (d *PlatformDirs) UserConfigDir() string {
	if dir := d.fromEnvOverride(userConfig); dir != "" {
		return dir
	}
	return d.resolveUserDir(userConfig)
}

// UserConfigDirs returns all user config directories in priority order.
// The first element is always UserConfigDir (for writing). Additional elements
// are fallback locations to search when reading. On non-XDG platforms with
// IncludeXDGFallbacks enabled, this includes XDG default paths for migration.
func (d *PlatformDirs) UserConfigDirs() []string {
	return d.userDirsWithFallbacks(userConfig)
}

// UserConfigPath returns a path within the user config directory.
func (d *PlatformDirs) UserConfigPath(elem ...string) string {
	return path(d.UserConfigDir(), elem...)
}

// UserDataDir returns the user-specific data directory.
func (d *PlatformDirs) UserDataDir() string {
	if dir := d.fromEnvOverride(userData); dir != "" {
		return dir
	}
	return d.resolveUserDir(userData)
}

// UserDataDirs returns all user data directories in priority order.
func (d *PlatformDirs) UserDataDirs() []string {
	return d.userDirsWithFallbacks(userData)
}

// UserDataPath returns a path within the user data directory.
func (d *PlatformDirs) UserDataPath(elem ...string) string {
	return path(d.UserDataDir(), elem...)
}

// UserCacheDir returns the user-specific cache directory.
func (d *PlatformDirs) UserCacheDir() string {
	if dir := d.fromEnvOverride(userCache); dir != "" {
		return dir
	}
	return d.resolveUserDir(userCache)
}

// UserCacheDirs returns all user cache directories in priority order.
func (d *PlatformDirs) UserCacheDirs() []string {
	return d.userDirsWithFallbacks(userCache)
}

// UserCachePath returns a path within the user cache directory.
func (d *PlatformDirs) UserCachePath(elem ...string) string {
	return path(d.UserCacheDir(), elem...)
}

// UserStateDir returns the user-specific state directory.
func (d *PlatformDirs) UserStateDir() string {
	if dir := d.fromEnvOverride(userState); dir != "" {
		return dir
	}
	return d.resolveUserDir(userState)
}

// UserStateDirs returns all user state directories in priority order.
func (d *PlatformDirs) UserStateDirs() []string {
	return d.userDirsWithFallbacks(userState)
}

// UserStatePath returns a path within the user state directory.
func (d *PlatformDirs) UserStatePath(elem ...string) string {
	return path(d.UserStateDir(), elem...)
}

// UserLogDir returns the user-specific log directory.
func (d *PlatformDirs) UserLogDir() string {
	if dir := d.fromEnvOverride(userLog); dir != "" {
		return dir
	}
	return d.resolveUserDir(userLog)
}

// UserLogDirs returns all user log directories in priority order.
func (d *PlatformDirs) UserLogDirs() []string {
	return d.userDirsWithFallbacks(userLog)
}

// UserLogPath returns a path within the user log directory.
func (d *PlatformDirs) UserLogPath(elem ...string) string {
	return path(d.UserLogDir(), elem...)
}

// UserRuntimeDir returns the user-specific runtime directory.
// Returns an error if the runtime directory cannot be determined
// (e.g., XDG_RUNTIME_DIR not set on Linux with no fallback).
func (d *PlatformDirs) UserRuntimeDir() (string, error) {
	if dir := d.fromEnvOverride(userRuntime); dir != "" {
		return dir, nil
	}
	return d.resolveRuntimeDir()
}

// UserRuntimePath returns a path within the user runtime directory.
func (d *PlatformDirs) UserRuntimePath(elem ...string) (string, error) {
	base, err := d.UserRuntimeDir()
	if err != nil {
		return "", err
	}
	return path(base, elem...), nil
}

// ---------------------------------------------------------------------
// System directories
// ---------------------------------------------------------------------

// SystemConfigDirs returns system-wide configuration directories
// in priority order (highest priority first).
func (d *PlatformDirs) SystemConfigDirs() []string {
	if dir := d.fromEnvOverride(systemConfig); dir != "" {
		return []string{dir}
	}
	return d.resolveSystemDirs(systemConfig)
}

// SystemConfigDir returns the primary system configuration directory.
// This is a convenience for SystemConfigDirs()[0].
func (d *PlatformDirs) SystemConfigDir() string {
	dirs := d.SystemConfigDirs()
	if len(dirs) == 0 {
		return ""
	}
	return dirs[0]
}

// SystemConfigPath returns a path within the primary system config directory.
func (d *PlatformDirs) SystemConfigPath(elem ...string) string {
	return path(d.SystemConfigDir(), elem...)
}

// SystemDataDirs returns system-wide data directories in priority order.
func (d *PlatformDirs) SystemDataDirs() []string {
	if dir := d.fromEnvOverride(systemData); dir != "" {
		return []string{dir}
	}
	return d.resolveSystemDirs(systemData)
}

// SystemDataDir returns the primary system data directory.
func (d *PlatformDirs) SystemDataDir() string {
	dirs := d.SystemDataDirs()
	if len(dirs) == 0 {
		return ""
	}
	return dirs[0]
}

// SystemDataPath returns a path within the primary system data directory.
func (d *PlatformDirs) SystemDataPath(elem ...string) string {
	return path(d.SystemDataDir(), elem...)
}

// SystemCacheDir returns the system-wide cache directory.
// Unlike SystemConfigDirs/SystemDataDirs, this is a single location (not a search path).
func (d *PlatformDirs) SystemCacheDir() string {
	if dir := d.fromEnvOverride(systemCache); dir != "" {
		return dir
	}
	return d.resolveSystemSingleDir(systemCache)
}

// SystemCachePath returns a path within the system cache directory.
func (d *PlatformDirs) SystemCachePath(elem ...string) string {
	return path(d.SystemCacheDir(), elem...)
}

// SystemStateDir returns the system-wide state directory.
// This is for persistent data that isn't user-facing (databases, etc.).
func (d *PlatformDirs) SystemStateDir() string {
	if dir := d.fromEnvOverride(systemState); dir != "" {
		return dir
	}
	return d.resolveSystemSingleDir(systemState)
}

// SystemStatePath returns a path within the system state directory.
func (d *PlatformDirs) SystemStatePath(elem ...string) string {
	return path(d.SystemStateDir(), elem...)
}

// SystemLogDir returns the system-wide log directory.
func (d *PlatformDirs) SystemLogDir() string {
	if dir := d.fromEnvOverride(systemLog); dir != "" {
		return dir
	}
	return d.resolveSystemSingleDir(systemLog)
}

// SystemLogPath returns a path within the system log directory.
func (d *PlatformDirs) SystemLogPath(elem ...string) string {
	return path(d.SystemLogDir(), elem...)
}

// SystemRuntimeDir returns the system-wide runtime directory.
// On Linux/BSD this is /run/{app}. On macOS and Windows, this returns
// an empty string as there is no equivalent concept.
func (d *PlatformDirs) SystemRuntimeDir() string {
	if dir := d.fromEnvOverride(systemRuntime); dir != "" {
		return dir
	}
	return d.resolveSystemSingleDir(systemRuntime)
}

// SystemRuntimePath returns a path within the system runtime directory.
// Returns empty string on platforms without system runtime directories.
func (d *PlatformDirs) SystemRuntimePath(elem ...string) string {
	base := d.SystemRuntimeDir()
	if base == "" {
		return ""
	}
	return path(base, elem...)
}

// ---------------------------------------------------------------------
// Find utilities
// ---------------------------------------------------------------------

// FindConfigFile finds a file in all config directories
// (user first, then system) and returns the first existing path.
func (d *PlatformDirs) FindConfigFile(filename string) (string, bool) {
	for _, p := range d.AllConfigPaths(filename) {
		if fileExists(p) {
			return p, true
		}
	}
	return "", false
}

// AllConfigPaths returns all possible paths for a config file,
// in priority order (user config first, then system configs).
// Does not check if files exist.
func (d *PlatformDirs) AllConfigPaths(filename string) []string {
	var paths []string
	paths = append(paths, d.UserConfigPath(filename))
	for _, dir := range d.SystemConfigDirs() {
		paths = append(paths, filepath.Join(dir, filename))
	}
	return paths
}

// ExistingConfigFiles returns paths to all existing instances of a
// config file across user and system directories, in priority order.
func (d *PlatformDirs) ExistingConfigFiles(filename string) []string {
	var existing []string
	for _, p := range d.AllConfigPaths(filename) {
		if fileExists(p) {
			existing = append(existing, p)
		}
	}
	return existing
}

// FindDataFile finds a file in all data directories
// (user first, then system) and returns the first existing path.
func (d *PlatformDirs) FindDataFile(filename string) (string, bool) {
	for _, p := range d.AllDataPaths(filename) {
		if fileExists(p) {
			return p, true
		}
	}
	return "", false
}

// AllDataPaths returns all possible paths for a data file,
// in priority order (user first, then system).
// Does not check if files exist.
func (d *PlatformDirs) AllDataPaths(filename string) []string {
	var paths []string
	paths = append(paths, d.UserDataPath(filename))
	for _, dir := range d.SystemDataDirs() {
		paths = append(paths, filepath.Join(dir, filename))
	}
	return paths
}

// ExistingDataFiles returns paths to all existing instances of a
// data file across user and system directories, in priority order.
func (d *PlatformDirs) ExistingDataFiles(filename string) []string {
	var existing []string
	for _, p := range d.AllDataPaths(filename) {
		if fileExists(p) {
			existing = append(existing, p)
		}
	}
	return existing
}

// FindCacheFile finds a file in all cache directories
// (user first, then system) and returns the first existing path.
func (d *PlatformDirs) FindCacheFile(filename string) (string, bool) {
	for _, p := range d.AllCachePaths(filename) {
		if fileExists(p) {
			return p, true
		}
	}
	return "", false
}

// AllCachePaths returns all possible paths for a cache file,
// in priority order (user first, then system).
// Does not check if files exist.
func (d *PlatformDirs) AllCachePaths(filename string) []string {
	return []string{
		d.UserCachePath(filename),
		d.SystemCachePath(filename),
	}
}

// ExistingCacheFiles returns paths to all existing instances of a
// cache file across user and system directories, in priority order.
func (d *PlatformDirs) ExistingCacheFiles(filename string) []string {
	var existing []string
	for _, p := range d.AllCachePaths(filename) {
		if fileExists(p) {
			existing = append(existing, p)
		}
	}
	return existing
}

// FindStateFile finds a file in all state directories
// (user first, then system) and returns the first existing path.
func (d *PlatformDirs) FindStateFile(filename string) (string, bool) {
	for _, p := range d.AllStatePaths(filename) {
		if fileExists(p) {
			return p, true
		}
	}
	return "", false
}

// AllStatePaths returns all possible paths for a state file,
// in priority order (user first, then system).
// Does not check if files exist.
func (d *PlatformDirs) AllStatePaths(filename string) []string {
	return []string{
		d.UserStatePath(filename),
		d.SystemStatePath(filename),
	}
}

// ExistingStateFiles returns paths to all existing instances of a
// state file across user and system directories, in priority order.
func (d *PlatformDirs) ExistingStateFiles(filename string) []string {
	var existing []string
	for _, p := range d.AllStatePaths(filename) {
		if fileExists(p) {
			existing = append(existing, p)
		}
	}
	return existing
}

// FindLogFile finds a file in all log directories
// (user first, then system) and returns the first existing path.
func (d *PlatformDirs) FindLogFile(filename string) (string, bool) {
	for _, p := range d.AllLogPaths(filename) {
		if fileExists(p) {
			return p, true
		}
	}
	return "", false
}

// AllLogPaths returns all possible paths for a log file,
// in priority order (user first, then system).
// Does not check if files exist.
func (d *PlatformDirs) AllLogPaths(filename string) []string {
	return []string{
		d.UserLogPath(filename),
		d.SystemLogPath(filename),
	}
}

// ExistingLogFiles returns paths to all existing instances of a
// log file across user and system directories, in priority order.
func (d *PlatformDirs) ExistingLogFiles(filename string) []string {
	var existing []string
	for _, p := range d.AllLogPaths(filename) {
		if fileExists(p) {
			existing = append(existing, p)
		}
	}
	return existing
}

// FindRuntimeFile finds a file in all runtime directories
// (user first, then system) and returns the first existing path.
// Note: System runtime directories don't exist on macOS/Windows.
func (d *PlatformDirs) FindRuntimeFile(filename string) (string, bool) {
	for _, p := range d.AllRuntimePaths(filename) {
		if p != "" && fileExists(p) {
			return p, true
		}
	}
	return "", false
}

// AllRuntimePaths returns all possible paths for a runtime file,
// in priority order (user first, then system).
// Does not check if files exist. Empty strings are included for
// platforms where certain runtime directories don't exist.
func (d *PlatformDirs) AllRuntimePaths(filename string) []string {
	var paths []string

	// User runtime - may error, in which case skip it
	if userRuntime, err := d.UserRuntimePath(filename); err == nil {
		paths = append(paths, userRuntime)
	}

	// System runtime - may be empty on some platforms
	if sysRuntime := d.SystemRuntimePath(filename); sysRuntime != "" {
		paths = append(paths, sysRuntime)
	}

	return paths
}

// ExistingRuntimeFiles returns paths to all existing instances of a
// runtime file across user and system directories, in priority order.
func (d *PlatformDirs) ExistingRuntimeFiles(filename string) []string {
	var existing []string
	for _, p := range d.AllRuntimePaths(filename) {
		if p != "" && fileExists(p) {
			existing = append(existing, p)
		}
	}
	return existing
}

// ---------------------------------------------------------------------
// Ensure utilities (create directories if needed)
// ---------------------------------------------------------------------

// EnsureUserConfigDir creates the user config directory if it doesn't
// exist and returns its path.
func (d *PlatformDirs) EnsureUserConfigDir() (string, error) {
	dir := d.UserConfigDir()
	return dir, os.MkdirAll(dir, 0o700)
}

// EnsureUserDataDir creates the user data directory if needed.
func (d *PlatformDirs) EnsureUserDataDir() (string, error) {
	dir := d.UserDataDir()
	return dir, os.MkdirAll(dir, 0o700)
}

// EnsureUserCacheDir creates the user cache directory if needed.
func (d *PlatformDirs) EnsureUserCacheDir() (string, error) {
	dir := d.UserCacheDir()
	return dir, os.MkdirAll(dir, 0o700)
}

// EnsureUserStateDir creates the user state directory if needed.
func (d *PlatformDirs) EnsureUserStateDir() (string, error) {
	dir := d.UserStateDir()
	return dir, os.MkdirAll(dir, 0o700)
}

// EnsureUserLogDir creates the user log directory if needed.
func (d *PlatformDirs) EnsureUserLogDir() (string, error) {
	dir := d.UserLogDir()
	return dir, os.MkdirAll(dir, 0o700)
}

// ---------------------------------------------------------------------
// Internal: env override helpers
// ---------------------------------------------------------------------

func (d *PlatformDirs) fromEnvOverride(dt dirType) string {
	if d.cfg.EnvOverrides == nil {
		return ""
	}
	envVar := d.cfg.EnvOverrides.get(dt)
	if envVar == "" {
		return ""
	}
	val := os.Getenv(envVar)
	if val == "" {
		return ""
	}
	if d.cfg.EnvOverrides.AppendAppName {
		return filepath.Join(val, d.appPath())
	}
	return val
}

// ---------------------------------------------------------------------
// Internal: path construction helpers
// ---------------------------------------------------------------------

func (d *PlatformDirs) appPath() string {
	base := d.cfg.AppName
	if d.cfg.Version != "" {
		base = filepath.Join(base, d.cfg.Version)
	}
	return base
}

func (d *PlatformDirs) windowsAppPath() string {
	if d.cfg.AppAuthor != "" {
		return filepath.Join(d.cfg.AppAuthor, d.appPath())
	}
	return d.appPath()
}

func (d *PlatformDirs) isXDGPlatform() bool {
	switch d.platform { //nolint:exhaustive // only XDG platforms need to be listed
	case PlatformLinux, PlatformFreeBSD, PlatformOpenBSD:
		return true
	default:
		return false
	}
}

// includeXDGFallbacks returns whether XDG fallbacks should be included in *Dirs() results.
// Defaults to true if not explicitly set.
func (d *PlatformDirs) includeXDGFallbacks() bool {
	if d.cfg.IncludeXDGFallbacks == nil {
		return true // default
	}
	return *d.cfg.IncludeXDGFallbacks
}

// userDirsWithFallbacks returns user directories with optional XDG fallbacks.
// On XDG platforms, returns just the primary directory.
// On non-XDG platforms with IncludeXDGFallbacks, includes XDG defaults as fallbacks.
func (d *PlatformDirs) userDirsWithFallbacks(dt dirType) []string {
	primary := d.resolveUserDirForFallbacks(dt)

	// On XDG platforms or with XDGOnAllPlatforms, no fallbacks needed
	if d.isXDGPlatform() || d.cfg.XDGOnAllPlatforms {
		return []string{primary}
	}

	// On non-XDG platforms, optionally include XDG fallback
	if !d.includeXDGFallbacks() {
		return []string{primary}
	}

	// Get XDG default path as fallback
	xdgFallback := d.xdgUserDirDefault(dt)
	if xdgFallback == "" || xdgFallback == primary {
		return []string{primary}
	}

	return []string{primary, xdgFallback}
}

// resolveUserDirForFallbacks resolves the primary user directory.
// Unlike resolveUserDir, this doesn't use env overrides (those are handled separately).
func (d *PlatformDirs) resolveUserDirForFallbacks(dt dirType) string {
	if dir := d.fromEnvOverride(dt); dir != "" {
		return dir
	}
	return d.resolveUserDir(dt)
}

// xdgUserDirDefault returns the XDG default path for a user directory type.
// This returns the default without checking XDG env vars.
func (d *PlatformDirs) xdgUserDirDefault(dt dirType) string {
	home, _ := os.UserHomeDir()

	switch dt { //nolint:exhaustive // only user dir types are supported
	case userConfig:
		return filepath.Join(home, ".config", d.appPath())
	case userData:
		return filepath.Join(home, ".local", "share", d.appPath())
	case userCache:
		return filepath.Join(home, ".cache", d.appPath())
	case userState:
		return filepath.Join(home, ".local", "state", d.appPath())
	case userLog:
		stateDir := filepath.Join(home, ".local", "state", d.appPath())
		return filepath.Join(stateDir, "log")
	default:
		return ""
	}
}

// ---------------------------------------------------------------------
// Internal: directory type enum
// ---------------------------------------------------------------------

type dirType int

const (
	userConfig dirType = iota
	userData
	userCache
	userState
	userLog
	userRuntime
	systemConfig
	systemData
	systemCache
	systemState
	systemLog
	systemRuntime
)

// ---------------------------------------------------------------------
// Internal: user directory resolution
// ---------------------------------------------------------------------

func (d *PlatformDirs) resolveUserDir(dt dirType) string {
	// On XDG platforms, always use XDG
	if d.isXDGPlatform() {
		return d.xdgUserDir(dt)
	}

	// On non-XDG platforms, check XDGOnAllPlatforms setting
	if d.cfg.XDGOnAllPlatforms {
		// Use XDG (env var or default)
		return d.xdgUserDir(dt)
	}

	// Only respect XDG env vars if explicitly set, no defaults
	if dir := d.xdgUserDirEnvOnly(dt); dir != "" {
		return dir
	}

	// Platform-native resolution
	switch d.platform { //nolint:exhaustive // XDG platforms handled above
	case PlatformMacOS:
		return d.macOSUserDir(dt)
	case PlatformWindows:
		return d.windowsUserDir(dt)
	default:
		// Shouldn't reach here, but fallback to XDG
		return d.xdgUserDir(dt)
	}
}

func (d *PlatformDirs) xdgUserDir(dt dirType) string {
	home, _ := os.UserHomeDir()

	switch dt { //nolint:exhaustive // only user dir types are supported
	case userConfig:
		if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
			return filepath.Join(dir, d.appPath())
		}
		return filepath.Join(home, ".config", d.appPath())

	case userData:
		if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
			return filepath.Join(dir, d.appPath())
		}
		return filepath.Join(home, ".local", "share", d.appPath())

	case userCache:
		if dir := os.Getenv("XDG_CACHE_HOME"); dir != "" {
			return filepath.Join(dir, d.appPath())
		}
		return filepath.Join(home, ".cache", d.appPath())

	case userState:
		if dir := os.Getenv("XDG_STATE_HOME"); dir != "" {
			return filepath.Join(dir, d.appPath())
		}
		return filepath.Join(home, ".local", "state", d.appPath())

	case userLog:
		// XDG doesn't define a log dir; convention is state/log
		stateDir := d.xdgUserDir(userState)
		return filepath.Join(stateDir, "log")

	default:
		return ""
	}
}

func (d *PlatformDirs) xdgUserDirEnvOnly(dt dirType) string {
	// Only returns non-empty if the env var is explicitly set
	var envVar string
	switch dt { //nolint:exhaustive // only user dir types are supported
	case userConfig:
		envVar = "XDG_CONFIG_HOME"
	case userData:
		envVar = "XDG_DATA_HOME"
	case userCache:
		envVar = "XDG_CACHE_HOME"
	case userState:
		envVar = "XDG_STATE_HOME"
	case userLog:
		if state := os.Getenv("XDG_STATE_HOME"); state != "" {
			return filepath.Join(state, d.appPath(), "log")
		}
		return ""
	default:
		return ""
	}

	if dir := os.Getenv(envVar); dir != "" {
		return filepath.Join(dir, d.appPath())
	}
	return ""
}

func (d *PlatformDirs) macOSUserDir(dt dirType) string {
	home, _ := os.UserHomeDir()
	lib := filepath.Join(home, "Library")

	switch dt { //nolint:exhaustive // only user dir types are supported
	case userConfig, userData, userState:
		return filepath.Join(lib, "Application Support", d.appPath())
	case userCache:
		return filepath.Join(lib, "Caches", d.appPath())
	case userLog:
		return filepath.Join(lib, "Logs", d.appPath())
	default:
		return ""
	}
}

func (d *PlatformDirs) windowsUserDir(dt dirType) string {
	var baseDir string

	switch dt { //nolint:exhaustive // only user dir types are supported
	case userConfig, userData, userState:
		if d.cfg.Roaming {
			baseDir = windowsRoamingAppData()
		} else {
			baseDir = windowsLocalAppData()
		}
		return filepath.Join(baseDir, d.windowsAppPath())

	case userCache:
		baseDir = windowsLocalAppData()
		return filepath.Join(baseDir, d.windowsAppPath(), "cache")

	case userLog:
		baseDir = windowsLocalAppData()
		return filepath.Join(baseDir, d.windowsAppPath(), "log")

	default:
		return ""
	}
}

// ---------------------------------------------------------------------
// Internal: runtime directory resolution
// ---------------------------------------------------------------------

func (d *PlatformDirs) resolveRuntimeDir() (string, error) {
	// Check XDG env var first (on XDG platforms or if XDGOnAllPlatforms)
	if d.isXDGPlatform() || d.cfg.XDGOnAllPlatforms {
		if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
			return filepath.Join(dir, d.appPath()), nil
		}
	}

	// Also check on non-XDG platforms if env var is explicitly set
	if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
		return filepath.Join(dir, d.appPath()), nil
	}

	switch d.platform {
	case PlatformLinux, PlatformFreeBSD, PlatformOpenBSD:
		// XDG_RUNTIME_DIR not set - fall back to temp directory
		// Note: This is technically non-compliant with XDG spec which says
		// the dir should not persist across reboots, but temp is reasonable
		return filepath.Join(os.TempDir(), fmt.Sprintf("%s-%d", d.cfg.AppName, os.Getuid())), nil

	case PlatformMacOS:
		// $TMPDIR is per-user on macOS
		return filepath.Join(os.TempDir(), d.appPath()), nil

	case PlatformWindows:
		return filepath.Join(windowsLocalAppData(), d.windowsAppPath(), "runtime"), nil

	case PlatformAuto:
		// PlatformAuto is resolved to a concrete platform in NewWithConfig.
		// This case should never be reached.
		return "", errors.New("PlatformAuto should have been resolved during construction")
	}
	return "", fmt.Errorf("cannot determine runtime directory for platform %s", d.platform)
}

// ---------------------------------------------------------------------
// Internal: system directory resolution
// ---------------------------------------------------------------------

func (d *PlatformDirs) resolveSystemDirs(dt dirType) []string {
	if d.isXDGPlatform() {
		return d.xdgSystemDirs(dt)
	}

	if d.cfg.XDGOnAllPlatforms {
		return d.xdgSystemDirs(dt)
	}

	// Check for XDG env vars even on non-XDG platforms
	if dirs := d.xdgSystemDirsEnvOnly(dt); len(dirs) > 0 {
		return dirs
	}

	// Platform-native resolution
	switch d.platform { //nolint:exhaustive // XDG platforms handled above
	case PlatformMacOS:
		return d.macOSSystemDirs(dt)
	case PlatformWindows:
		return d.windowsSystemDirs(dt)
	default:
		return d.xdgSystemDirs(dt)
	}
}

func (d *PlatformDirs) xdgSystemDirs(dt dirType) []string {
	var envVar, defaultVal string

	switch dt { //nolint:exhaustive // only system config/data use search paths
	case systemConfig:
		envVar = "XDG_CONFIG_DIRS"
		defaultVal = "/etc/xdg"
	case systemData:
		envVar = "XDG_DATA_DIRS"
		defaultVal = "/usr/local/share:/usr/share"
	default:
		return nil
	}

	val := os.Getenv(envVar)
	if val == "" {
		val = defaultVal
	}

	parts := strings.Split(val, ":")
	var dirs []string
	for _, p := range parts {
		if p != "" {
			dirs = append(dirs, filepath.Join(p, d.appPath()))
		}
	}
	return dirs
}

func (d *PlatformDirs) xdgSystemDirsEnvOnly(dt dirType) []string {
	var envVar string
	switch dt { //nolint:exhaustive // only system config/data use search paths
	case systemConfig:
		envVar = "XDG_CONFIG_DIRS"
	case systemData:
		envVar = "XDG_DATA_DIRS"
	default:
		return nil
	}

	val := os.Getenv(envVar)
	if val == "" {
		return nil
	}

	parts := strings.Split(val, ":")
	var dirs []string
	for _, p := range parts {
		if p != "" {
			dirs = append(dirs, filepath.Join(p, d.appPath()))
		}
	}
	return dirs
}

func (d *PlatformDirs) macOSSystemDirs(dt dirType) []string {
	switch dt { //nolint:exhaustive // only system config/data use search paths
	case systemConfig, systemData:
		return []string{filepath.Join("/Library", "Application Support", d.appPath())}
	default:
		return nil
	}
}

func (d *PlatformDirs) windowsSystemDirs(dt dirType) []string {
	switch dt { //nolint:exhaustive // only system config/data use search paths
	case systemConfig, systemData:
		return []string{filepath.Join(windowsProgramData(), d.windowsAppPath())}
	default:
		return nil
	}
}

// ---------------------------------------------------------------------
// Internal: single system directory resolution (cache, state, log, runtime)
// ---------------------------------------------------------------------

func (d *PlatformDirs) resolveSystemSingleDir(dt dirType) string {
	switch d.platform { //nolint:exhaustive // PlatformAuto resolved during construction
	case PlatformLinux, PlatformFreeBSD, PlatformOpenBSD:
		return d.fhsSystemDir(dt)
	case PlatformMacOS:
		return d.macOSSystemSingleDir(dt)
	case PlatformWindows:
		return d.windowsSystemSingleDir(dt)
	}
	return d.fhsSystemDir(dt)
}

// fhsSystemDir returns FHS-compliant system directories for Linux/BSD.
func (d *PlatformDirs) fhsSystemDir(dt dirType) string {
	switch dt { //nolint:exhaustive // only system single-dir types
	case systemCache:
		return filepath.Join("/var", "cache", d.appPath())
	case systemState:
		return filepath.Join("/var", "lib", d.appPath())
	case systemLog:
		return filepath.Join("/var", "log", d.appPath())
	case systemRuntime:
		return filepath.Join("/run", d.appPath())
	default:
		return ""
	}
}

func (d *PlatformDirs) macOSSystemSingleDir(dt dirType) string {
	switch dt { //nolint:exhaustive // only system single-dir types
	case systemCache:
		return filepath.Join("/Library", "Caches", d.appPath())
	case systemState:
		// macOS doesn't distinguish state from data at the system level
		return filepath.Join("/Library", "Application Support", d.appPath())
	case systemLog:
		return filepath.Join("/Library", "Logs", d.appPath())
	case systemRuntime:
		// No equivalent on macOS
		return ""
	default:
		return ""
	}
}

func (d *PlatformDirs) windowsSystemSingleDir(dt dirType) string {
	programData := windowsProgramData()
	base := filepath.Join(programData, d.windowsAppPath())

	switch dt { //nolint:exhaustive // only system single-dir types
	case systemCache:
		return filepath.Join(base, "cache")
	case systemState:
		// Windows doesn't distinguish state from data
		return base
	case systemLog:
		return filepath.Join(base, "log")
	case systemRuntime:
		// No equivalent on Windows
		return ""
	default:
		return ""
	}
}

// ---------------------------------------------------------------------
// Internal: file utilities
// ---------------------------------------------------------------------

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
