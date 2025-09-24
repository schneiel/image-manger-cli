// Package dedupaction provides different strategies for handling duplicate files.
package dedupaction

import (
	"github.com/schneiel/ImageManagerGo/core/image"
	shared "github.com/schneiel/ImageManagerGo/core/strategies/shared"
)

// Strategy defines the interface for handling duplicate files in an OOP manner.
// Implementing classes should provide specific logic for different action types.
type Strategy interface {
	// Execute performs the action on the group of identified duplicates
	Execute(original *image.Image, duplicate *image.Image) error
	// GetResources returns the resource manager for this strategy. The caller must call Setup()
	// before executing actions and Teardown() afterward to properly manage resources like files,
	// directories, or CSV writers.
	GetResources() shared.ActionResource
}
