// Package dedupkeep provides different strategies for deciding which duplicate file to keep.
package dedupkeep

// Func defines the function signature for deciding which file to keep.
type Func func(_ []string) (toKeep string, toRemove []string)
