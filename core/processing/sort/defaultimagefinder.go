package sort

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

// DefaultImageFinder is responsible for scanning a directory to find image files
// based on a list of allowed extensions.
type DefaultImageFinder struct {
	allowedExtensions []string
}

// NewDefaultImageFinder creates a new image finder with the specified allowed file extensions.
func NewDefaultImageFinder(allowedExtensions []string) ImageFinder {
	// To ensure case-insensitive matching, convert all extensions to lower case once.
	lowerCaseExtensions := make([]string, len(allowedExtensions))

	for i, ext := range allowedExtensions {
		lowerCaseExtensions[i] = strings.ToLower(ext)
	}

	return &DefaultImageFinder{allowedExtensions: lowerCaseExtensions}
}

// Find scans the given root path recursively and returns a slice of paths
// to files that match the allowed extensions.
func (f *DefaultImageFinder) Find(rootPath string) ([]string, error) {
	var imagePaths []string
	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// This error is typically logged by the caller (ImageProcessor)
			// or handled as a non-fatal error for a single file/directory.
			return err // Propagate errors (e.g., permission denied)
		}

		if !d.IsDir() && f.isImage(path) {
			imagePaths = append(imagePaths, path)
		}
		return nil
	}

	err := filepath.WalkDir(rootPath, walkFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", rootPath, err)
	}

	return imagePaths, nil
}

// isImage checks if a file is an image based on its extension. The check is case-insensitive.
func (f *DefaultImageFinder) isImage(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return slices.Contains(f.allowedExtensions, ext)
}
