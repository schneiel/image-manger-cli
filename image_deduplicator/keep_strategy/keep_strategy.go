package keep_strategy

// KeepStrategy defines the interface for selecting which file to keep from a group of duplicates.
type KeepStrategy interface {
	// Select chooses one file to keep and returns the rest to be removed.
	Select(paths []string) (toKeep string, toRemove []string)
}
