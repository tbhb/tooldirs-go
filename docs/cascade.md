# Cascade design document

This document describes the design decisions, behavioral semantics, and configuration options for the cascade functionality in toolpaths.

## Overview

The `Cascade` type provides layered path resolution for CLI tools and applications. It answers the question: "given a file or directory name, where should I look for it and in what order?"

Configuration and data are modeled as a cascade of scopes with defined precedence. Reads search scopes from highest to lowest priority, returning the first match (or collecting all matches for merging). Writes target a specific scope. This pattern appears across many tools: git config (local → global → system), EditorConfig (nearest → farthest), npm (project → user → global), and enterprise tools with managed policies.

`Cascade` builds on top of `Dirs` (platform directory resolution) and the `FindUp` methods (directory traversal). Most CLI tools use `Cascade` directly; simpler tools or libraries that just need platform directories use `Dirs`.

## Core concepts

### Scopes

A scope represents a layer in the cascade hierarchy. Each scope has a name, a priority (lower numbers = higher priority for reads), a base path, and flags controlling behavior.

Built-in scopes cover common patterns:

| Scope             | Typical base path                                 | Use case                          |
| ----------------- | ------------------------------------------------- | --------------------------------- |
| `ScopeManagedReq` | `/etc/{app}/managed/required/`                    | Enterprise policies, applied last |
| `ScopeSystem`     | `/etc/{app}/` or platform equivalent              | System-wide defaults              |
| `ScopeUser`       | `~/.config/{app}/` or platform equivalent         | User preferences                  |
| `ScopeProject`    | `{project}/.{app}/` or `{project}/.config/{app}/` | Project-specific config           |
| `ScopeLocal`      | `{project}/.{app}.local/`                         | Git-ignored local overrides       |
| `ScopeManagedRec` | `/etc/{app}/managed/recommended/`                 | Enterprise defaults, applied first|

The default cascade for a typical CLI tool is: Local → Project → User → System. Enterprise tools can add managed layers that sandwich user choices.

### Cascade ordering

Read operations search scopes in priority order. The first scope containing the requested file wins for single-file lookups. For merge scenarios, all matching files are returned in priority order.

Write operations target a specific scope. The library does not implicitly choose where to write.

### Project detection

Many tools scope configuration to a project root. `Cascade` finds project roots by walking up the directory tree looking for markers (VCS directories, config files, lock files) using the `FindUp` methods.

Config-dir awareness handles an edge case: when the working directory is inside the config directory itself (e.g., `~/.config/myapp/policies.d/`), the library recognizes this and returns the project root directly rather than walking up and potentially finding a parent project.

### Project directory structure

Within a project, tools typically need multiple directory types beyond just config: cache, data, logs, state. Each directory type can have multiple location patterns (e.g., `.myapp/cache` vs `.myapp-cache`), named subdirectories, and known files.

The library models this with three layers of extensibility:

1. Project directory types (config, cache, data, log, state)
2. Subdirectories within those types (policies.d/, hooks/, plugins/)
3. Known files with naming variants

### File variants

Tools often support multiple file formats or naming conventions. A tool might accept `config.yaml`, `config.yml`, `config.json`, or `.config.yaml`. The library checks variants in order, returning the first match.

## Configuration

### Cascade configuration

```go
type CascadeConfig struct {
    // AppName is required. Used for directory names within each scope.
    AppName string

    // Dirs provides platform directory resolution for User and System scopes.
    // If nil, creates a Dirs with default configuration.
    Dirs *Dirs

    // Scopes defines the cascade hierarchy. If nil, uses DefaultScopes().
    // Scopes are searched in priority order (lowest priority number first).
    Scopes []ScopeConfig

    // ProjectMarkers define filesystem entries that indicate a project root.
    // Defaults to VCS markers (.git, .hg, .svn, .bzr) plus app-specific
    // directories (.{app}, .config/{app}).
    ProjectMarkers []string

    // ConfigDirPatterns for config-dir-aware project detection.
    // When cwd is inside a matching pattern, the library returns the
    // parent as the project root without walking up.
    // Defaults to [".{app}", ".config/{app}"].
    ConfigDirPatterns []string

    // StopMarkers halt upward traversal during project detection.
    // Useful for monorepo boundaries. Defaults to empty (traverse to root).
    StopMarkers []string

    // ProjectDirs defines directory types within the project.
    // If nil, uses DefaultProjectDirs().
    ProjectDirs map[string]ProjectDirConfig

    // ProjectFiles defines known files with their variants and locations.
    // If nil, no known files are configured.
    ProjectFiles map[string]ProjectFileConfig

    // FileVariants controls which naming patterns to check when resolving files.
    // Defaults to VariantBoth (checks both "file.yaml" and ".file.yaml").
    FileVariants VariantStyle

    // Getwd provides the current working directory. If nil, uses os.Getwd.
    // Useful for testing.
    Getwd func() (string, error)
}
```

### Scope configuration

```go
type ScopeConfig struct {
    // Name identifies this scope (e.g., "user", "project", "system").
    Name string

    // Priority determines search order. Lower numbers = higher priority.
    // Reads check scopes in ascending priority order.
    Priority int

    // Writable indicates whether this scope accepts writes.
    // Managed/system scopes are typically read-only.
    Writable bool

    // BasePath returns the base directory for this scope.
    // Receives the Cascade instance for access to Dirs and project root.
    BasePath func(c *Cascade) (string, error)

    // Subdirs are checked within the base path. For example, a project
    // scope might check both ".myapp" and ".config/myapp" subdirectories.
    // If empty, files are resolved directly under BasePath.
    Subdirs []string
}
```

### Project directory configuration

```go
// ProjectDirConfig defines a directory type at the project level.
type ProjectDirConfig struct {
    // Patterns are paths relative to project root, checked in order.
    // First existing one wins for reads; first one is used for writes.
    // Supports {app} placeholder for the app name.
    Patterns []string

    // Subdirs defines named subdirectories within this directory type.
    // Each name maps to path variants to check (relative to the dir pattern).
    Subdirs map[string][]string
}

// ProjectFileConfig defines a known file for convenient resolution.
type ProjectFileConfig struct {
    // Variants are filename variants to check, in order.
    Variants []string

    // Dir is which project directory type this file lives in.
    // Empty string means project root directly.
    Dir string

    // Subdirs to also check within Dir (in addition to Dir root).
    // Useful for files that might live in either .myapp/config.yaml
    // or .myapp/conf.d/config.yaml.
    Subdirs []string
}
```

### Variant styles

```go
type VariantStyle int

const (
    // VariantBoth checks "name.ext" and ".name.ext"
    VariantBoth VariantStyle = iota
    // VariantDotted checks only ".name.ext"
    VariantDotted
    // VariantPlain checks only "name.ext"
    VariantPlain
    // VariantCustom uses a provided FileResolver
    VariantCustom
)
```

## Primary interface

### Construction

```go
// Cascade provides layered path resolution across multiple scopes.
type Cascade struct {
    // contains filtered or unexported fields
}

// NewCascade creates a Cascade with the given app name and options.
func NewCascade(appName string, opts ...CascadeOption) (*Cascade, error)

// Dirs returns the underlying platform directory resolver.
func (c *Cascade) Dirs() *Dirs
```

### Project detection

```go
// ProjectRoot returns the detected project root, or empty string if none found.
// Results are cached after first call.
func (c *Cascade) ProjectRoot() (string, error)

// IsInsideProject returns true if a project root was detected.
func (c *Cascade) IsInsideProject() bool
```

### Scope access

```go
// ScopeDir returns the base directory for the named scope.
func (c *Cascade) ScopeDir(name string) (string, error)

// Scopes returns all configured scopes in priority order.
func (c *Cascade) Scopes() []ScopeConfig
```

### Project directory types

```go
// ProjectDir returns the primary and alternate paths for a project directory type.
// The primary is the first existing directory; alternates are other existing locations.
// For writes, use DefaultProjectDir which returns the first pattern regardless of existence.
func (c *Cascade) ProjectDir(name string) (primary string, alternates []string, err error)

// Convenience methods for common directory types:
func (c *Cascade) ProjectConfigDir() (string, []string, error)
func (c *Cascade) ProjectCacheDir() (string, []string, error)
func (c *Cascade) ProjectDataDir() (string, []string, error)
func (c *Cascade) ProjectStateDir() (string, []string, error)
func (c *Cascade) ProjectLogDir() (string, []string, error)

// DefaultProjectDir returns the default path for a project directory type.
// This is the first pattern, used for writes. Does not check existence.
func (c *Cascade) DefaultProjectDir(name string) (string, error)
```

### Project subdirectories

```go
// ProjectSubdir returns paths for a named subdirectory within a project directory type.
func (c *Cascade) ProjectSubdir(dirType, subdirName string) (primary string, alternates []string, err error)

// DefaultProjectSubdir returns the default path for a subdirectory.
func (c *Cascade) DefaultProjectSubdir(dirType, subdirName string) (string, error)
```

### Project file resolution

```go
// ResolveProjectFile finds a file within a project directory type.
func (c *Cascade) ResolveProjectFile(dirType, basename string) (primary string, alternates []string, err error)

// ResolveProjectFileIn finds a file within a specific subdirectory of a project directory type.
func (c *Cascade) ResolveProjectFileIn(dirType, subdir, basename string) (primary string, alternates []string, err error)

// ResolveKnownFile resolves a file by its configured name (from ProjectFiles).
func (c *Cascade) ResolveKnownFile(name string) (primary string, alternates []string, err error)

// DefaultProjectFile returns the default path for a file in a project directory type.
func (c *Cascade) DefaultProjectFile(dirType, basename string) (string, error)

// DefaultKnownFile returns the default path for a known file.
func (c *Cascade) DefaultKnownFile(name string) (string, error)
```

### Cross-scope file resolution

```go
// ResolveFile finds a file across scopes. Returns the highest-priority
// existing file, any additional existing files in other scopes, and an error.
// If no file exists in any scope, primary is empty and alternates is nil.
func (c *Cascade) ResolveFile(basename string, scopes ...string) (primary string, alternates []string, err error)

// ResolveFileIn is like ResolveFile but searches only the specified scopes.
func (c *Cascade) ResolveFileIn(basename string, scopes []string) (primary string, alternates []string, err error)

// AllFilePaths returns all candidate paths for a file across scopes,
// whether or not they exist. Useful for documentation or debugging.
func (c *Cascade) AllFilePaths(basename string, scopes ...string) []string

// ExistingFiles returns paths to all existing instances of a file
// across scopes, in priority order.
func (c *Cascade) ExistingFiles(basename string, scopes ...string) []string
```

### Cross-scope directory resolution

```go
// ResolveDir finds a directory across scopes.
func (c *Cascade) ResolveDir(basename string, scopes ...string) (primary string, alternates []string, err error)

// ResolveDirIn is like ResolveDir but searches only the specified scopes.
func (c *Cascade) ResolveDirIn(basename string, scopes []string) (primary string, alternates []string, err error)

// AllDirPaths returns all candidate paths for a directory across scopes.
func (c *Cascade) AllDirPaths(basename string, scopes ...string) []string

// ExistingDirs returns paths to all existing instances of a directory.
func (c *Cascade) ExistingDirs(basename string, scopes ...string) []string
```

### Write targeting

```go
// DefaultPath returns the path where a new file should be written
// in the specified scope. Does not check existence.
func (c *Cascade) DefaultPath(basename string, scope string) (string, error)

// DefaultDir returns the directory path in the specified scope.
func (c *Cascade) DefaultDir(basename string, scope string) (string, error)
```

### Path joining

```go
// JoinProject joins path elements to the project root.
func (c *Cascade) JoinProject(elem ...string) (string, error)

// JoinProjectDir joins path elements to a project directory type.
func (c *Cascade) JoinProjectDir(dirType string, elem ...string) (string, error)

// JoinScope joins path elements to a scope's base directory.
func (c *Cascade) JoinScope(scope string, elem ...string) (string, error)
```

## Functional options

```go
// WithCascadeDirs provides a custom Dirs for platform directory resolution.
func WithCascadeDirs(dirs *Dirs) CascadeOption

// WithScopes replaces the default scope configuration.
func WithScopes(scopes ...ScopeConfig) CascadeOption

// WithProjectMarkers sets custom project root markers.
func WithProjectMarkers(markers ...string) CascadeOption

// WithConfigDirPatterns sets patterns for config-dir-aware detection.
func WithConfigDirPatterns(patterns ...string) CascadeOption

// WithStopMarkers sets markers that halt upward traversal.
func WithStopMarkers(markers ...string) CascadeOption

// WithProjectDirs configures project directory types.
func WithProjectDirs(dirs map[string]ProjectDirConfig) CascadeOption

// WithProjectFiles configures known project files.
func WithProjectFiles(files map[string]ProjectFileConfig) CascadeOption

// WithFileVariants controls file naming pattern matching.
func WithFileVariants(style VariantStyle) CascadeOption

// WithFileResolver provides custom file variant resolution.
func WithFileResolver(resolver FileResolver) CascadeOption

// WithGetwd provides a custom working directory function (for testing).
func WithGetwd(fn func() (string, error)) CascadeOption

// WithoutProjectDetection disables project scope entirely.
// Useful for system daemons or tools without a project concept.
func WithoutProjectDetection() CascadeOption

// WithManagedLayers adds managed/required and managed/recommended scopes
// for enterprise policy support.
func WithManagedLayers() CascadeOption
```

## Default configurations

### DefaultProjectDirs

The default project directory structure for CLI tools:

```go
func DefaultProjectDirs(appName string) map[string]ProjectDirConfig {
    dotted := "." + appName
    xdg := ".config/" + appName

    return map[string]ProjectDirConfig{
        "config": {
            Patterns: []string{xdg, dotted},
        },
        "cache": {
            Patterns: []string{dotted + "/cache"},
        },
        "data": {
            Patterns: []string{dotted + "/data"},
        },
        "state": {
            Patterns: []string{dotted + "/state"},
        },
        "log": {
            Patterns: []string{dotted + "/log"},
        },
    }
}
```

### DefaultScopes

The default scope configuration for CLI tools with project support:

```go
func DefaultScopes() []ScopeConfig {
    return []ScopeConfig{
        {
            Name:     "local",
            Priority: 10,
            Writable: true,
            BasePath: localBasePath,
            // local files live directly in project root
        },
        {
            Name:     "project",
            Priority: 20,
            Writable: true,
            BasePath: projectBasePath,
            Subdirs:  []string{".config/{app}", ".{app}"},
        },
        {
            Name:     "user",
            Priority: 30,
            Writable: true,
            BasePath: userBasePath,
            // uses Dirs.UserConfigDir()
        },
        {
            Name:     "system",
            Priority: 40,
            Writable: false,
            BasePath: systemBasePath,
            // uses Dirs.SystemConfigDir()
        },
    }
}
```

### ManagedScopes

For enterprise tools with managed policy layers:

```go
func ManagedScopes() []ScopeConfig {
    base := DefaultScopes()
    return append([]ScopeConfig{
        {
            Name:     "managed-required",
            Priority: 0, // Highest priority, applied last
            Writable: false,
            BasePath: managedRequiredBasePath,
        },
    }, append(base, ScopeConfig{
        Name:     "managed-recommended",
        Priority: 100, // Lowest priority, applied first
        Writable: false,
        BasePath: managedRecommendedBasePath,
    })...)
}
```

This creates the sandwich pattern where enterprise policies wrap user preferences:

1. managed-required (priority 0) - enforced, user cannot override
2. local (priority 10)
3. project (priority 20)
4. user (priority 30)
5. system (priority 40)
6. managed-recommended (priority 100) - defaults, user can override

## File resolution

### Variant matching

The library generates candidate filenames from a basename according to the configured `VariantStyle`:

| Style           | Input         | Candidates                    |
| --------------- | ------------- | ----------------------------- |
| `VariantBoth`   | `config.yaml` | `config.yaml`, `.config.yaml` |
| `VariantDotted` | `config.yaml` | `.config.yaml`                |
| `VariantPlain`  | `config.yaml` | `config.yaml`                 |

For more complex patterns, implement `FileResolver`:

```go
type FileResolver interface {
    // Candidates returns all filenames to check for a given basename.
    Candidates(basename string) []string
}

// Built-in resolvers
var (
    // YAMLResolver checks .yaml and .yml extensions
    YAMLResolver FileResolver

    // JSONResolver checks .json extension
    JSONResolver FileResolver

    // MultiFormatResolver checks .yaml, .yml, .json, .toml
    MultiFormatResolver FileResolver
)
```

### Resolution algorithm

For `ResolveFile(basename, scopes...)`:

1. If scopes is empty, use all configured scopes in priority order
2. For each scope in priority order:
   a. Get the scope's base path
   b. For each subdir in the scope (or root if no subdirs):
      c. For each candidate filename from the resolver:
         d. Check if the file exists
         e. If yes and this is the first match, record as primary
         f. If yes and not first, append to alternates
3. Return primary, alternates, nil (or empty, nil, nil if nothing found)

For `ResolveProjectFile(dirType, basename)`:

1. Get the project directory type configuration
2. For each pattern in the directory type:
   a. For each candidate filename from the resolver:
      b. Check if the file exists
      c. If yes and this is the first match, record as primary
      d. If yes and not first, append to alternates
3. Return primary, alternates, nil

## Usage examples

### Basic CLI tool

```go
cascade, err := toolpaths.NewCascade("mytool")
if err != nil {
    log.Fatal(err)
}

// Find config file (searches local → project → user → system)
configPath, _, err := cascade.ResolveFile("config.yaml")
if err != nil {
    log.Fatal(err)
}
if configPath != "" {
    loadConfig(configPath)
}

// Access project directories
cacheDir, _, _ := cascade.ProjectCacheDir()
logDir, _, _ := cascade.ProjectLogDir()

// Write new config to user scope
userConfigPath, err := cascade.DefaultPath("config.yaml", "user")
if err != nil {
    log.Fatal(err)
}
writeConfig(userConfigPath, defaultConfig)
```

### Complex tool with custom project structure

This example shows configuration similar to agenthooks:

```go
cascade, _ := toolpaths.NewCascade("agenthooks",
    toolpaths.WithProjectDirs(map[string]toolpaths.ProjectDirConfig{
        "config": {
            Patterns: []string{".config/agenthooks", ".agenthooks"},
            Subdirs: map[string][]string{
                "policies":       {"policies.d"},
                "local-policies": {"policies.local.d"},
            },
        },
        "cache": {
            Patterns: []string{".agenthooks/cache", ".agenthooks-cache"},
        },
        "log": {
            Patterns: []string{".agenthooks/log", ".agenthooks/logs"},
        },
    }),
    toolpaths.WithProjectFiles(map[string]toolpaths.ProjectFileConfig{
        "config": {
            Variants: []string{"agenthooks.yaml", ".agenthooks.yaml"},
            Dir:      "config",
        },
        "local-config": {
            Variants: []string{"agenthooks.local.yaml", ".agenthooks.local.yaml"},
            Dir:      "config",
        },
        "policies": {
            Variants: []string{"agenthooks-policies.yaml", ".agenthooks-policies.yaml"},
            Dir:      "config",
        },
        "local-policies": {
            Variants: []string{"agenthooks-policies.local.yaml", ".agenthooks-policies.local.yaml"},
            Dir:      "config",
        },
    }),
)

// Resolve known files by name
configFile, alts, _ := cascade.ResolveKnownFile("config")
policiesFile, _, _ := cascade.ResolveKnownFile("policies")

// Access named subdirectories
policiesDir, _, _ := cascade.ProjectSubdir("config", "policies")
localPoliciesDir, _, _ := cascade.ProjectSubdir("config", "local-policies")

// Get default paths for creating new files
defaultConfig, _ := cascade.DefaultKnownFile("config")
defaultPoliciesDir, _ := cascade.DefaultProjectSubdir("config", "policies")
```

### Merging layered configuration

```go
cascade, _ := toolpaths.NewCascade("mytool")

// Get all config files in priority order
configs := cascade.ExistingFiles("config.yaml")

// Load in reverse order so highest priority wins
merged := make(map[string]any)
for i := len(configs) - 1; i >= 0; i-- {
    cfg := loadYAML(configs[i])
    mergeMaps(merged, cfg)
}
```

### Tool without project concept

```go
// System daemon with only user and system config
cascade, _ := toolpaths.NewCascade("mydaemon",
    toolpaths.WithoutProjectDetection(),
    toolpaths.WithScopes(
        toolpaths.ScopeConfig{Name: "user", Priority: 10, Writable: true, BasePath: userBase},
        toolpaths.ScopeConfig{Name: "system", Priority: 20, Writable: false, BasePath: systemBase},
    ),
)
```

### Enterprise tool with managed policies

```go
cascade, _ := toolpaths.NewCascade("corptool",
    toolpaths.WithManagedLayers(),
)

// Required policies always apply (cannot be overridden)
requiredPolicy, _, _ := cascade.ResolveFileIn("policy.yaml", []string{"managed-required"})

// User preferences with recommended defaults
userPolicy, alternates, _ := cascade.ResolveFile("policy.yaml")
// alternates may include managed-recommended as fallback
```

### EditorConfig-style cascading

```go
cascade, _ := toolpaths.NewCascade("myeditor",
    toolpaths.WithProjectMarkers(".myeditor"),
)

// Collect all config files from cwd to project root
// Stop early if we hit one with root=true
matches := cascade.Dirs().FindAllUp(cwd, ".myeditor")

var configs []string
for _, m := range matches {
    configs = append(configs, m.Path())
    if hasRootDirective(m.Path()) {
        break
    }
}

// Apply in reverse order (farthest first, nearest wins)
for i := len(configs) - 1; i >= 0; i-- {
    applyConfig(configs[i])
}
```

### Custom file resolution

```go
// Tool that accepts config.yaml, config.yml, config.json, or config.toml
cascade, _ := toolpaths.NewCascade("polytool",
    toolpaths.WithFileResolver(toolpaths.MultiFormatResolver),
)

// Finds first existing file regardless of extension
configPath, _, _ := cascade.ResolveFile("config")
```

### Monorepo with workspace boundaries

```go
cascade, _ := toolpaths.NewCascade("monotool",
    toolpaths.WithStopMarkers("pnpm-workspace.yaml", "lerna.json"),
)

// Project detection stops at workspace root, won't escape to parent monorepo
projectRoot, _ := cascade.ProjectRoot()
```

## Testing

### FakeCascade

The library provides a test double for controlling path resolution without touching the filesystem:

```go
fake := toolpaths.NewFakeCascade("/tmp/test")

// Configure scope directories
fake.SetScopeDir("user", "/tmp/test/user-config")
fake.SetScopeDir("project", "/tmp/test/myproject/.mytool")

// Configure project directories
fake.SetProjectDir("config", "/tmp/test/myproject/.mytool")
fake.SetProjectDir("cache", "/tmp/test/myproject/.mytool/cache")

// Configure existing files
fake.SetExisting("/tmp/test/user-config/config.yaml")
fake.SetExisting("/tmp/test/myproject/.mytool/local.yaml")

// Configure project root
fake.SetProjectRoot("/tmp/test/myproject")

// Use in tests
primary, alternates, err := fake.ResolveFile("config.yaml")
configDir, _, _ := fake.ProjectConfigDir()
```

### Integration with FakeDirs

For tests that need to control both platform directories and cascade behavior:

```go
fakeDirs := toolpaths.NewFakeDirs("/tmp/test")
cascade, _ := toolpaths.NewCascade("mytool",
    toolpaths.WithCascadeDirs(fakeDirs),
    toolpaths.WithGetwd(func() (string, error) {
        return "/tmp/test/projects/myapp", nil
    }),
)
```

## Design rationale

### Why scopes instead of just paths?

Scopes encode semantics, not just locations. A scope knows whether it's writable, what its priority is, and how to compute its base path. This enables:

- Clear separation of "where to read" vs "where to write"
- Consistent handling of missing scopes (no project → skip project scope)
- Extensibility for custom hierarchies without reimplementing resolution

### Why separate project directory types from scopes?

Scopes represent the cascade hierarchy (local → project → user → system). Project directory types represent the internal structure of a project (config, cache, data, log). These are orthogonal concerns:

- A tool might have project-level cache and logs even if it only uses the "project" scope for config
- The cascade determines which scope's config file wins; the project dir type determines where cache goes
- Some directory types (cache, log) typically don't participate in cascading at all

### Why three-way returns for resolution?

The `(primary, alternates, error)` return pattern supports both "first match wins" and "merge all" use cases with a single method. The caller decides the semantics:

```go
// First match wins
config, _, _ := cascade.ResolveFile("config.yaml")
load(config)

// Merge all
config, alts, _ := cascade.ResolveFile("config.yaml")
for _, alt := range append([]string{config}, alts...) {
    merge(alt)
}
```

### Why is project detection opt-out rather than opt-in?

Most CLI tools have a project concept. The default configuration assumes this. Tools without projects (daemons, system utilities) explicitly opt out with `WithoutProjectDetection()`. This matches the common case while keeping the API simple.

### Why config-dir-aware detection?

Without it, running a command from inside `.myapp/policies.d/` would walk up and potentially find a parent project's markers. The user is clearly working in *this* project's config; the library should recognize that and return the correct project root. This matches user intent and prevents surprising behavior.

### Why known files?

Tools often have a handful of well-known files with complex variant rules. Rather than making callers remember "config can be `agenthooks.yaml` or `.agenthooks.yaml` in either `.agenthooks/` or `.config/agenthooks/`", the library lets you define this once and refer to it by name:

```go
// Define once
WithProjectFiles(map[string]ProjectFileConfig{
    "config": {Variants: []string{"agenthooks.yaml", ".agenthooks.yaml"}, Dir: "config"},
})

// Use by name
cascade.ResolveKnownFile("config")
cascade.DefaultKnownFile("config")
```

This reduces errors and makes the API more discoverable.

## Relationship to Dirs and FindUp

`Cascade` builds on top of `Dirs` and `FindUp`:

| Concern                     | Component                              |
| --------------------------- | -------------------------------------- |
| Platform directory paths    | `Dirs`                                 |
| Upward directory traversal  | `FindUp*` methods on `Dirs`            |
| Layered scope hierarchy     | `Cascade`                              |
| Project detection           | `Cascade` (uses `FindUp` internally)   |
| Project directory structure | `Cascade` (ProjectDirConfig)           |
| Multi-variant file matching | `Cascade`                              |
| Known file resolution       | `Cascade` (ProjectFileConfig)          |

A typical import just uses `Cascade`:

```go
cascade, _ := toolpaths.NewCascade("myapp")
```

Access the underlying `Dirs` when needed:

```go
// For FindUp or direct platform directory access
dirs := cascade.Dirs()
dirs.FindUp(cwd, ".git")
dirs.UserCacheDir()
```

Most users interact only with `Cascade`. The `Dirs` type is useful for simpler tools that don't need layered resolution, or for accessing `FindUp` methods directly.
