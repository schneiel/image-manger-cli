package dedup

import (
	"github.com/schneiel/ImageManagerGo/core/image"
)

// Scanner finds potential duplicate files and groups them.
type Scanner interface {
	Scan(rootPath string) (FileGroup, error)
}

// Hasher calculates image hashes for a given list of file paths.
type Hasher interface {
	HashFiles(files []string) ([]*image.HashInfo, error)
}

// Grouper takes a list of image hashes and groups them into duplicate sets.
type Grouper interface {
	Group(hashes []*image.HashInfo) ([]DuplicateGroup, error)
}

// Strategy defines the action to be taken on duplicate images.
type Strategy interface {
	Execute(original *image.Image, duplicate *image.Image) error
}
