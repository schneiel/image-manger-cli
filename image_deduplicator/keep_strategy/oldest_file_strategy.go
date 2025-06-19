package keep_strategy

import (
	"os"
	"sort"
)

// OldestFileStrategy keeps the file with the oldest modification date.
type OldestFileStrategy struct{}

// Select keeps the file with the oldest modification time.
func (s *OldestFileStrategy) Select(paths []string) (string, []string) {
	sortedPaths := make([]string, len(paths))
	copy(sortedPaths, paths)
	sort.Slice(sortedPaths, func(i, j int) bool {
		infoI, errI := os.Stat(sortedPaths[i])
		infoJ, errJ := os.Stat(sortedPaths[j])
		if errI != nil || errJ != nil {
			return sortedPaths[i] < sortedPaths[j] // Fallback to path
		}
		return infoI.ModTime().Before(infoJ.ModTime())
	})
	return sortedPaths[0], sortedPaths[1:]
}
