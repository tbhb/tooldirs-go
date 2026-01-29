package toolpaths

import (
	"path/filepath"
)

// Match represents a found marker during upward traversal.
type Match struct {
	Dir    string // Directory containing the marker
	Marker string // The marker that matched (filename or dirname)
}

// Path returns the full path to the marker.
func (m Match) Path() string {
	return filepath.Join(m.Dir, m.Marker)
}

// FindUp walks up from start, returning the first directory containing any of
// the specified markers. Markers can be files or directories. When multiple
// markers are specified, the walker checks them in order; the first existing
// marker in a directory wins.
func (d *PlatformDirs) FindUp(start string, markers ...string) (string, string, bool) {
	matches := d.walkUp(start, markers, nil, nil, false)
	if len(matches) == 0 {
		return "", "", false
	}
	return matches[0].Dir, matches[0].Marker, true
}

// FindUpFunc walks up from start with a predicate. A marker only matches if it
// exists AND match(markerPath) returns true. This enables content inspection
// without the library needing to understand file formats.
func (d *PlatformDirs) FindUpFunc(
	start string,
	markers []string,
	match func(markerPath string) bool,
) (string, string, bool) {
	matches := d.walkUp(start, markers, nil, match, false)
	if len(matches) == 0 {
		return "", "", false
	}
	return matches[0].Dir, matches[0].Marker, true
}

// FindUpUntil walks up from start, stopping when a directory contains any of
// the stopAt markers. If a directory contains both a target marker and a stop
// marker, the target matches before traversal stops.
func (d *PlatformDirs) FindUpUntil(
	start string,
	markers, stopAt []string,
) (string, string, bool) {
	matches := d.walkUp(start, markers, stopAt, nil, false)
	if len(matches) == 0 {
		return "", "", false
	}
	return matches[0].Dir, matches[0].Marker, true
}

// FindUpUntilFunc combines predicate validation with stop markers.
func (d *PlatformDirs) FindUpUntilFunc(
	start string,
	markers, stopAt []string,
	match func(markerPath string) bool,
) (string, string, bool) {
	matches := d.walkUp(start, markers, stopAt, match, false)
	if len(matches) == 0 {
		return "", "", false
	}
	return matches[0].Dir, matches[0].Marker, true
}

// FindAllUp returns all directories containing any marker, ordered nearest to
// farthest from start.
func (d *PlatformDirs) FindAllUp(start string, markers ...string) []Match {
	return d.walkUp(start, markers, nil, nil, true)
}

// FindAllUpFunc returns all matching directories, filtering through a predicate.
func (d *PlatformDirs) FindAllUpFunc(
	start string,
	markers []string,
	match func(markerPath string) bool,
) []Match {
	return d.walkUp(start, markers, nil, match, true)
}

// FindAllUpUntil collects all matches until traversal encounters a stop marker.
func (d *PlatformDirs) FindAllUpUntil(start string, markers, stopAt []string) []Match {
	return d.walkUp(start, markers, stopAt, nil, true)
}

// FindAllUpUntilFunc combines collection, predicate, and stop behavior.
func (d *PlatformDirs) FindAllUpUntilFunc(
	start string,
	markers, stopAt []string,
	match func(markerPath string) bool,
) []Match {
	return d.walkUp(start, markers, stopAt, match, true)
}

// walkUp is the internal traversal function. It walks from start toward the
// filesystem root, checking for markers in each directory.
func (d *PlatformDirs) walkUp(
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
		if match, found := d.checkMarkers(dir, markers, matchFn); found {
			results = append(results, match)
			if !collectAll {
				return results
			}
		}

		if shouldStop(dir, stopAt) {
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

// cleanAbsPath returns a cleaned absolute path.
func cleanAbsPath(path string) string {
	dir := filepath.Clean(path)
	if !filepath.IsAbs(dir) {
		if abs, err := filepath.Abs(dir); err == nil {
			dir = abs
		}
	}
	return dir
}

// checkMarkers checks if any marker exists in the directory.
func (d *PlatformDirs) checkMarkers(
	dir string,
	markers []string,
	matchFn func(string) bool,
) (Match, bool) {
	for _, m := range markers {
		markerPath := filepath.Join(dir, m)
		if fileExists(markerPath) {
			if matchFn == nil || matchFn(markerPath) {
				return Match{Dir: dir, Marker: m}, true
			}
		}
	}
	return Match{}, false
}

// shouldStop checks if any stop marker exists in the directory.
func shouldStop(dir string, stopAt []string) bool {
	for _, s := range stopAt {
		if fileExists(filepath.Join(dir, s)) {
			return true
		}
	}
	return false
}
