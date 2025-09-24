// Package dedupkeep provides function-based implementations for deciding which duplicate file to keep.
package dedupkeep

import (
	"os"
	"path/filepath"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
)

// OldestFile returns a Func that keeps the oldest file based on modification time.
func OldestFile(fs filesystem.FileSystem) Func {
	return func(paths []string) (toKeep string, toRemove []string) {
		if len(paths) == 0 {
			return "", []string{}
		}

		toKeep = paths[0]
		toRemove = []string{} // Initialize as empty slice, not nil
		oldestTime := getModTime(fs, toKeep)

		for _, path := range paths[1:] {
			modTime := getModTime(fs, path)
			if modTime.ModTime().Before(oldestTime.ModTime()) {
				toRemove = append(toRemove, toKeep)
				toKeep = path
				oldestTime = modTime
				continue
			}
			toRemove = append(toRemove, path)
		}
		return toKeep, toRemove
	}
}

// ShortestPath returns a Func that keeps the file with the shortest path.
func ShortestPath() Func {
	return func(paths []string) (toKeep string, toRemove []string) {
		if len(paths) == 0 {
			return "", []string{}
		}

		toKeep = paths[0]
		toRemove = []string{} // Initialize as empty slice, not nil
		for _, path := range paths[1:] {
			if len(path) < len(toKeep) {
				toRemove = append(toRemove, toKeep)
				toKeep = path
			} else {
				toRemove = append(toRemove, path)
			}
		}
		return toKeep, toRemove
	}
}

// getModTime is a helper function to get modification time with error handling.
func getModTime(fs filesystem.FileSystem, path string) os.FileInfo {
	info, err := fs.Stat(path)
	if err != nil {
		return &DummyFileInfo{name: filepath.Base(path)}
	}
	return info
}
