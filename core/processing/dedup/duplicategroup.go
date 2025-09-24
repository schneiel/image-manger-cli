// Package dedup provides functionality for detecting and managing duplicate images
// using perceptual hashing techniques.
package dedup

// DuplicateGroup represents a set of files that have been identified
// as duplicates of each other based on their hash similarity.
type DuplicateGroup []string
