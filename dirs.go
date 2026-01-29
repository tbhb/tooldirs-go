package toolpaths

// Platform represents the detected or overridden operating system.
type Platform int

const (
	// PlatformAuto indicates automatic platform detection based on runtime.GOOS.
	// This is the zero value and default behavior.
	PlatformAuto Platform = iota
	PlatformLinux
	PlatformMacOS
	PlatformWindows
	PlatformFreeBSD
	PlatformOpenBSD
)

func (p Platform) String() string {
	switch p {
	case PlatformAuto:
		return "auto"
	case PlatformLinux:
		return "linux"
	case PlatformMacOS:
		return "macos"
	case PlatformWindows:
		return "windows"
	case PlatformFreeBSD:
		return "freebsd"
	case PlatformOpenBSD:
		return "openbsd"
	default:
		return "unknown"
	}
}

// Config controls how directory paths are resolved.
type Config struct {
	// AppName is required. Used as the directory name.
	AppName string

	// AppAuthor is optional. On Windows, paths become {AppAuthor}/{AppName}.
	// Ignored on other platforms.
	AppAuthor string

	// Version is optional. If set, appended as a subdirectory.
	Version string

	// Roaming controls Windows behavior only.
	// true = FOLDERID_RoamingAppData
	// false = FOLDERID_LocalAppData
	Roaming bool

	// XDGOnAllPlatforms controls whether XDG conventions are used on
	// non-Linux/BSD platforms (macOS, Windows).
	//
	// When false (default): XDG env vars are respected if explicitly set,
	// but XDG default paths (e.g., ~/.config) are not used.
	//
	// When true: XDG env vars AND default paths are used as the primary
	// resolution strategy, with platform-native as fallback.
	XDGOnAllPlatforms bool

	// IncludeXDGFallbacks controls whether User*Dirs() methods include
	// XDG default locations as fallbacks on non-XDG platforms.
	//
	// When true (default): UserConfigDirs on macOS returns both the native
	// path and the XDG default, e.g., [~/Library/Application Support/myapp, ~/.config/myapp].
	// This enables migration from XDG paths to native paths.
	//
	// When false: Only the native path is returned.
	//
	// This option has no effect if XDGOnAllPlatforms is true (XDG paths are already primary).
	// This option has no effect on XDG platforms (Linux, FreeBSD, OpenBSD) where XDG is native.
	IncludeXDGFallbacks *bool

	// EnvOverrides allows specifying app-specific environment variables
	// that take precedence over all other resolution strategies.
	// If the env var is set and non-empty, its value is used directly
	// (with AppName/Version appended according to AppendAppName setting).
	EnvOverrides *EnvOverrides

	// Platform overrides OS detection. Useful for testing.
	// Leave as PlatformAuto (zero value) for automatic detection.
	Platform Platform
}

// EnvOverrides specifies app-specific environment variables for each
// directory type. If set and non-empty, these take absolute precedence.
type EnvOverrides struct {
	// If true, AppName (and Version if set) are appended to env var values.
	// If false, env var value is used as-is.
	AppendAppName bool

	UserConfig  string // e.g., "MYAPP_CONFIG_HOME"
	UserData    string // e.g., "MYAPP_DATA_HOME"
	UserCache   string // e.g., "MYAPP_CACHE_HOME"
	UserState   string // e.g., "MYAPP_STATE_HOME"
	UserLog     string // e.g., "MYAPP_LOG_HOME"
	UserRuntime string // e.g., "MYAPP_RUNTIME_DIR"

	SystemConfig  string // e.g., "MYAPP_SYSTEM_CONFIG"
	SystemData    string // e.g., "MYAPP_SYSTEM_DATA"
	SystemCache   string // e.g., "MYAPP_SYSTEM_CACHE"
	SystemState   string // e.g., "MYAPP_SYSTEM_STATE"
	SystemLog     string // e.g., "MYAPP_SYSTEM_LOG"
	SystemRuntime string // e.g., "MYAPP_SYSTEM_RUNTIME"
}

// get returns the env var name for the given directory type.
// Returns empty string if the type is not recognized.
func (e *EnvOverrides) get(dt dirType) string {
	switch dt {
	case userConfig:
		return e.UserConfig
	case userData:
		return e.UserData
	case userCache:
		return e.UserCache
	case userState:
		return e.UserState
	case userLog:
		return e.UserLog
	case userRuntime:
		return e.UserRuntime
	case systemConfig:
		return e.SystemConfig
	case systemData:
		return e.SystemData
	case systemCache:
		return e.SystemCache
	case systemState:
		return e.SystemState
	case systemLog:
		return e.SystemLog
	case systemRuntime:
		return e.SystemRuntime
	default:
		return ""
	}
}

// Dirs defines the interface for platform directory resolution.
// Use this interface in application code to enable testing with FakeDirs.
type Dirs interface {
	// User config directories
	// UserConfigDir returns the primary user config directory (for writing)
	UserConfigDir() string
	// UserConfigDirs returns all user config directories in priority order (for reading)
	UserConfigDirs() []string
	// UserConfigPath joins path elements to the primary user config directory
	UserConfigPath(elem ...string) string

	// User data directories
	UserDataDir() string
	UserDataDirs() []string
	UserDataPath(elem ...string) string

	// User cache directories
	UserCacheDir() string
	UserCacheDirs() []string
	UserCachePath(elem ...string) string

	// User state directories
	UserStateDir() string
	UserStateDirs() []string
	UserStatePath(elem ...string) string

	// User log directories
	UserLogDir() string
	UserLogDirs() []string
	UserLogPath(elem ...string) string

	// User runtime directory
	UserRuntimeDir() (string, error)
	UserRuntimePath(elem ...string) (string, error)

	// System config directories
	SystemConfigDirs() []string
	SystemConfigDir() string
	SystemConfigPath(elem ...string) string

	// System data directories
	SystemDataDirs() []string
	SystemDataDir() string
	SystemDataPath(elem ...string) string

	// System cache directory
	SystemCacheDir() string
	SystemCachePath(elem ...string) string

	// System state directory
	SystemStateDir() string
	SystemStatePath(elem ...string) string

	// System log directory
	SystemLogDir() string
	SystemLogPath(elem ...string) string

	// System runtime directory
	SystemRuntimeDir() string
	SystemRuntimePath(elem ...string) string

	// Find utilities
	FindConfigFile(filename string) (string, bool)
	AllConfigPaths(filename string) []string
	ExistingConfigFiles(filename string) []string

	FindDataFile(filename string) (string, bool)
	AllDataPaths(filename string) []string
	ExistingDataFiles(filename string) []string

	FindCacheFile(filename string) (string, bool)
	AllCachePaths(filename string) []string
	ExistingCacheFiles(filename string) []string

	FindStateFile(filename string) (string, bool)
	AllStatePaths(filename string) []string
	ExistingStateFiles(filename string) []string

	FindLogFile(filename string) (string, bool)
	AllLogPaths(filename string) []string
	ExistingLogFiles(filename string) []string

	FindRuntimeFile(filename string) (string, bool)
	AllRuntimePaths(filename string) []string
	ExistingRuntimeFiles(filename string) []string

	// Ensure utilities create directories if they don't exist
	EnsureUserConfigDir() (string, error)
	EnsureUserDataDir() (string, error)
	EnsureUserCacheDir() (string, error)
	EnsureUserStateDir() (string, error)
	EnsureUserLogDir() (string, error)

	// Project discovery methods walk up the directory tree to find markers.
	// These are primitives for finding project roots, workspace boundaries,
	// and cascading configuration files.

	// FindUp returns the first directory containing any of the specified markers.
	FindUp(start string, markers ...string) (string, string, bool)

	// FindUpFunc adds a predicate. A marker only matches if match(markerPath) returns true.
	FindUpFunc(
		start string,
		markers []string,
		match func(markerPath string) bool,
	) (string, string, bool)

	// FindUpUntil stops traversal when a directory contains any stopAt marker.
	FindUpUntil(start string, markers, stopAt []string) (string, string, bool)

	// FindUpUntilFunc combines predicate validation with stop markers.
	FindUpUntilFunc(
		start string,
		markers, stopAt []string,
		match func(markerPath string) bool,
	) (string, string, bool)

	// FindAllUp returns all directories containing any marker, nearest to farthest.
	FindAllUp(start string, markers ...string) []Match

	// FindAllUpFunc filters matches through a predicate.
	FindAllUpFunc(
		start string,
		markers []string,
		match func(markerPath string) bool,
	) []Match

	// FindAllUpUntil collects matches until encountering a stop marker.
	FindAllUpUntil(start string, markers, stopAt []string) []Match

	// FindAllUpUntilFunc combines collection, predicate, and stop behavior.
	FindAllUpUntilFunc(
		start string,
		markers, stopAt []string,
		match func(markerPath string) bool,
	) []Match
}

// Compile-time check that PlatformDirs implements Dirs.
var _ Dirs = (*PlatformDirs)(nil)
