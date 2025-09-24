// Package dedup provides functionality for detecting and managing duplicate images
// using perceptual hashing techniques.
package dedup

// FileGroup represents a set of files that are potential duplicates,
// for example, because they share the same file size.
type FileGroup [][]string
