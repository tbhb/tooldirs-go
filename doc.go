// Package toolpaths provides platform-appropriate directory paths for
// application configuration, data, cache, state, logs, and runtime files.
//
// It implements the XDG Base Directory Specification on Linux/BSD, and uses
// native conventions on macOS and Windows. The library supports:
//
//   - User-specific directories (config, data, cache, state, log, runtime)
//   - System-wide directories (config, data)
//   - XDG environment variable overrides on all platforms (opt-in)
//   - App-specific environment variable overrides
//   - File path resolution helpers
//   - Find utilities for layered configuration
package toolpaths
