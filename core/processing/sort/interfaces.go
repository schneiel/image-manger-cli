// Package sort provides functionality for analyzing and processing images.
package sort

import (
	"github.com/schneiel/ImageManagerGo/core/image"
	"github.com/schneiel/ImageManagerGo/core/strategies/shared"
)

// Strategy defines the interface for executing sort operations in an OOP manner.
// Implementing classes should provide specific logic for different action types.
type Strategy interface {
	// Execute performs the action on the source and destination paths
	Execute(sourcePath, destinationPath string) error
	// GetResources returns the resource manager for this strategy
	GetResources() shared.ActionResource
}

// ImageFinder defines the interface for finding images in a directory structure.
type ImageFinder interface {
	Find(rootPath string) ([]string, error)
}

// ImageAnalyzer defines the interface for analyzing individual image files.
type ImageAnalyzer interface {
	Analyze(filePath string) (image.Image, error)
}

// ImageProcessor defines the interface for processing multiple images.
type ImageProcessor interface {
	Process(dirPath string) []image.Image
}
