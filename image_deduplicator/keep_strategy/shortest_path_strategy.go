package keep_strategy

import "sort"

// ShortestPathStrategy keeps the file with the shortest absolute path.
type ShortestPathStrategy struct{}

// Select keeps the file with the shortest path length.
func (s *ShortestPathStrategy) Select(paths []string) (string, []string) {
	sortedPaths := make([]string, len(paths))
	copy(sortedPaths, paths)
	sort.Slice(sortedPaths, func(i, j int) bool {
		if len(sortedPaths[i]) != len(sortedPaths[j]) {
			return len(sortedPaths[i]) < len(sortedPaths[j])
		}
		return sortedPaths[i] < sortedPaths[j] // Fallback to alphabetical
	})
	return sortedPaths[0], sortedPaths[1:]
}
