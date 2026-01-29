package tooldirs_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tbhb/tooldirs-go"
)

// createDirHierarchy creates a directory hierarchy for testing.
// It returns the base temp directory (cleaned up automatically by t.TempDir).
func createDirHierarchy(t *testing.T, structure map[string]string) string {
	t.Helper()
	base := t.TempDir()

	for path, content := range structure {
		fullPath := filepath.Join(base, path)
		dir := filepath.Dir(fullPath)
		err := os.MkdirAll(dir, 0o755)
		require.NoError(t, err, "failed to create directory %s", dir)

		if content == "[dir]" {
			err = os.MkdirAll(fullPath, 0o755)
		} else {
			err = os.WriteFile(fullPath, []byte(content), 0o644)
		}
		require.NoError(t, err, "failed to create %s", fullPath)
	}

	return base
}

func TestFindUp(t *testing.T) {
	t.Run("finds marker in start directory", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"project/go.mod": "module test",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		dir, marker, found := dirs.FindUp(filepath.Join(base, "project"), "go.mod")
		assert.True(t, found)
		assert.Equal(t, filepath.Join(base, "project"), dir)
		assert.Equal(t, "go.mod", marker)
	})

	t.Run("finds marker in parent directory", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"project/go.mod":      "module test",
			"project/src/main.go": "package main",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		dir, marker, found := dirs.FindUp(filepath.Join(base, "project", "src"), "go.mod")
		assert.True(t, found)
		assert.Equal(t, filepath.Join(base, "project"), dir)
		assert.Equal(t, "go.mod", marker)
	})

	t.Run("returns false when marker not found", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"project/src/main.go": "package main",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		dir, marker, found := dirs.FindUp(filepath.Join(base, "project", "src"), "go.mod")
		assert.False(t, found)
		assert.Empty(t, dir)
		assert.Empty(t, marker)
	})

	t.Run("finds first existing marker from multiple candidates", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"project/package.json": "{}",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		// go.mod doesn't exist, package.json does
		dir, marker, found := dirs.FindUp(
			filepath.Join(base, "project"),
			"go.mod",
			"package.json",
		)
		assert.True(t, found)
		assert.Equal(t, filepath.Join(base, "project"), dir)
		assert.Equal(t, "package.json", marker)
	})

	t.Run("prefers earlier markers when multiple exist", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"project/go.mod":       "module test",
			"project/package.json": "{}",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		// Both exist, should return go.mod (first in list)
		dir, marker, found := dirs.FindUp(
			filepath.Join(base, "project"),
			"go.mod",
			"package.json",
		)
		assert.True(t, found)
		assert.Equal(t, filepath.Join(base, "project"), dir)
		assert.Equal(t, "go.mod", marker)
	})

	t.Run("works with directory markers", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"project/.git/config": "git config",
			"project/src/main.go": "package main",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		dir, marker, found := dirs.FindUp(filepath.Join(base, "project", "src"), ".git")
		assert.True(t, found)
		assert.Equal(t, filepath.Join(base, "project"), dir)
		assert.Equal(t, ".git", marker)
	})

	t.Run("returns false for empty markers", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"project/go.mod": "module test",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		dir, marker, found := dirs.FindUp(filepath.Join(base, "project"))
		assert.False(t, found)
		assert.Empty(t, dir)
		assert.Empty(t, marker)
	})
}

func TestFindUpFunc(t *testing.T) {
	t.Run("applies predicate to found marker", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"project/Cargo.toml": "[package]\nname = \"test\"",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		// Predicate checks for [workspace] section
		hasWorkspace := func(path string) bool {
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				return false
			}
			return strings.Contains(string(content), "[workspace]")
		}

		// Should not match because file doesn't have [workspace]
		dir, marker, found := dirs.FindUpFunc(
			filepath.Join(base, "project"),
			[]string{"Cargo.toml"},
			hasWorkspace,
		)
		assert.False(t, found)
		assert.Empty(t, dir)
		assert.Empty(t, marker)
	})

	t.Run("matches when predicate returns true", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"project/Cargo.toml": "[workspace]\nmembers = [\"crate1\"]",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		hasWorkspace := func(path string) bool {
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				return false
			}
			return strings.Contains(string(content), "[workspace]")
		}

		dir, marker, found := dirs.FindUpFunc(
			filepath.Join(base, "project"),
			[]string{"Cargo.toml"},
			hasWorkspace,
		)
		assert.True(t, found)
		assert.Equal(t, filepath.Join(base, "project"), dir)
		assert.Equal(t, "Cargo.toml", marker)
	})

	t.Run("skips markers that fail predicate and finds later ones", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"project/crate/Cargo.toml": "[package]\nname = \"crate\"",
			"project/Cargo.toml":       "[workspace]\nmembers = [\"crate\"]",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		hasWorkspace := func(path string) bool {
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				return false
			}
			return strings.Contains(string(content), "[workspace]")
		}

		dir, marker, found := dirs.FindUpFunc(
			filepath.Join(base, "project", "crate"),
			[]string{"Cargo.toml"},
			hasWorkspace,
		)
		assert.True(t, found)
		assert.Equal(t, filepath.Join(base, "project"), dir)
		assert.Equal(t, "Cargo.toml", marker)
	})
}

func TestFindUpUntil(t *testing.T) {
	t.Run("stops at stop marker", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"repo/.git/config":         "git config",
			"repo/project/src/main.go": "package main",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		// Look for go.mod, stop at .git
		dir, marker, found := dirs.FindUpUntil(
			filepath.Join(base, "repo", "project", "src"),
			[]string{"go.mod"},
			[]string{".git"},
		)
		assert.False(t, found)
		assert.Empty(t, dir)
		assert.Empty(t, marker)
	})

	t.Run("finds marker before stop marker in same directory", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"repo/.git/config": "git config",
			"repo/go.mod":      "module test",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		// Both go.mod and .git are in same directory
		dir, marker, found := dirs.FindUpUntil(
			filepath.Join(base, "repo"),
			[]string{"go.mod"},
			[]string{".git"},
		)
		assert.True(t, found)
		assert.Equal(t, filepath.Join(base, "repo"), dir)
		assert.Equal(t, "go.mod", marker)
	})

	t.Run("finds marker in subdirectory before reaching stop marker", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"repo/.git/config":         "git config",
			"repo/project/go.mod":      "module test",
			"repo/project/src/main.go": "package main",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		dir, marker, found := dirs.FindUpUntil(
			filepath.Join(base, "repo", "project", "src"),
			[]string{"go.mod"},
			[]string{".git"},
		)
		assert.True(t, found)
		assert.Equal(t, filepath.Join(base, "repo", "project"), dir)
		assert.Equal(t, "go.mod", marker)
	})
}

func TestFindUpUntilFunc(t *testing.T) {
	t.Run("combines predicate and stop behavior", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"repo/.git/config":      "git config",
			"repo/crate/Cargo.toml": "[package]\nname = \"crate\"",
			"repo/Cargo.toml":       "[workspace]\nmembers = [\"crate\"]",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		hasWorkspace := func(path string) bool {
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				return false
			}
			return strings.Contains(string(content), "[workspace]")
		}

		dir, marker, found := dirs.FindUpUntilFunc(
			filepath.Join(base, "repo", "crate"),
			[]string{"Cargo.toml"},
			[]string{".git"},
			hasWorkspace,
		)
		assert.True(t, found)
		assert.Equal(t, filepath.Join(base, "repo"), dir)
		assert.Equal(t, "Cargo.toml", marker)
	})
}

func TestFindAllUp(t *testing.T) {
	t.Run("collects all matching directories", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"project/.myconfig":         "global config",
			"project/src/.myconfig":     "src config",
			"project/src/pkg/.myconfig": "pkg config",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		matches := dirs.FindAllUp(filepath.Join(base, "project", "src", "pkg"), ".myconfig")

		require.Len(t, matches, 3)
		// Results should be nearest to farthest
		assert.Equal(t, filepath.Join(base, "project", "src", "pkg"), matches[0].Dir)
		assert.Equal(t, filepath.Join(base, "project", "src"), matches[1].Dir)
		assert.Equal(t, filepath.Join(base, "project"), matches[2].Dir)

		for _, m := range matches {
			assert.Equal(t, ".myconfig", m.Marker)
		}
	})

	t.Run("returns empty slice when no matches", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"project/src/main.go": "package main",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		matches := dirs.FindAllUp(filepath.Join(base, "project", "src"), ".myconfig")
		assert.Empty(t, matches)
	})
}

func TestFindAllUpFunc(t *testing.T) {
	t.Run("filters matches through predicate", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"project/.editorconfig":     "root = false",
			"project/src/.editorconfig": "indent_size = 2",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		// Only match .editorconfig files that don't have root = true
		notRoot := func(path string) bool {
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				return false
			}
			return !strings.Contains(string(content), "root = true")
		}

		matches := dirs.FindAllUpFunc(
			filepath.Join(base, "project", "src"),
			[]string{".editorconfig"},
			notRoot,
		)
		require.Len(t, matches, 2)
	})
}

func TestFindAllUpUntil(t *testing.T) {
	t.Run("collects until stop marker", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"repo/.git/config":           "git config",
			"repo/.myconfig":             "repo config",
			"repo/project/.myconfig":     "project config",
			"repo/project/src/.myconfig": "src config",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		matches := dirs.FindAllUpUntil(
			filepath.Join(base, "repo", "project", "src"),
			[]string{".myconfig"},
			[]string{".git"},
		)

		require.Len(t, matches, 3)
		assert.Equal(t, filepath.Join(base, "repo", "project", "src"), matches[0].Dir)
		assert.Equal(t, filepath.Join(base, "repo", "project"), matches[1].Dir)
		assert.Equal(t, filepath.Join(base, "repo"), matches[2].Dir)
	})
}

func TestFindAllUpUntilFunc(t *testing.T) {
	t.Run("combines all behaviors", func(t *testing.T) {
		base := createDirHierarchy(t, map[string]string{
			"repo/.git/config":           "git config",
			"repo/.editorconfig":         "root = true",
			"repo/project/.editorconfig": "indent = 2",
		})

		dirs, err := tooldirs.New("testapp")
		require.NoError(t, err)

		// Match only non-root editorconfigs, stop at .git
		notRoot := func(path string) bool {
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				return false
			}
			return !strings.Contains(string(content), "root = true")
		}

		matches := dirs.FindAllUpUntilFunc(
			filepath.Join(base, "repo", "project"),
			[]string{".editorconfig"},
			[]string{".git"},
			notRoot,
		)

		require.Len(t, matches, 1)
		assert.Equal(t, filepath.Join(base, "repo", "project"), matches[0].Dir)
	})
}

func TestMatch(t *testing.T) {
	t.Run("Path returns full path", func(t *testing.T) {
		m := tooldirs.Match{
			Dir:    "/home/user/project",
			Marker: "go.mod",
		}
		assert.Equal(t, filepath.Join("/home/user/project", "go.mod"), m.Path())
	})
}

// Tests using FakeDirs

func TestFakeDirsFindUp(t *testing.T) {
	t.Run("uses ExistingFiles map", func(t *testing.T) {
		fake := tooldirs.NewFakeDirs("/base")
		fake.SetExisting("/home/user/project/go.mod")

		dir, marker, found := fake.FindUp("/home/user/project/src", "go.mod")
		assert.True(t, found)
		assert.Equal(t, "/home/user/project", dir)
		assert.Equal(t, "go.mod", marker)
	})

	t.Run("returns false when marker not in ExistingFiles", func(t *testing.T) {
		fake := tooldirs.NewFakeDirs("/base")
		// Don't set any existing files

		dir, marker, found := fake.FindUp("/home/user/project/src", "go.mod")
		assert.False(t, found)
		assert.Empty(t, dir)
		assert.Empty(t, marker)
	})

	t.Run("respects stop markers", func(t *testing.T) {
		fake := tooldirs.NewFakeDirs("/base")
		fake.SetExisting("/home/user/.git")
		fake.SetExisting("/home/go.mod")

		// Should stop at .git in /home/user, not find go.mod in /home
		dir, marker, found := fake.FindUpUntil(
			"/home/user/project",
			[]string{"go.mod"},
			[]string{".git"},
		)
		assert.False(t, found)
		assert.Empty(t, dir)
		assert.Empty(t, marker)
	})

	t.Run("finds marker in same dir as stop marker", func(t *testing.T) {
		fake := tooldirs.NewFakeDirs("/base")
		fake.SetExisting("/home/user/project/.git")
		fake.SetExisting("/home/user/project/go.mod")

		dir, marker, found := fake.FindUpUntil(
			"/home/user/project",
			[]string{"go.mod"},
			[]string{".git"},
		)
		assert.True(t, found)
		assert.Equal(t, "/home/user/project", dir)
		assert.Equal(t, "go.mod", marker)
	})
}

func TestFakeDirsFindAllUp(t *testing.T) {
	t.Run("collects all matching directories", func(t *testing.T) {
		fake := tooldirs.NewFakeDirs("/base")
		fake.SetExisting("/home/user/project/.myconfig")
		fake.SetExisting("/home/user/.myconfig")
		fake.SetExisting("/home/.myconfig")

		matches := fake.FindAllUp("/home/user/project", ".myconfig")

		require.Len(t, matches, 3)
		assert.Equal(t, "/home/user/project", matches[0].Dir)
		assert.Equal(t, "/home/user", matches[1].Dir)
		assert.Equal(t, "/home", matches[2].Dir)
	})
}
