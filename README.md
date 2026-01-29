# tooldirs

Platform-appropriate directory paths for Go apps.

tooldirs resolves where configuration files, app data, caches, logs, and runtime files should live on Linux, macOS, and Windows. It follows each platform's native conventions while providing fallback paths for migration scenarios.

## Installation

```bash
go get github.com/tbhb/tooldirs-go
```

Requires Go 1.25 or later.

## Quick start

```go
package main

import (
    "fmt"
    "github.com/tbhb/tooldirs-go"
)

func main() {
    dirs := tooldirs.New("myapp")

    // Write config files here
    fmt.Println(dirs.UserConfigDir())
    // Linux:   ~/.config/myapp
    // macOS:   ~/Library/Application Support/myapp
    // Windows: C:\Users\<user>\AppData\Local\myapp

    // Find a config file (user dirs first, then system)
    if path, found := dirs.FindConfigFile("config.yaml"); found {
        fmt.Printf("Found config at: %s\n", path)
    }

    // Create the config directory if needed
    configDir, err := dirs.EnsureUserConfigDir()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Config dir ready: %s\n", configDir)
}
```

## Directory types

The library supports six directory types:

| Type      | Purpose                                           | Example files                       |
| --------- | ------------------------------------------------- | ----------------------------------- |
| `Config`  | User-editable configuration files                 | config.yaml, settings.json          |
| `Data`    | App data that is not configuration                | databases, plugins, themes          |
| `Cache`   | Non-essential data you can regenerate             | HTTP cache, thumbnails              |
| `State`   | Data that should persist but is not configuration | recently opened files, undo history |
| `Log`     | App log files                                     | app.log, error.log                  |
| `Runtime` | Ephemeral files that should not survive a reboot  | Unix sockets, PID files             |

## Platform behavior

### Linux, FreeBSD, OpenBSD

Implements the XDG Base Directory Specification:

| Type      | Default                    |
| --------- | -------------------------- |
| `Config`  | `~/.config/{app}`          |
| `Data`    | `~/.local/share/{app}`     |
| `Cache`   | `~/.cache/{app}`           |
| `State`   | `~/.local/state/{app}`     |
| `Log`     | `~/.local/state/{app}/log` |
| `Runtime` | `$XDG_RUNTIME_DIR/{app}`   |

The library respects XDG environment variables (`XDG_CONFIG_HOME`, etc.) when set.

### macOS

Uses Library directories:

| Type      | Location                              |
| --------- | ------------------------------------- |
| `Config`  | `~/Library/Application Support/{app}` |
| `Data`    | `~/Library/Application Support/{app}` |
| `State`   | `~/Library/Application Support/{app}` |
| `Cache`   | `~/Library/Caches/{app}`              |
| `Log`     | `~/Library/Logs/{app}`                |
| `Runtime` | `$TMPDIR/{app}`                       |

`Config`, `Data`, and `State` resolve to the same location since macOS does not distinguish between them at the filesystem level.

### Windows

Uses Known Folders:

| Type      | Location (Roaming=false)       | Location (Roaming=true)        |
| --------- | ------------------------------ | ------------------------------ |
| `Config`  | `%LOCALAPPDATA%\{app}`         | `%APPDATA%\{app}`              |
| `Data`    | `%LOCALAPPDATA%\{app}`         | `%APPDATA%\{app}`              |
| `State`   | `%LOCALAPPDATA%\{app}`         | `%APPDATA%\{app}`              |
| `Cache`   | `%LOCALAPPDATA%\{app}\cache`   | `%LOCALAPPDATA%\{app}\cache`   |
| `Log`     | `%LOCALAPPDATA%\{app}\log`     | `%LOCALAPPDATA%\{app}\log`     |
| `Runtime` | `%LOCALAPPDATA%\{app}\runtime` | `%LOCALAPPDATA%\{app}\runtime` |

Cache, log, and runtime always use LOCALAPPDATA regardless of the Roaming setting.

## API overview

### Single path vs many paths

`*Dir()` methods return a single path for writing new files:

```go
configDir := dirs.UserConfigDir()  // Primary config directory
```

`*Dirs()` methods return a slice of paths in priority order. The first element is the primary directory; the rest are fallback locations for reading:

```go
// On macOS, returns [~/Library/Application Support/myapp, ~/.config/myapp]
// The second path enables migration from XDG paths to native paths
configDirs := dirs.UserConfigDirs()
```

### Path helpers

Build paths within a directory:

```go
settingsPath := dirs.UserConfigPath("settings.json")
// ~/Library/Application Support/myapp/settings.json
```

### Find utilities

Find files across user and system directories:

```go
// First existing file in priority order
path, found := dirs.FindConfigFile("config.yaml")

// All candidate paths (without checking existence)
paths := dirs.AllConfigPaths("config.yaml")

// All paths where the file actually exists
existing := dirs.ExistingConfigFiles("config.yaml")
```

### Ensure utilities

Create directories with mode 0700 if they do not exist:

```go
dir, err := dirs.EnsureUserConfigDir()
dir, err := dirs.EnsureUserDataDir()
dir, err := dirs.EnsureUserCacheDir()
dir, err := dirs.EnsureUserStateDir()
dir, err := dirs.EnsureUserLogDir()
```

### Diagnostics

Debug path resolution:

```go
info := dirs.Diagnose()
fmt.Print(info.String())

env := dirs.DiagnoseEnv()
fmt.Print(dirs.DiagnoseEnvString())
```

## Configuration

```go
dirs := tooldirs.NewWithConfig(tooldirs.Config{
    // Required: application name used as directory name
    AppName: "myapp",

    // Optional: vendor name (Windows only, creates {AppAuthor}\{AppName})
    AppAuthor: "MyCompany",

    // Optional: version subdirectory for major version separation
    Version: "2",

    // Optional: use %APPDATA% instead of %LOCALAPPDATA% (Windows only)
    Roaming: true,

    // Optional: use XDG paths on all platforms
    XDGOnAllPlatforms: true,

    // Optional: include XDG fallbacks in *Dirs() on non-XDG platforms
    // Default is true for migration support
    IncludeXDGFallbacks: &falseVal,

    // Optional: app-specific environment variable overrides
    EnvOverrides: &tooldirs.EnvOverrides{
        AppendAppName: true,
        UserConfig:    "MYAPP_CONFIG_HOME",
        UserData:      "MYAPP_DATA_HOME",
    },
})
```

## Testing

The `FakeDirs` type implements the `Dirs` interface for testing without filesystem or environment interaction:

```go
import "github.com/tbhb/tooldirs-go"

// Create a fake with all paths under a base directory
fake := tooldirs.NewFakeDirs("/tmp/test-app")
// fake.UserConfigDir() returns "/tmp/test-app/config"
// fake.UserDataDir() returns "/tmp/test-app/data"

// Or use a temp directory with automatic cleanup
fake, cleanup := tooldirs.NewFakeDirsWithTempDir("test")
defer cleanup()

// Control file existence for find utilities
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
    dirs tooldirs.Dirs
}

func NewApp(dirs tooldirs.Dirs) *App {
    return &App{dirs: dirs}
}

// In production
app := NewApp(tooldirs.New("myapp"))

// In tests
fake := tooldirs.NewFakeDirs("/tmp/test")
app := NewApp(fake)
```

## Related projects

- [`adrg/xdg`](https://github.com/adrg/xdg) - XDG Base Directory Specification for Go. Provides platform-native defaults on macOS and Windows but exposes a global singleton API returning base directories without app names. Apps must construct subdirectory paths manually.
- [`platformdirs`](https://github.com/platformdirs/platformdirs) - Python library that inspired this project
- [`os.UserConfigDir`](https://pkg.go.dev/os#UserConfigDir) - Go standard library (no app name, no search paths)

## AI disclosure

Claude Code helped develop this project. The maintainer reviewed and tested all code.

## License

MIT License. See [LICENSE](MIT-LICENSE) for details.
