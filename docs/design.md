# toolpaths design document

This document describes the design decisions, behavioral semantics, and configuration options for the toolpaths package.

## Overview

toolpaths resolves platform-appropriate filesystem paths for app data. The library addresses a fundamental tension: apps want portable code, but users expect platform-native behavior. A settings file "should" live in `~/Library/Application Support` on macOS, `%APPDATA%` on Windows, and `~/.config` on Linux.

The library also addresses a secondary concern: migration. Apps that used XDG paths on all platforms may want to transition to native paths while still finding existing user data.

## Core concepts

### Single vs many path methods

The API follows XDG Base Directory Specification semantics with two method families:

`*Dir()` methods return a single path representing the primary location for that directory type. Write new files here. For example, `UserConfigDir()` returns the primary user configuration directory.

`*Dirs()` methods return a slice of paths in priority order. The first element is always the primary directory (same as `*Dir()`). The remaining elements are fallback locations to search when reading. This enables layered configuration and migration scenarios. For example, `UserConfigDirs()` on macOS returns `[~/Library/Application Support/myapp, ~/.config/myapp]` by default.

### Directory types

The library supports six directory types for both user and system scopes:

`Config` is for user-editable configuration files. These are typically small, human-readable files that control app behavior. Examples: `config.yaml`, `settings.json`, keybindings.

`Data` is for app data that is not configuration. This includes databases, downloaded content, generated indices, or any persistent data the app manages. Examples: `history.db`, `plugins/`, `themes/`.

`Cache` is for non-essential data you can regenerate. The system or user may delete cache contents at any time. Examples include HTTP cache, compiled templates, and thumbnail previews.

State is for data that should persist but is not configuration. This is an XDG 0.8 addition that distinguishes "stuff the app tracks" from "stuff the user configures." Examples: recently opened files, cursor positions, undo history.

`Log` is for app log files. On XDG platforms this is a subdirectory of state. On macOS and Windows it has dedicated locations.

`Runtime` is for ephemeral files that should not survive a reboot. This includes Unix sockets, PID files, and other IPC mechanisms. Runtime directories have stricter requirements (user-owned, restricted permissions) and may not exist on all platforms.

### User vs system scope directories

User directories are per-user locations, typically under the home directory. Apps running as a normal user should use these.

System directories are machine-wide locations, typically requiring elevated privileges to write. System services, installers, or apps that need to share data across all users use these.

## Platform behaviors

### Linux, FreeBSD, OpenBSD (XDG platforms)

These platforms follow the XDG Base Directory Specification natively.

User directories resolve as follows. If the corresponding XDG environment variable exists, that value takes precedence. Otherwise the XDG default applies.

| Type      | Environment variable | Default                       |
| --------- | -------------------- | ----------------------------- |
| `Config`  | `XDG_CONFIG_HOME`    | `~/.config`                   |
| `Data`    | `XDG_DATA_HOME`      | `~/.local/share`              |
| `Cache`   | `XDG_CACHE_HOME`     | `~/.cache`                    |
| `State`   | `XDG_STATE_HOME`     | `~/.local/state`              |
| `Log`     | (derived from State) | `~/.local/state/{app}/log`    |
| `Runtime` | `XDG_RUNTIME_DIR`    | `/tmp/{app}-{uid}` (fallback) |

The library appends the app name (and optional version) to form the final path. For example, with `AppName: "myapp"` and `Version: "2"`, `UserConfigDir()` returns `~/.config/myapp/2`.

System directories follow XDG for `Config` and `Data` (which define search paths), and FHS for the rest:

| Type      | Environment variable | Default                       |
| --------- | -------------------- | ----------------------------- |
| `Config`  | `XDG_CONFIG_DIRS`    | `/etc/xdg`                    |
| `Data`    | `XDG_DATA_DIRS`      | `/usr/local/share:/usr/share` |
| `Cache`   | (none)               | `/var/cache`                  |
| `State`   | (none)               | `/var/lib`                    |
| `Log`     | (none)               | `/var/log`                    |
| `Runtime` | (none)               | `/run`                        |

`SystemConfigDirs` and `SystemDataDirs` return slices because XDG defines these as colon-separated search paths. The other system directories are single locations per FHS conventions.

### macOS

macOS has its own conventions rooted in the Library directory structure.

User directories resolve to:

| Type      | Location                              |
| --------- | ------------------------------------- |
| `Config`  | `~/Library/Application Support/{app}` |
| `Data`    | `~/Library/Application Support/{app}` |
| `State`   | `~/Library/Application Support/{app}` |
| `Cache`   | `~/Library/Caches/{app}`              |
| `Log`     | `~/Library/Logs/{app}`                |
| `Runtime` | `$TMPDIR/{app}`                       |

Note that `Config`, `Data`, and `State` all resolve to the same location. This reflects macOS conventions where Application Support serves as a general-purpose per-app container. Apps that need to distinguish these should use subdirectories.

System directories resolve to:

| Type      | Location                             |
| --------- | ------------------------------------ |
| `Config`  | `/Library/Application Support/{app}` |
| `Data`    | `/Library/Application Support/{app}` |
| `State`   | `/Library/Application Support/{app}` |
| `Cache`   | `/Library/Caches/{app}`              |
| `Log`     | `/Library/Logs/{app}`                |
| `Runtime` | (none)                               |

macOS lacks a system runtime directory, so `SystemRuntimeDir` returns an empty string.

### Windows

Windows uses Known Folders accessed via shell API functions.

User directories resolve based on the `Roaming` configuration option:

| Type      | Roaming=false (default)                 | Roaming=true                            |
| --------- | --------------------------------------- | --------------------------------------- |
| `Config`  | `%LOCALAPPDATA%\{author}\{app}`         | `%APPDATA%\{author}\{app}`              |
| `Data`    | `%LOCALAPPDATA%\{author}\{app}`         | `%APPDATA%\{author}\{app}`              |
| `State`   | `%LOCALAPPDATA%\{author}\{app}`         | `%APPDATA%\{author}\{app}`              |
| `Cache`   | `%LOCALAPPDATA%\{author}\{app}\cache`   | `%LOCALAPPDATA%\{author}\{app}\cache`   |
| `Log`     | `%LOCALAPPDATA%\{author}\{app}\log`     | `%LOCALAPPDATA%\{author}\{app}\log`     |
| `Runtime` | `%LOCALAPPDATA%\{author}\{app}\runtime` | `%LOCALAPPDATA%\{author}\{app}\runtime` |

`Cache`, `Log`, and `Runtime` always use `LOCALAPPDATA` regardless of the `Roaming` setting, as these should not roam between machines.

If you set `AppAuthor`, paths include an author directory level: `{AppAuthor}\{AppName}`. This matches Windows conventions where vendor names provide vendor-level organization. If `AppAuthor` is empty, only `AppName` appears in the path.

System directories resolve to:

| Type      | Location                             |
| --------- | ------------------------------------ |
| `Config`  | `%ProgramData%\{author}\{app}`       |
| `Data`    | `%ProgramData%\{author}\{app}`       |
| `State`   | `%ProgramData%\{author}\{app}`       |
| `Cache`   | `%ProgramData%\{author}\{app}\cache` |
| `Log`     | `%ProgramData%\{author}\{app}\log`   |
| `Runtime` | (none)                               |

Windows lacks a system runtime directory, so `SystemRuntimeDir` returns an empty string.

The library uses `golang.org/x/sys/windows` to call `SHGetKnownFolderPath` for Known Folder resolution. This handles cases where users have relocated their profile directories. The library falls back to environment variables if the system call fails.

## Configuration options

### `AppName` (required)

The app name, used as the directory name. This should be a valid directory name on all target platforms. Avoid special characters, spaces, and leading dots.

### `AppAuthor` (optional, Windows only)

The vendor or author name. On Windows, this creates an extra directory level: `{AppAuthor}\{AppName}`. Other platforms ignore this option, as they use flat naming.

### `Version` (optional)

An optional version string appended as a subdirectory. Use this when different major versions of an app should maintain separate data. The full path becomes `{app}/{version}` on Unix or `{author}\{app}\{version}` on Windows.

### `Roaming` (Windows only)

Controls whether Windows user directories use roaming or local app data.

When false (default), directories resolve under `%LOCALAPPDATA%`. Data stays on the local machine.

When true, directories resolve under `%APPDATA%`. In domain environments, this data may roam to other machines the user logs into.

`Cache`, `Log`, and `Runtime` always use local app data regardless of this setting.

### `XDGOnAllPlatforms`

Controls whether XDG conventions apply on non-XDG platforms (macOS, Windows).

When false (default), platform-native paths apply for `*Dir()` methods. XDG environment variables still take precedence if explicitly set. The `*Dirs()` methods include XDG default paths as fallbacks (controlled by `IncludeXDGFallbacks`).

When true, XDG paths apply everywhere as the primary locations. This provides uniform behavior across platforms at the cost of non-native paths. Useful for apps with an existing XDG-based user base or strong cross-platform consistency requirements.

### `IncludeXDGFallbacks`

Controls whether `*Dirs()` methods include XDG default locations as fallbacks on non-XDG platforms. This is a pointer to `bool` (`*bool`) in the `Config` struct.

When nil or true (default), `UserConfigDirs()` on macOS returns both the native path and the XDG default: `[~/Library/Application Support/myapp, ~/.config/myapp]`. This enables migration from XDG paths to native paths.

When set to false, only the native path appears: `[~/Library/Application Support/myapp]`.

This option has no effect if `XDGOnAllPlatforms` is true, since XDG paths are already primary in that mode. This option also has no effect on XDG platforms (Linux, FreeBSD, OpenBSD) where XDG is native.

```go
// Default behavior (nil = true)
dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
    AppName: "myapp",
})

// Explicitly disable XDG fallbacks
falseVal := false
dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
    AppName:             "myapp",
    IncludeXDGFallbacks: &falseVal,
})
```

### `EnvOverrides`

Lets you specify app-specific environment variables that take precedence over all other resolution strategies. This helps apps that want to provide their own override mechanism.

```go
EnvOverrides: &toolpaths.EnvOverrides{
    AppendAppName: true,
    UserConfig:    "MYAPP_CONFIG_HOME",
    UserData:      "MYAPP_DATA_HOME",
    // ... etc
}
```

If `AppendAppName` is true, the library appends `AppName` (and `Version`) to the environment variable value. If false, it uses the value as-is.

### `Platform`

Overrides automatic platform detection. Useful for testing or for apps that need to generate paths for a different platform.

The zero value is `PlatformAuto`, which triggers runtime detection based on `runtime.GOOS`. You can explicitly set any platform value including `PlatformLinux` to override detection:

```go
dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
    AppName:  "testapp",
    Platform: toolpaths.PlatformLinux, // Works on any OS
})
```

## Resolution precedence

For each directory type, resolution follows this precedence order:

1. App-specific environment override (if EnvOverrides configured and variable set)
2. XDG environment variable (if on XDG platform, or if XDGOnAllPlatforms, or if variable explicitly set)
3. Platform-native default

For `*Dirs()` methods, after resolving the Home path, the library may append XDG defaults as fallbacks (on non-XDG platforms with `IncludeXDGFallbacks` enabled).

## Behavioral matrices

### `UserConfigDir` resolution

| Platform | `XDGOnAllPlatforms` | XDG set | Result                                |
| -------- | ------------------- | ------- | ------------------------------------- |
| Linux    | false               | no      | `~/.config/{app}`                     |
| Linux    | false               | yes     | `$XDG_CONFIG_HOME/{app}`              |
| Linux    | true                | no      | `~/.config/{app}`                     |
| Linux    | true                | yes     | `$XDG_CONFIG_HOME/{app}`              |
| macOS    | false               | no      | `~/Library/Application Support/{app}` |
| macOS    | false               | yes     | `$XDG_CONFIG_HOME/{app}`              |
| macOS    | true                | no      | `~/.config/{app}`                     |
| macOS    | true                | yes     | `$XDG_CONFIG_HOME/{app}`              |
| Windows  | false               | no      | `%LOCALAPPDATA%\{app}`                |
| Windows  | false               | yes     | `$XDG_CONFIG_HOME/{app}`              |
| Windows  | true                | no      | `~/.config/{app}`                     |
| Windows  | true                | yes     | `$XDG_CONFIG_HOME/{app}`              |

### `UserConfigDirs` resolution (`IncludeXDGFallbacks=true`)

| Platform | XDGOnAllPlatforms | Result                                                   |
| -------- | ----------------- | -------------------------------------------------------- |
| Linux    | false             | `[~/.config/{app}]`                                      |
| Linux    | true              | `[~/.config/{app}]`                                      |
| macOS    | false             | `[~/Library/Application Support/{app}, ~/.config/{app}]` |
| macOS    | true              | `[~/.config/{app}]`                                      |
| Windows  | false             | `[%LOCALAPPDATA%\{app}, ~/.config/{app}]`                |
| Windows  | true              | `[~/.config/{app}]`                                      |

### `UserConfigDirs` resolution (`IncludeXDGFallbacks=false`)

| Platform | XDGOnAllPlatforms | Result                                  |
| -------- | ----------------- | --------------------------------------- |
| Linux    | false             | `[~/.config/{app}]`                     |
| Linux    | true              | `[~/.config/{app}]`                     |
| macOS    | false             | `[~/Library/Application Support/{app}]` |
| macOS    | true              | `[~/.config/{app}]`                     |
| Windows  | false             | `[%LOCALAPPDATA%\{app}]`                |
| Windows  | true              | `[~/.config/{app}]`                     |

### `SystemConfigDirs` resolution

| Platform | `XDGOnAllPlatforms` | XDG set | Result                                 |
| -------- | ------------------- | ------- | -------------------------------------- |
| Linux    | false               | no      | `[/etc/xdg/{app}]`                     |
| Linux    | false               | yes     | (parsed from `$XDG_CONFIG_DIRS`)       |
| macOS    | false               | no      | `[/Library/Application Support/{app}]` |
| macOS    | false               | yes     | (parsed from `$XDG_CONFIG_DIRS`)       |
| macOS    | true                | no      | `[/etc/xdg/{app}]`                     |
| Windows  | false               | no      | `[%ProgramData%\{app}]`                |
| Windows  | true                | no      | `[/etc/xdg/{app}]`                     |

## Design opinions

### `Config`, `Data`, and `State` collapse on macOS and Windows

On macOS, `Config`, `Data`, and `State` all resolve to `~/Library/Application Support`. On Windows, they all resolve to `%LOCALAPPDATA%` (or `%APPDATA%` if roaming). This reflects the reality that these platforms do not distinguish between these concepts at the filesystem level.

Apps that need separation should use subdirectories within the single resolved path.

### `Log` is a subdirectory of state on XDG platforms

XDG does not define a log directory. Following common practice, `Log` maps to a subdirectory of state: `$XDG_STATE_HOME/{app}/log`. On macOS and Windows, dedicated log locations exist and the library uses them.

### `Runtime` directory fallback

Login managers on conforming systems should set `XDG_RUNTIME_DIR`. When missing, the library falls back to a temporary directory with the UID appended for uniqueness: `/tmp/{app}-{uid}`. This fallback does not guarantee the security properties that `XDG_RUNTIME_DIR` should provide (user ownership, restricted permissions, `tmpfs` backing).

Apps with strict runtime directory requirements should check for `XDG_RUNTIME_DIR` explicitly or use `UserRuntimeDir`'s error return to detect the fallback case.

### No automatic directory creation

The library does not create directories automatically. The `*Dir()` methods return paths that may or may not exist. Apps should use `Ensure*Dir()` methods when they need to guarantee a directory exists:

```go
// EnsureUserConfigDir creates the directory with mode 0700 if it doesn't exist
path, err := dirs.EnsureUserConfigDir()
if err != nil {
    return fmt.Errorf("failed to create config directory: %w", err)
}
```

Available ensure methods:

- `EnsureUserConfigDir()`
- `EnsureUserDataDir()`
- `EnsureUserCacheDir()`
- `EnsureUserStateDir()`
- `EnsureUserLogDir()`

All ensure methods create directories with mode `0700` (user-only access) for security. This design avoids side effects during path resolution and gives apps explicit control over when the library creates directories.

### XDG environment variables respected on all platforms

Even when `XDGOnAllPlatforms` is false, explicitly set XDG environment variables take precedence over platform-native defaults. This allows users to override behavior on any platform, which is useful for development, testing, or users who prefer XDG semantics everywhere.

### Default to including XDG fallbacks

`IncludeXDGFallbacks` defaults to true because migration support is valuable and the cost (an extra path in the search list) is minimal. Apps that want strict native-only behavior can set this to false.

### System directories do not include fallbacks

Unlike user directories, system directories do not include XDG fallbacks on non-XDG platforms. System paths are typically controlled by administrators and installers, not migrated by end users.

## Find utilities

The library provides find utilities for common operations on layered directories:

`Find*File(filename)` returns the first existing path across all directories for that type (user directories first, then system directories).

`All*Paths(filename)` returns all candidate paths without checking existence.

`Existing*Files(filename)` returns all paths where the file actually exists.

These search across both user and system directories in priority order. The typical pattern is:

```go
// Find config file (user overrides system)
if path, ok := dirs.FindConfigFile("config.yaml"); ok {
    return loadConfig(path)
}

// Or load all configs for merging
for _, path := range dirs.ExistingConfigFiles("config.yaml") {
    mergeConfig(path)
}
```

## Package manager compatibility

### Linux app isolation

App isolation systems set XDG environment variables to point to `~/.var/app/{app-id}/`. The library respects these automatically, so isolated apps work correctly without special handling.

### Snap

Snap rewrites `$HOME` to `$SNAP_USER_DATA`. Since the library uses `os.UserHomeDir()` which respects `$HOME`, Snap apps automatically get isolated paths.

### Other package managers

These package managers do not change user data locations. Apps installed via Nix or Homebrew use standard paths.

### General principle

The library does not detect or special-case any package manager. It respects environment variables, which is the correct integration point for isolated environments. Package managers that need custom paths should set the appropriate environment variables.

## Testing considerations

### Platform override

Use `Config.Platform` to test behavior for different platforms:

```go
dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
    AppName:  "testapp",
    Platform: toolpaths.PlatformMacOS,
})
```

The zero value `PlatformAuto` triggers runtime detection. Any explicit platform value, including `PlatformLinux`, overrides detection:

```go
// Test Linux paths on macOS or Windows
dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
    AppName:  "testapp",
    Platform: toolpaths.PlatformLinux,
})
```

### Environment isolation

Tests should use `t.Setenv` (or a similar helper) to set and restore environment variables. Key variables to consider:

- `HOME` (affects all user directories)
- `XDG_CONFIG_HOME`, `XDG_DATA_HOME`, etc. (XDG overrides)
- `APPDATA`, `LOCALAPPDATA`, `ProgramData` (Windows fallbacks)

### App-specific overrides

For integration tests, use EnvOverrides to control paths:

```go
dirs, err := toolpaths.NewWithConfig(toolpaths.Config{
    AppName: "testapp",
    EnvOverrides: &toolpaths.EnvOverrides{
        AppendAppName: false,
        UserConfig:    "TEST_CONFIG_DIR",
    },
})
os.Setenv("TEST_CONFIG_DIR", "/tmp/test-config")
```

### `FakeDirs` test double

The library provides a `FakeDirs` type that implements the `Dirs` interface for testing. This allows tests to control directory paths without touching the filesystem or environment:

```go
// Create a fake with all paths under a base directory
fake := toolpaths.NewFakeDirs("/tmp/test-app")
// fake.UserConfigDir() returns "/tmp/test-app/config"
// fake.UserDataDir() returns "/tmp/test-app/data"

// Or use a temp directory with automatic cleanup
fake, cleanup := toolpaths.NewFakeDirsWithTempDir("test")
defer cleanup()

// Control file existence for search utilities
fake.SetExisting("/tmp/test-app/config/settings.yaml")
path, found := fake.FindConfigFile("settings.yaml") // found == true

// Configure error returns
fake.EnsureErrors["config"] = errors.New("permission denied")
_, err := fake.EnsureUserConfigDir() // returns error

// Actually create directories in tests
fake.CreateDirs = true
fake.EnsureUserConfigDir() // creates the directory
```

Use the `Dirs` interface in app code to enable dependency injection:

```go
type App struct {
    dirs toolpaths.Dirs
}

func NewApp(dirs toolpaths.Dirs) *App {
    return &App{dirs: dirs}
}

// In production
dirs, err := toolpaths.New("myapp")
if err != nil {
    log.Fatal(err)
}
app := NewApp(dirs)

// In tests
fake := toolpaths.NewFakeDirs("/tmp/test")
app := NewApp(fake)
```

## Project discovery methods

The `Dirs` interface includes methods for walking up the directory tree to find project roots, workspace boundaries, and configuration files. These methods are intentionally low-level: they handle traversal and marker detection, leaving content inspection and semantic interpretation to callers.

Placing these on the interface rather than as package-level functions enables testing through `FakeDirs`. Tests can control which markers exist and where, without touching the real filesystem.

### Core types

```go
// Match represents a found marker during upward traversal.
type Match struct {
    Dir    string // Directory containing the marker
    Marker string // The marker that matched (filename or dirname)
}

// Path returns the full path to the marker.
func (m Match) Path() string
```

### Methods

Eight methods cover the combinations of single vs. all results, with vs. without match functions, and with vs. without stop markers.

#### Single result methods

`FindUp` walks up from `start`, returning the first directory containing any of the specified markers.

```go
FindUp(start string, markers ...string) (dir, marker string, found bool)
```

`FindUpFunc` adds a predicate. A marker only matches if it exists AND `match(markerPath)` returns true.

```go
FindUpFunc(start string, markers []string, match func(markerPath string) bool) (dir, marker string, found bool)
```

`FindUpUntil` stops traversal when a directory contains any of the `stopAt` markers.

```go
FindUpUntil(start string, markers, stopAt []string) (dir, marker string, found bool)
```

`FindUpUntilFunc` combines the predicate and stop behavior.

```go
FindUpUntilFunc(start string, markers, stopAt []string, match func(markerPath string) bool) (dir, marker string, found bool)
```

#### Methods returning all matches

`FindAllUp` returns all directories containing any marker, ordered nearest to farthest.

```go
FindAllUp(start string, markers ...string) []Match
```

`FindAllUpFunc` filters matches through a predicate.

```go
FindAllUpFunc(start string, markers []string, match func(markerPath string) bool) []Match
```

`FindAllUpUntil` collects matches until the traversal encounters a stop marker.

```go
FindAllUpUntil(start string, markers, stopAt []string) []Match
```

`FindAllUpUntilFunc` combines collection, predicate, and stop behavior.

```go
FindAllUpUntilFunc(start string, markers, stopAt []string, match func(markerPath string) bool) []Match
```

### Behavioral semantics

#### Traversal order

All methods start at `start` and walk toward the filesystem root. The walker checks each directory in turn. Results appear in order of proximity, with nearest directories first.

#### Marker matching

Markers can be files or directories. A marker matches if it exists in the current directory. When you specify more than one marker, the walker checks them in order; the first existing marker in a directory wins.

#### Predicate validation

In `*Func` variants, the predicate receives the full path to the existing marker. The marker only counts as a match if the predicate returns true. This enables content inspection without the library needing to understand file formats.

```go
// Only match Cargo.toml files that contain a [workspace] section
dirs.FindUpFunc(cwd, []string{"Cargo.toml"}, func(path string) bool {
    content, err := os.ReadFile(path)
    if err != nil {
        return false
    }
    return bytes.Contains(content, []byte("[workspace]"))
})
```

#### Stop markers

In `*Until` variants, traversal stops when a directory contains any stop marker. The walker checks for the target marker before checking stop conditions. If a directory contains both a target marker and a stop marker, the target marker matches and the method returns it, then traversal stops.

This ordering matters for the common case where project root and VCS root share the same directory (for example, a repo root with both `.git` and `go.mod`).

#### Empty results

`FindUp` and related single-result methods return `found=false` when the walker reaches the filesystem root (or a stop marker) without finding any marker. `FindAllUp` and related methods return an empty slice.

### Usage examples

#### Find project root by marker

```go
// Go module root
dir, _, found := dirs.FindUp(cwd, "go.mod")

// Rust package root
dir, _, found := dirs.FindUp(cwd, "Cargo.toml")

// Node package root
dir, _, found := dirs.FindUp(cwd, "package.json")
```

#### Find project root with VCS boundary

```go
// Find go.mod but don't leave the git repo
dir, _, found := dirs.FindUpUntil(cwd, []string{"go.mod"}, []string{".git"})
```

#### Find workspace root with content inspection

Cargo workspaces require checking file content:

```go
dir, _, found := dirs.FindUpFunc(cwd, []string{"Cargo.toml"}, func(path string) bool {
    content, err := os.ReadFile(path)
    if err != nil {
        return false
    }
    return bytes.Contains(content, []byte("[workspace]"))
})
```

#### Cascading configuration files

Some tools merge configuration from many files (nearest takes precedence):

```go
// Find all .myconfig files from cwd to repo root
matches := dirs.FindAllUpUntil(cwd, []string{".myconfig"}, []string{".git"})

// Load in reverse order so nearest overrides farthest
for i := len(matches) - 1; i >= 0; i-- {
    mergeConfig(matches[i].Path())
}
```

#### Stop traversal based on file content

EditorConfig stops at files containing `root=true`. Since the stop condition requires content inspection, handle it in caller code:

```go
matches := dirs.FindAllUp(cwd, ".editorconfig")

var configs []string
for _, m := range matches {
    configs = append(configs, m.Path())
    if hasRootTrue(m.Path()) {
        break
    }
}
```

#### Prioritized marker search

Go workspaces take precedence over modules:

```go
// Prefer go.work over go.mod
if dir, _, found := dirs.FindUp(cwd, "go.work"); found {
    return dir, nil
}
if dir, _, found := dirs.FindUp(cwd, "go.mod"); found {
    return dir, nil
}
return "", ErrProjectNotFound
```

#### Tool configuration search

Find a tool's configuration, preferring project-local over global:

```go
// Check project first (bounded by VCS), then home directory
if path, _, found := dirs.FindUpUntil(cwd, []string{".mytool.yaml"}, []string{".git"}); found {
    return path, nil
}
// Fall back to user config directory
return dirs.UserConfigPath(".mytool.yaml"), nil
```

### Design rationale

#### Primitives over policies

These methods handle traversal mechanics. They do not encode opinions about what "project" or "workspace" means. Different ecosystems have incompatible definitions: developers use a Go workspace (`go.work`) for local development and typically exclude it from version control, while a Cargo workspace (`[workspace]` in `Cargo.toml`) forms part of the project structure. The primitives let each tool define its own semantics.

#### No automatic content inspection

The library does not read file contents. Callers provide match functions when content matters. This keeps the library focused on filesystem traversal and avoids dependencies on TOML, YAML, or JSON parsers.

#### Markers are existence-based

A marker matches if it exists. The library does not distinguish between files and directories for matching purposes. If you need to match only files or only directories, use a match function:

```go
// Only match if .git is a directory (not a gitlink file)
dirs.FindUpFunc(cwd, []string{".git"}, func(path string) bool {
    info, err := os.Stat(path)
    return err == nil && info.IsDir()
})
```

#### Stop after match on combined marker/stop directories

When a directory contains both a target marker and a stop marker, the target matches before traversal stops. This handles the common case of a repo root containing both `.git` and the project marker. Without this ordering, tools would fail to find projects at repo roots when using `.git` as a stop marker.

#### Symlink handling

The traversal methods use logical path semantics rather than physical path semantics. When walking up from a start directory, the methods traverse the path as the user sees it, not the path after resolving symlinks.

Consider a user in `~/projects/myapp` where `projects` is a symlink to `/mnt/data/projects`. Logical traversal visits `~/projects/myapp` → `~/projects` → `~/` → `/`. Physical traversal would instead visit `/mnt/data/projects/myapp` → `/mnt/data/projects` → `/mnt/data` → `/mnt` → `/`.

Logical semantics match user expectations and align with how most tools behave (git, cargo, npm all use logical paths). Users navigate by their mental model of the filesystem, not where bytes are stored. Physical traversal would produce surprising results when symlinks cross mount points or reference paths outside the user's home directory.

Implementation uses `filepath.Clean()` on the start path without calling `filepath.EvalSymlinks()`.

For marker detection, the methods use `os.Stat()` which follows symlinks. A symlinked marker matches if its target exists. This handles common cases like git worktrees and submodules, where `.git` may be a file (gitlink) pointing elsewhere rather than a directory. Broken symlinks do not match since their targets do not exist.

Match functions in the `*Func` variants receive the logical path to the marker (the path found during traversal), not a resolved path. This keeps behavior consistent with what callers would see from `ls` in that directory. Callers needing the physical path can call `filepath.EvalSymlinks()` themselves.

#### Testing with `FakeDirs`

`FakeDirs` supports the `FindUp*` methods through its existing marker system. Use `SetExisting()` to declare which markers exist at which paths:

```go
fake := toolpaths.NewFakeDirs("/tmp/test-app")

// Simulate a project structure
fake.SetExisting("/home/user/projects/myapp/go.mod")
fake.SetExisting("/home/user/projects/myapp/.git")
fake.SetExisting("/home/user/projects/.git")  // parent repo

// Test finding the nearest go.mod
dir, marker, found := fake.FindUp("/home/user/projects/myapp/cmd/server", "go.mod")
// dir = "/home/user/projects/myapp", marker = "go.mod", found = true

// Test with stop marker
dir, _, found = fake.FindUpUntil(
    "/home/user/projects/myapp/cmd/server",
    []string{"go.mod"},
    []string{".git"},
)
// Stops at myapp because it contains both go.mod and .git
```

The fake implementation walks the logical path and checks `SetExisting()` entries for each marker at each directory level. This lets tests verify traversal logic without creating real directory structures.

## Comparison with alternatives

### `github.com/adrg/xdg`

The `adrg/xdg` library provides platform-native defaults on macOS (Library directories) and Windows (Known Folders) when XDG environment variables are not set. The key architectural differences from toolpaths are:

- Global singleton API with package-level variables (`xdg.ConfigHome`, `xdg.DataHome`) versus toolpaths' instance-based API with configuration options
- Returns base directories without app names; apps must construct subdirectory paths manually
- No built-in support for app versioning or author prefixes on Windows
- No fallback search paths for migration scenarios; `*Dirs` variables contain platform-specific search paths but not XDG-to-native migration fallbacks

toolpaths handles app-name-scoped paths, versioning, and layered directory search automatically. Use `XDGOnAllPlatforms` if you want XDG paths as the primary location on all platforms regardless of native conventions.

### Python `platformdirs`

This library draws inspiration from the Python `platformdirs` package. Key differences:

- Go version uses `*Dir()` / `*Dirs()` naming rather than Python's `user_*_dir` / `user_*_path`. The `*Dir()` methods return the primary directory while `*Dirs()` returns a slice including fallbacks.
- Go version includes XDG fallbacks by default for migration support (controlled by `IncludeXDGFallbacks`)
- Go version provides find utilities as methods (`FindConfigFile`, `AllConfigPaths`, `ExistingConfigFiles`) rather than requiring manual iteration
- Go version provides `Ensure*Dir()` methods to create directories with appropriate permissions

### `os.UserConfigDir`, `os.UserCacheDir`

Go's standard library provides `os.UserConfigDir()` and `os.UserCacheDir()` which return platform-appropriate base directories. These functions:

- Do not append app names
- Do not support versioning
- Do not provide data, state, log, or runtime directories
- Do not support system-wide directories
- Do not support search paths or fallbacks

toolpaths builds on the same platform conventions while providing a complete solution for app directory management.
